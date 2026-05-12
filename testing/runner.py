"""
@file runner.py
@brief Vstupní CLI pro benchmarkové a experimentální testování platformy RIoT.
@author Vojtěch Hubáček
@defgroup riot_testing Testing
@ingroup riot
@see README.md
"""

import argparse
import shutil
import subprocess
import sys
import time
from dataclasses import replace
from pathlib import Path

import yaml

REPO_ROOT = Path(__file__).resolve().parents[1]
if str(REPO_ROOT) not in sys.path:
    sys.path.insert(0, str(REPO_ROOT))

from testing.core.config_loader import load_testing_config
from testing.core.context import build_runtime_context
from testing.core.reporter import Reporter
from testing.core.result_store import ResultStore
from testing.datasets.dataset_builder import DatasetBuilder
from testing.experiments import build_experiment
from testing.generators.generator_manager import GeneratorManager
from testing.generators.transport.direct_isc_adapter import DirectISCAdapter
from testing.setup.setup_env import EnvironmentSetupService
from testing.setup.setup_entities import EntitySetupService
from testing.setup.setup_kpis import KPISetupService

REPROCESS_QUEUE_NAMES = [
    "kpi-reprocess-requests",
    "time-series-reprocess-read-request",
    "time-series-reprocess-read-response",
    "time-series-delete-requests-tsdb",
    "time-series-delete-requests-mpu",
    "kpi-fulfillment-check-results",
    "time-series-kpi-results",
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="RIoT testing harness")
    subparsers = parser.add_subparsers(dest="command", required=True)

    setup = subparsers.add_parser("setup")
    setup.add_argument(
        "--phase",
        choices=["all", "types", "refresh", "kpis"],
        default="all",
    )
    subparsers.add_parser("validate")

    prepare = subparsers.add_parser("prepare-dataset")
    prepare.add_argument("--dataset", required=True)
    prepare.add_argument("--rebuild", action="store_true")
    prepare.add_argument("--resume", action="store_true")

    prime = subparsers.add_parser("prime")
    prime.add_argument("--dataset")
    prime.add_argument("--experiment")
    prime.add_argument("--scenario")
    prime.add_argument("--all", action="store_true")
    prime.add_argument("--rebuild", action="store_true")
    prime.add_argument("--resume", action="store_true")
    prime.add_argument("--full-reset", action="store_true", help="Reset Docker volumes/runtime state before priming the dataset.")
    prime.add_argument("--api-key", help="Full-rights API key used after a full reset.")
    prime.add_argument("--validate-timeout-seconds", type=float, default=900.0)
    prime.add_argument("--validate-poll-seconds", type=float, default=5.0)
    prime.add_argument("--skip-down-v", action="store_true")
    prime.add_argument("--skip-clean-docker-dir", action="store_true")
    prime.add_argument("--skip-prune", action="store_true")
    prime.add_argument("--skip-build", action="store_true")
    prime.add_argument("--keep-outputs", action="store_true")

    inject = subparsers.add_parser("inject-window")
    inject.add_argument("--source", required=True)
    inject.add_argument("--minutes", type=int, required=True)
    inject.add_argument("--profile", choices=["normal", "high", "stress"], default="normal")
    inject.add_argument("--bootstrap", action="store_true")
    inject.add_argument("--dry-run", action="store_true")

    run = subparsers.add_parser("run")
    run.add_argument("--experiment")
    run.add_argument("--scenario")
    run.add_argument("--all", action="store_true")
    run.add_argument("--dry-run", action="store_true")

    report = subparsers.add_parser("report")
    report.add_argument("--experiment")
    report.add_argument(
        "--format",
        choices=["console", "md", "csv", "json"],
        default="console",
    )

    suite = subparsers.add_parser("suite")
    suite.add_argument("--dataset", default="D3")
    suite.add_argument("--api-key")
    suite.add_argument("--dry-run", action="store_true")
    suite.add_argument("--resume", action="store_true")
    suite.add_argument("--validate-timeout-seconds", type=float, default=900.0)
    suite.add_argument("--validate-poll-seconds", type=float, default=5.0)
    suite.add_argument(
        "--experiments",
        nargs="+",
        default=["E1", "E2", "E3", "E4"],
        help="Ordered list of experiments to execute in the suite.",
    )
    suite.add_argument("--skip-down-v", action="store_true")
    suite.add_argument("--skip-clean-docker-dir", action="store_true")
    suite.add_argument("--skip-prune", action="store_true")
    suite.add_argument("--skip-build", action="store_true")
    suite.add_argument("--keep-outputs", action="store_true")
    return parser.parse_args()


