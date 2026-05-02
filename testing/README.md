# RIoT Testing

This directory contains a configurable testing harness for RIoT thesis benchmarks.

The goal is to keep experiment definitions declarative and reusable.

Testing is split into two modes on purpose:

- `E1 ingest` uses the real online path through RabbitMQ and MPU
- `E2/E3/E4` use a fast historical dataset seed for TSDB-oriented experiments

This avoids measuring slow historical dataset creation as if it were online ingest throughput.

Configuration lives in:

- source simulation profiles live in `config/sources.yaml`
- KPI definitions live in `config/kpis.yaml`
- dataset build definitions live in `config/datasets.yaml`
- experiment scenarios live in `config/experiments.yaml`

The harness is intentionally split into small modules:

- `runner.py` orchestrates commands
- `core/` holds shared models, config loading, reporting and state
- `clients/` contains HTTP, GraphQL, REST and metrics clients
- `setup/` prepares SD types, instances and KPI definitions
- `datasets/` builds and registers datasets
- `experiments/` implements benchmark types
- `generators/` contains source generator abstractions

## Current state

The scaffold is implemented and ready for extension. The following parts already work:

- configuration loading
- local state and dataset registry persistence
- experiment resolution from config
- `setup` command flow
- `prepare-dataset` command flow
- `prime` command flow
- `run --dry-run`
- result serialization and markdown/csv summaries
- fast-forward source simulation for `WAZE`, `MHD` and `NDIC`
- realistic per-source lifecycle modelling for dataset sizing and ingest baselines
- automatic checkpoint materialization for smaller datasets when building a larger one
- direct synthetic ISC injection for `WAZE` and `NDIC` via RabbitMQ management API
- direct synthetic ISC injection for `MHD` via RabbitMQ management API
- real queue-idle based execution path for `E1` ingest and `E2` reprocess when the stack is running

The built-in generators currently model:

- `WAZE`: minute ticks, lower night activity, roughly 150-200 active segments per minute, 3-15 minute active lifetimes, recurring segment returns
- `MHD`: minute ticks, near-zero traffic between `00:30-04:30`, 30-50 minute trip lifetimes, smooth delay evolution
- `NDIC`: sparse minute ticks, low arrival rate, 1-2 hour active lifetimes, slower churn

All generators run in `fast-forward` mode. A one-month dataset therefore simulates one month of timestamps, but it does not wait one month of wall-clock time.

Dataset preparation now materializes compatible checkpoints automatically. For example, building `D3` also stores dataset registry entries for `D1` and `D2` if they share the same source set and load profile.

The following parts still need project-specific API payload completion:

- final reprocess trigger strategy for repeated benchmarks
- `MHD` direct injection alignment and optional real preprocessor-facing adapters
- optional metrics polling from RabbitMQ / Prometheus

Important constraint:

- `SELECTED` KPI definitions cannot be created before the system has created real `SDInstance`
  records for the given source. In practice this means:
  1. run `setup` to create `SDType`s,
  2. ingest source data until instances exist,
  3. run `setup` again to bind and create `SELECTED` KPI definitions.

Recommended setup flow in the current harness:

1. `python testing/runner.py setup --phase types`
2. prepare or ingest data so real `SDInstance` records are created
3. `python testing/runner.py setup --phase refresh`
4. `python testing/runner.py setup --phase kpis`

Recommended full workflow now:

1. `python testing/runner.py validate`
2. `python testing/runner.py prime --experiment E3 --rebuild`
3. `python testing/runner.py run --experiment E3`
4. `python testing/runner.py run --experiment E2`
5. `python testing/runner.py run --experiment E1`
6. `python testing/runner.py run --experiment E4`

## Commands

```bash
python testing/runner.py setup
python testing/runner.py validate

python testing/runner.py prepare-dataset --dataset D1
python testing/runner.py prepare-dataset --dataset D2 --resume
python testing/runner.py prime --experiment E3 --rebuild
python testing/runner.py prime --dataset D3 --rebuild

python testing/runner.py inject-window --source waze --minutes 10 --profile normal
python testing/runner.py inject-window --source ndic --minutes 15 --profile high
python testing/runner.py inject-window --source mhd --minutes 10 --profile normal

python testing/runner.py run --all --dry-run
python testing/runner.py run --experiment E2
python testing/runner.py run --scenario H4

python testing/runner.py report --format md
```

## Outputs

- raw run outputs: `testing/outputs/raw/`
- generated summaries: `testing/outputs/reports/`
- local ids and dataset registry: `testing/.state/`

## Recommended next steps

1. Run the stack and verify the live `inject-window`, `E1` and `E2` paths end-to-end.
2. Tighten KPI point counting and queue-based completion heuristics based on live measurements.
3. Add project-specific metrics polling in `clients/metrics_client.py`.
4. Replace the remaining estimated parts of `E4` with live storage measurements.
