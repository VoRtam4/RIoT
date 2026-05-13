# RIoT Testing

Benchmark testy určené pro systém RIoT.

## Příprava prostředí

Testovací nástroj je určený pro běh nad spuštěnou instancí RIoT. Před použitím musí být dostupný alespoň Backend Core, RabbitMQ AMQP a RabbitMQ Management API; pro metriky se používá také Prometheus, pokud je nakonfigurovaný v `config/env.yaml`.

Nástroj používá Python 3.10+ a externí balíčky `PyYAML`, `requests` a `pika`. Pokud není v prostředí připravený samostatný dependency soubor, lze je nainstalovat například takto:

```bash
python -m venv .venv
source .venv/bin/activate
pip install PyYAML requests pika
```

Před spuštěním je potřeba zkontrolovat `config/env.yaml`, hlavně adresy backendu, RabbitMQ a platný API klíč.

## Přehled složek a souborů

- `runner.py`: vstupní CLI pro spouštění testů
- `config/`: deklarativní konfigurace zdrojů, datasetů, KPI a experimentů
- `clients/`: klienti pro backend, GraphQL, REST, RabbitMQ a metriky
- `core/`: sdílené modely, načítání konfigurace, reporty, result store
- `datasets/`: příprava a materializace datasetů
- `experiments/`: implementace typů benchmarků
- `generators/`: generátory syntetických zdrojových dat
- `setup/`: registrace typů, instancí a KPI
- `.state/`: lokální stav harnessu
- `outputs/`: raw výsledky a reporty

## Konfigurační soubory

- `config/env.yaml`: endpointy, API klíče, RabbitMQ, timeouty, polling
- `config/sources.yaml`: zdroje, jejich parametry a profily zátěže; v základě jsou definovány zdroje Waze, MHD, NDIC
- `config/datasets.yaml`: datasety DX, délka historie, load profile, build mode
- `config/kpis.yaml`: KPI definice KX pro benchmarky a `extended_kpis`
- `config/experiments.yaml`: experimenty `E1` až `E4`, scénáře HX, IX, SX, RX, počet opakování

## Spouštění a přepínače

### `validate`

- `python testing/runner.py validate`
- bez přepínačů, jen ověření prostředí

### `setup`

- `python testing/runner.py setup`
- `--phase all`: celé nastavení
- `--phase types`: jen registrace SD types
- `--phase refresh`: refresh instancí ze systému
- `--phase kpis`: vytvoření nebo navázání KPI

### `prepare-dataset`

- `python testing/runner.py prepare-dataset --dataset D3`
- `--dataset DX`: zvolený dataset
- `--rebuild`: vynutí nové vytvoření
- `--resume`: naváže na rozpracovaný build

### `prime`

- `python testing/runner.py prime --experiment E3`
- `--dataset DX`: připraví konkrétní dataset
- `--experiment EX`: odvodí dataset z experimentu
- `--scenario HX`: odvodí dataset ze scénáře
- `--all`: připraví vše potřebné pro zapnuté experimenty
- `--rebuild`: znovu vytvoří dataset
- `--resume`: naváže na předchozí přípravu
- `--full-reset`: před primingem resetuje runtime stack
- `--api-key ...`: použije zadaný API key
- `--validate-timeout-seconds N`: timeout validace
- `--validate-poll-seconds N`: interval pollingu validace
- `--skip-down-v`: přeskočí `docker compose down -v`
- `--skip-clean-docker-dir`: nemaže runtime docker adresář
- `--skip-prune`: přeskočí docker prune
- `--skip-build`: nepřestaví image
- `--keep-outputs`: nemaže `testing/outputs`

### `inject-window`

- `python testing/runner.py inject-window --source src_name --minutes 15`
- `--source src_name`: zdroj
- `--minutes N`: délka simulovaného okna
- `--profile normal|high|stress`: profil zátěže
- `--bootstrap`: vytvoří i potřebné instance
- `--dry-run`: jen vypíše plán

### `run`

- `python testing/runner.py run --experiment EX`
- `--experiment EX`: spustí celý experiment
- `--scenario RX`: spustí jeden scénář
- `--all`: spustí vše, co je zapnuté v konfiguraci
- `--dry-run`: jen vypíše plán běhu

### `report`

- `python testing/runner.py report --format md`
- `--experiment EX`: filtr jen na jeden experiment
- `--format console|md|csv|json`: formát výstupu

### `suite`

- `python testing/runner.py suite`
- `--dataset DX`: výchozí dataset pro části suite, které ho potřebují
- `--api-key ...`: nepoužije interaktivní dotaz
- `--dry-run`: vypíše plán suite
- `--resume`: přeskočí již hotové scénáře
- `--validate-timeout-seconds N`: timeout validace
- `--validate-poll-seconds N`: interval pollingu validace
- `--experiments EX EX EX EX`: vlastní pořadí nebo podmnožina experimentů
- `--skip-down-v`: přeskočí `docker compose down -v`
- `--skip-clean-docker-dir`: nemaže runtime docker adresář
- `--skip-prune`: přeskočí docker prune
- `--skip-build`: nepřestaví image
- `--keep-outputs`: nemaže staré výstupy

## Outputy

- `testing/.state/`: runtime stav, registry datasetů, uložené identifikátory
- `testing/outputs/raw/`: raw výsledky jednotlivých běhů
- `testing/outputs/reports/`: `suite_summary.*`, `latest_summary.*`, další přehledy

## Příklady použití

```bash
# kontrola prostředí
python testing/runner.py validate

# příprava typů a KPI
python testing/runner.py setup --phase types
python testing/runner.py setup --phase kpis

# příprava datasetu
python testing/runner.py prepare-dataset --dataset DX --rebuild

# spuštění jednoho experimentu
python testing/runner.py run --experiment EX

# spuštění jednoho scénáře
python testing/runner.py run --scenario RX

# krátké online okno dat
python testing/runner.py inject-window --source src_name --minutes 15 --profile high

# celý orchestrated běh
python testing/runner.py suite --experiments E1 E2 E3 E4

# jen plán suite
python testing/runner.py suite --dry-run

# export reportu
python testing/runner.py report --format md
```