def resolve_scenarios(config, args: argparse.Namespace):
    selected = []
    if args.all:
        for experiment in config.experiments.values():
            if experiment.enabled:
                selected.extend(experiment.scenarios)
        return selected
    if args.experiment:
        experiment = config.experiments.get(args.experiment)
        if experiment is None:
            raise SystemExit(f"Unknown experiment: {args.experiment}")
        return experiment.scenarios
    if args.scenario:
        for experiment in config.experiments.values():
            for scenario in experiment.scenarios:
                if scenario.id == args.scenario:
                    return [scenario]
        raise SystemExit(f"Unknown scenario: {args.scenario}")
    raise SystemExit("Specify --all, --experiment or --scenario")


def command_setup(ctx, config) -> int:
    env_service = EnvironmentSetupService()
    entity_service = EntitySetupService()
    kpi_service = KPISetupService()

    checks = env_service.validate(ctx)
    print("Environment checks:")
    for name, ok in checks.items():
        print(f"- {name}: {'ok' if ok else 'failed'}")

    phase = getattr(ctx, "setup_phase", "all")
    if phase in {"all", "refresh", "kpis"} and not checks.get("rabbitmq_management", False):
        print("Warning: rabbitmq_management is unavailable; refresh/kpi-dependent steps may fail later.")
    state = ctx.state_manager.load_runtime_state()
    if phase in {"all", "types", "refresh", "kpis"}:
        # Always refresh SDType IDs from the live backend before refreshing
        # instances or creating KPIs. After DB/container restarts the persisted
        # numeric IDs in testing/.state can be stale even if the SDType UIDs
        # stayed the same.
        state = entity_service.ensure_sd_types(ctx, config)
    if phase in {"all", "refresh"}:
        for source_name in config.sources.keys():
            if source_name in state.get("sources", {}):
                entity_service.refresh_instances_for_source(ctx, state, source_name)
    if phase in {"all", "kpis"}:
        state = kpi_service.ensure_kpis(ctx, config, state)
    ctx.state_manager.save_runtime_state(state)
    print("Setup completed.")
    return 0


def command_validate(ctx) -> int:
    env_service = EnvironmentSetupService()
    checks = env_service.validate(ctx)
    print("Validation:")
    failed = False
    for name, ok in checks.items():
        print(f"- {name}: {'ok' if ok else 'failed'}")
        failed = failed or not ok
    return 1 if failed else 0


def command_prepare_dataset(ctx, config, args: argparse.Namespace) -> int:
    dataset = config.datasets.get(args.dataset)
    if dataset is None:
        raise SystemExit(f"Unknown dataset: {args.dataset}")

    builder = DatasetBuilder()
    record = builder.build(ctx, config, dataset, rebuild=args.rebuild, resume=args.resume)
    materialized = builder.resolve_checkpoint_ids(config, dataset)
    print("Dataset prepared:")
    print(f"- id: {record.dataset_id}")
    print(f"- raw points: {record.raw_points}")
    print(f"- kpi points: {record.kpi_points}")
    print(f"- storage bytes: {record.storage_size_bytes}")
    print(f"- materialized checkpoints: {', '.join(materialized)}")
    return 0


def command_prime(repo_root: Path, ctx, config, args: argparse.Namespace) -> int:
    env_service = EnvironmentSetupService()
    entity_service = EntitySetupService()
    kpi_service = KPISetupService()

    if getattr(args, "full_reset", False):
        _run_suite_cleanup(repo_root, args)
        _run_compose_up(repo_root, build=not getattr(args, "skip_build", False))
        api_key = args.api_key.strip() if getattr(args, "api_key", None) else _prompt_for_api_key()
        _persist_api_key(repo_root, api_key)
        config = _with_api_key(config, api_key)
        ctx = build_runtime_context(repo_root, config)
        _wait_until_valid(
            ctx,
            timeout_seconds=getattr(args, "validate_timeout_seconds", 900.0),
            poll_seconds=getattr(args, "validate_poll_seconds", 5.0),
        )
        _run_setup_phase(ctx, config, "types")
        _run_setup_phase(ctx, config, "refresh")
        _run_setup_phase(ctx, config, "kpis")

    checks = env_service.validate(ctx)
    print("Environment checks:")
    for name, ok in checks.items():
        print(f"- {name}: {'ok' if ok else 'failed'}")

    target_dataset_id = _resolve_prime_dataset_id(config, args)
    if target_dataset_id:
        dataset = config.datasets.get(target_dataset_id)
        if dataset is None:
            raise SystemExit(f"Unknown dataset: {target_dataset_id}")
        builder = DatasetBuilder()
        record = builder.build(ctx, config, dataset, rebuild=args.rebuild, resume=args.resume)
        materialized = builder.resolve_checkpoint_ids(config, dataset)
        print("Dataset prepared:")
        print(f"- id: {record.dataset_id}")
        print(f"- raw points: {record.raw_points}")
        print(f"- kpi points: {record.kpi_points}")
        print(f"- storage bytes: {record.storage_size_bytes}")
        print(f"- materialized checkpoints: {', '.join(materialized)}")

    state = entity_service.ensure_sd_types(ctx, config)
    for source_name in config.sources.keys():
        if source_name in state.get("sources", {}):
            entity_service.refresh_instances_for_source(ctx, state, source_name)
    state = kpi_service.ensure_kpis(ctx, config, state)
    if target_dataset_id:
        requested = kpi_service.enqueue_reprocess_for_kpis(ctx, config, state)
        if requested:
            print(f"Triggered KPI reprocess after dataset prime: {requested} definitions.")
            _wait_for_reprocess_clean(ctx)
    ctx.state_manager.save_runtime_state(state)
    print("Prime completed.")
    return 0


def command_inject_window(ctx, config, args: argparse.Namespace) -> int:
    if args.source not in config.sources:
        raise SystemExit(f"Unknown source: {args.source}")
    if args.source not in {"waze", "ndic", "mhd"}:
        raise SystemExit("Direct ISC injection is currently implemented only for sources 'waze', 'ndic' and 'mhd'")
    if not args.dry_run and not ctx.rabbitmq_amqp_client.is_configured():
        raise SystemExit("RabbitMQ AMQP client is not configured")

    manager = GeneratorManager.from_config(config.sources)
    events = manager.generate_window_events(
        duration_minutes=args.minutes,
        load_profile=args.profile,
        source_name=args.source,
        include_bootstrap=args.bootstrap,
    )
    if args.dry_run:
        unique_instances = len({event.instance_uid for event in events})
        print("Synthetic window preview:")
        print(f"- source: {args.source}")
        print(f"- profile: {args.profile}")
        print(f"- minutes: {args.minutes}")
        print(f"- bootstrap included: {'yes' if args.bootstrap else 'no'}")
        print(f"- generated events: {len(events)}")
        print(f"- unique instances: {unique_instances}")
        return 0

    runtime_state = ctx.state_manager.load_runtime_state()
    adapter = DirectISCAdapter(ctx.rabbitmq_amqp_client, runtime_state)
    result = adapter.push_events(
        source_name=args.source,
        sd_type_uid=config.sources[args.source].payload["sd_type_uid"],
        events=events,
    )
    ctx.state_manager.save_runtime_state(runtime_state)
    print("Synthetic window injected:")
    print(f"- source: {args.source}")
    print(f"- profile: {args.profile}")
    print(f"- minutes: {args.minutes}")
    print(f"- generated events: {len(events)}")
    print(f"- registrations published: {result['registrations']}")
    print(f"- state messages published: {result['published_events']}")
    return 0


def command_run(ctx, config, args: argparse.Namespace) -> int:
    scenarios = resolve_scenarios(config, args)
    reporter = Reporter()
    store = ResultStore(ctx)

    if args.dry_run:
        print("Planned scenarios:")
        for scenario in scenarios:
            print(
                f"- {scenario.id}: experiment={scenario.experiment_id} "
                f"type={scenario.experiment_type} dataset={scenario.dataset_id}"
            )
        return 0

    all_results = []
    for scenario in scenarios:
        experiment = build_experiment(ctx, scenario.experiment_type)
        results = experiment.run(ctx, scenario)
        store.save_results(results)
        all_results.extend(results)

    summaries = reporter.summarize(all_results)
    reporter.print_console_summary(summaries)
    reporter.save_markdown(summaries, ctx.run.output_dir / "reports" / "latest_summary.md")
    reporter.save_csv(summaries, ctx.run.output_dir / "reports" / "latest_summary.csv")
    return 0


def _resolve_prime_dataset_id(config, args: argparse.Namespace) -> str | None:
    if getattr(args, "dataset", None):
        return args.dataset

    scenarios = []
    if getattr(args, "all", False):
        for experiment in config.experiments.values():
            if experiment.enabled:
                scenarios.extend(experiment.scenarios)
    elif getattr(args, "experiment", None):
        experiment = config.experiments.get(args.experiment)
        if experiment is None:
            raise SystemExit(f"Unknown experiment: {args.experiment}")
        scenarios = experiment.scenarios
    elif getattr(args, "scenario", None):
        scenarios = resolve_scenarios(config, args)

    dataset_ids = [scenario.dataset_id for scenario in scenarios if scenario.dataset_id]
    if not dataset_ids:
        return None

    datasets = [config.datasets[dataset_id] for dataset_id in dataset_ids]
    datasets.sort(key=lambda item: item.history_duration_hours)
    return datasets[-1].id


def command_report(ctx, args: argparse.Namespace) -> int:
    store = ResultStore(ctx)
    reporter = Reporter()
    results = store.load_results(experiment_id=args.experiment)
    summaries = reporter.summarize(results)

    if args.format == "console":
        reporter.print_console_summary(summaries)
    elif args.format == "md":
        target = ctx.run.output_dir / "reports" / "summary.md"
        reporter.save_markdown(summaries, target)
        print(target)
    elif args.format == "csv":
        target = ctx.run.output_dir / "reports" / "summary.csv"
        reporter.save_csv(summaries, target)
        print(target)
    else:
        target = ctx.run.output_dir / "reports" / "summary.json"
        store.save_summaries_json(summaries, target)
        print(target)
    return 0


def command_suite(repo_root: Path, config, args: argparse.Namespace) -> int:
    if args.dry_run:
        print("Planned suite flow:")
        print("- text order: E1 storage, E2 ingest, E3 reprocess grouped by D1 -> D2 -> D3, then E4 history export")
        print("- initial cleanup: docker compose down -v, remove RIoT/docker, remove testing/.state, remove testing/outputs, docker system prune -a --volumes -f")
        print("- docker compose up -d --build")
        print("- prompt for full-rights API key")
        print("- wait until validate is fully green")
        print("- setup phases: types, refresh, kpis")
        print("- prime datasets as needed for storage, reprocess and history export")
        if "E3" in args.experiments:
            print("- E3 runs grouped by dataset with full runtime reset between groups:")
            for dataset_id, scenario_ids in _group_reprocess_scenarios_by_dataset(config).items():
                print(f"  - {dataset_id}: {', '.join(scenario_ids)}")
        if args.resume:
            print("- resume mode: keep prior raw results and skip already finished scenarios")
        return 0

    if not args.resume:
        _run_suite_cleanup(repo_root, args)
    _run_compose_up(repo_root, build=not args.skip_build)

    api_key = args.api_key.strip() if args.api_key else _prompt_for_api_key()
    _persist_api_key(repo_root, api_key)
    config = _with_api_key(config, api_key)
    ctx = build_runtime_context(repo_root, config)

    _wait_until_valid(ctx, timeout_seconds=args.validate_timeout_seconds, poll_seconds=args.validate_poll_seconds)
    _run_setup_phase(ctx, config, "types")
    _run_setup_phase(ctx, config, "refresh")
    _run_setup_phase(ctx, config, "kpis")

    store = ResultStore(ctx)
    reporter = Reporter()
    completed = _completed_scenarios(store, config)
    current_live_dataset_id: str | None = None
    requested = [experiment_id for experiment_id in ["E1", "E2", "E3", "E4"] if experiment_id in args.experiments]

    if "E1" in requested and not _experiment_complete("E1", completed, config):
        ctx, config, current_live_dataset_id = _ensure_dataset_ready(
            repo_root=repo_root,
            ctx=ctx,
            config=config,
            dataset_id="D3",
            current_live_dataset_id=current_live_dataset_id,
            validate_timeout_seconds=args.validate_timeout_seconds,
            validate_poll_seconds=args.validate_poll_seconds,
            prompt_api_key=not args.api_key,
            fixed_api_key=args.api_key.strip() if args.api_key else None,
            reset_runtime=False,
            resume=args.resume,
        )
        for scenario in config.experiments["E1"].scenarios:
            if _scenario_complete(scenario, completed):
                continue
            _run_suite_scenario(ctx, scenario, store, reporter)
            completed[(scenario.experiment_id, scenario.id)] = scenario.repetitions

    if "E2" in requested and not _experiment_complete("E2", completed, config):
        _restart_benchmark_services(repo_root)
        _wait_until_valid(ctx, timeout_seconds=args.validate_timeout_seconds, poll_seconds=args.validate_poll_seconds)
        _run_setup_phase(ctx, config, "refresh")
        _run_setup_phase(ctx, config, "kpis")
        current_live_dataset_id = None
        for scenario in config.experiments["E2"].scenarios:
            if _scenario_complete(scenario, completed):
                continue
            _run_suite_scenario(ctx, scenario, store, reporter)
            completed[(scenario.experiment_id, scenario.id)] = scenario.repetitions

    if "E3" in requested:
        grouped_scenarios = _group_reprocess_scenarios_by_dataset(config)
        for dataset_id, scenario_ids in grouped_scenarios.items():
            pending = [
                next(item for item in config.experiments["E3"].scenarios if item.id == scenario_id)
                for scenario_id in scenario_ids
                if not _scenario_complete(next(item for item in config.experiments["E3"].scenarios if item.id == scenario_id), completed)
            ]
            if not pending:
                continue
            ctx, config, current_live_dataset_id = _ensure_dataset_ready(
                repo_root=repo_root,
                ctx=ctx,
                config=config,
                dataset_id=dataset_id,
                current_live_dataset_id=current_live_dataset_id,
                validate_timeout_seconds=args.validate_timeout_seconds,
                validate_poll_seconds=args.validate_poll_seconds,
                prompt_api_key=not args.api_key,
                fixed_api_key=args.api_key.strip() if args.api_key else None,
                reset_runtime=True,
                resume=args.resume,
            )
            for scenario in pending:
                _run_suite_scenario(ctx, scenario, store, reporter)
                completed[(scenario.experiment_id, scenario.id)] = scenario.repetitions

    if "E4" in requested and not _experiment_complete("E4", completed, config):
        if current_live_dataset_id != "D3":
            ctx, config, current_live_dataset_id = _ensure_dataset_ready(
                repo_root=repo_root,
                ctx=ctx,
                config=config,
                dataset_id="D3",
                current_live_dataset_id=current_live_dataset_id,
                validate_timeout_seconds=args.validate_timeout_seconds,
                validate_poll_seconds=args.validate_poll_seconds,
                prompt_api_key=not args.api_key,
                fixed_api_key=args.api_key.strip() if args.api_key else None,
                reset_runtime=True,
                resume=args.resume,
            )
        for scenario in config.experiments["E4"].scenarios:
            if _scenario_complete(scenario, completed):
                continue
            _run_suite_scenario(ctx, scenario, store, reporter)
            completed[(scenario.experiment_id, scenario.id)] = scenario.repetitions

    all_results = store.load_results()
    summaries = reporter.summarize(all_results)
    reporter.print_console_summary(summaries)
    reports_dir = ctx.run.output_dir / "reports"
    md_path = reports_dir / "suite_summary.md"
    csv_path = reports_dir / "suite_summary.csv"
    json_path = reports_dir / "suite_summary.json"
    reporter.save_markdown(summaries, md_path)
    reporter.save_csv(summaries, csv_path)
    store.save_summaries_json(summaries, json_path)
    print("Suite completed.")
    print(f"- markdown report: {md_path}")
    print(f"- csv report: {csv_path}")
    print(f"- json report: {json_path}")
    return 0


def _with_api_key(config, api_key: str):
    return replace(config, env=replace(config.env, api_key=api_key))


def _persist_api_key(repo_root: Path, api_key: str) -> None:
    env_path = repo_root / "testing" / "config" / "env.yaml"
    with env_path.open("r", encoding="utf-8") as handle:
        payload = yaml.safe_load(handle) or {}
    payload["api_key"] = api_key
    with env_path.open("w", encoding="utf-8") as handle:
        yaml.safe_dump(payload, handle, allow_unicode=False, sort_keys=False)


def _prompt_for_api_key() -> str:
    print("Create a full-rights API key in the RIoT UI, then paste it here.")
    while True:
        value = input("API key: ").strip()
        if value:
            return value
        print("API key must not be empty.")


def _run_setup_phase(ctx, config, phase: str) -> None:
    ctx.setup_phase = phase
    code = command_setup(ctx, config)
    if code != 0:
        raise RuntimeError(f"Setup phase failed: {phase}")


def _wait_until_valid(ctx, timeout_seconds: float, poll_seconds: float) -> None:
    env_service = EnvironmentSetupService()
    started = time.perf_counter()
    last_checks = None
    while True:
        checks = env_service.validate(ctx)
        if checks != last_checks:
            print("Validation:")
            for name, ok in checks.items():
                print(f"- {name}: {'ok' if ok else 'failed'}")
            last_checks = checks
        if all(checks.values()):
            break
        if time.perf_counter() - started > timeout_seconds:
            raise TimeoutError(f"Environment did not become valid within {timeout_seconds} seconds. Last checks={checks!r}")
        time.sleep(poll_seconds)

    _wait_for_isc_consumers_ready(ctx, timeout_seconds=timeout_seconds, poll_seconds=poll_seconds)


def _wait_for_isc_consumers_ready(ctx, timeout_seconds: float, poll_seconds: float) -> None:
    if not ctx.rabbitmq_management_client.is_configured():
        return

    required_consumers = {
        "sd-type-registration-requests": 1,
        "sd-instance-registration-requests": 1,
        "set-of-sd-types-updates": 1,
        "message-processing-unit-connection-notifications": 1,
    }
    print("Waiting for ISC consumers:")
    ctx.rabbitmq_management_client.wait_for_queue_consumers(
        required_consumers,
        poll_interval_seconds=poll_seconds,
        timeout_seconds=timeout_seconds,
    )
    for queue_name, minimum_consumers in required_consumers.items():
        print(f"- {queue_name}: >= {minimum_consumers} consumer ready")


def _run_suite_cleanup(repo_root: Path, args: argparse.Namespace) -> None:
    if not args.skip_down_v:
        _run_checked(["docker", "compose", "down", "-v"], cwd=repo_root)
    if not args.skip_clean_docker_dir:
        shutil.rmtree(repo_root / "docker", ignore_errors=True)
    shutil.rmtree(repo_root / "testing" / ".state", ignore_errors=True)
    if not args.keep_outputs:
        shutil.rmtree(repo_root / "testing" / "outputs", ignore_errors=True)
    if not args.skip_prune:
        _run_checked(["docker", "system", "prune", "-a", "--volumes", "-f"], cwd=repo_root)


def _run_compose_up(repo_root: Path, build: bool) -> None:
    command = ["docker", "compose", "up", "-d"]
    if build:
        command.append("--build")
    _run_checked(command, cwd=repo_root)


def _restart_benchmark_services(repo_root: Path) -> None:
    _run_checked(
        [
            "docker",
            "compose",
            "restart",
            "riot-backend-core",
            "riot-message-processing-unit",
            "riot-time-series-store",
        ],
        cwd=repo_root,
    )


def _ensure_dataset_ready(
    repo_root: Path,
    ctx,
    config,
    dataset_id: str,
    current_live_dataset_id: str | None,
    validate_timeout_seconds: float,
    validate_poll_seconds: float,
    prompt_api_key: bool,
    fixed_api_key: str | None,
    reset_runtime: bool,
    resume: bool,
):
    if reset_runtime:
        _reset_runtime_stack(repo_root)
        _run_compose_up(repo_root, build=False)
        api_key = fixed_api_key if fixed_api_key is not None else _prompt_for_api_key()
        _persist_api_key(repo_root, api_key)
        config = _with_api_key(config, api_key)
        ctx = build_runtime_context(repo_root, config)
        _wait_until_valid(ctx, timeout_seconds=validate_timeout_seconds, poll_seconds=validate_poll_seconds)
        _run_setup_phase(ctx, config, "types")
        _run_setup_phase(ctx, config, "refresh")
        _run_setup_phase(ctx, config, "kpis")
        current_live_dataset_id = None
    if current_live_dataset_id != dataset_id:
        prime_args = argparse.Namespace(
            dataset=dataset_id,
            experiment=None,
            scenario=None,
            all=False,
            rebuild=True,
            resume=resume,
        )
        command_prime(repo_root, ctx, config, prime_args)
        _run_setup_phase(ctx, config, "refresh")
        _run_setup_phase(ctx, config, "kpis")
        _wait_for_reprocess_clean(ctx)
        current_live_dataset_id = dataset_id
    return ctx, config, current_live_dataset_id


def _run_suite_scenario(ctx, scenario, store: ResultStore, reporter: Reporter) -> None:
    experiment = build_experiment(ctx, scenario.experiment_type)
    results = experiment.run(ctx, scenario)
    store.save_results(results)
    _update_suite_reports(ctx, store, reporter)


def _update_suite_reports(ctx, store: ResultStore, reporter: Reporter) -> None:
    all_results = store.load_results()
    summaries = reporter.summarize(all_results)
    reports_dir = ctx.run.output_dir / "reports"
    reporter.save_markdown(summaries, reports_dir / "suite_summary.md")
    reporter.save_csv(summaries, reports_dir / "suite_summary.csv")
    store.save_summaries_json(summaries, reports_dir / "suite_summary.json")


def _wait_for_reprocess_clean(ctx, timeout_seconds: float = 1800.0) -> None:
    if not ctx.rabbitmq_management_client.is_configured():
        return
    ctx.rabbitmq_management_client.wait_for_queues_idle(
        REPROCESS_QUEUE_NAMES,
        poll_interval_seconds=ctx.env.poll_interval_seconds,
        consecutive_idle_polls=3,
        timeout_seconds=timeout_seconds,
    )


def _reset_runtime_stack(repo_root: Path) -> None:
    _run_checked(["docker", "compose", "down", "-v"], cwd=repo_root)
    shutil.rmtree(repo_root / "docker", ignore_errors=True)
    runtime_state_path = repo_root / "testing" / ".state" / "runtime_state.json"
    runtime_state_path.unlink(missing_ok=True)


def _group_reprocess_scenarios_by_dataset(config) -> dict[str, list[str]]:
    experiment = config.experiments["E3"]
    grouped: dict[str, list[str]] = {}
    ordered = sorted(
        experiment.scenarios,
        key=lambda scenario: (
            int(config.datasets[scenario.dataset_id].history_duration_hours),
            scenario.id,
        ),
    )
    for scenario in ordered:
        grouped.setdefault(scenario.dataset_id, []).append(scenario.id)
    return grouped


def _completed_scenarios(store: ResultStore, config) -> dict[tuple[str, str], int]:
    counts: dict[tuple[str, str], int] = {}
    for result in store.load_results():
        key = (result.experiment_id, result.scenario_id)
        counts[key] = counts.get(key, 0) + 1
    completed: dict[tuple[str, str], int] = {}
    for experiment in config.experiments.values():
        for scenario in experiment.scenarios:
            key = (scenario.experiment_id, scenario.id)
            if counts.get(key, 0) >= scenario.repetitions:
                completed[key] = scenario.repetitions
    return completed


def _scenario_complete(scenario, completed: dict[tuple[str, str], int]) -> bool:
    return completed.get((scenario.experiment_id, scenario.id), 0) >= scenario.repetitions


def _experiment_complete(experiment_id: str, completed: dict[tuple[str, str], int], config) -> bool:
    experiment = config.experiments[experiment_id]
    return all(_scenario_complete(scenario, completed) for scenario in experiment.scenarios)


def _run_checked(command: list[str], cwd: Path) -> None:
    print("$", " ".join(command))
    subprocess.run(command, cwd=cwd, check=True)


def main() -> int:
    args = parse_args()
    repo_root = REPO_ROOT
    config = load_testing_config(repo_root / "testing" / "config")
    if args.command == "suite":
        return command_suite(repo_root, config, args)
    ctx = build_runtime_context(repo_root, config)
    if args.command == "setup":
        ctx.setup_phase = args.phase

    if args.command == "setup":
        return command_setup(ctx, config)
    if args.command == "validate":
        return command_validate(ctx)
    if args.command == "prepare-dataset":
        return command_prepare_dataset(ctx, config, args)
    if args.command == "prime":
        return command_prime(repo_root, ctx, config, args)
    if args.command == "inject-window":
        return command_inject_window(ctx, config, args)
    if args.command == "run":
        return command_run(ctx, config, args)
    if args.command == "report":
        return command_report(ctx, args)
    raise SystemExit(f"Unsupported command: {args.command}")


if __name__ == "__main__":
    sys.exit(main())
