# Time Series Store

Time Series Store je interní modul pro ukládání a čtení časových dat v InfluxDB. Přijímá surové body a výsledky KPI z Message Processing Unit, zapisuje je do odpovídajících měření a poskytuje ostatním backendovým částem streamované čtení nad časovou historií.

Modul řeší běžné čtení dat pro API a exporty, zjišťování dostupných hodnot tagů, agregované čtení KPI v časových oknech a čtení surových bodů pro přepočet KPI. S okolím komunikuje přes RabbitMQ, zatímco fyzické uložení dat probíhá v InfluxDB.

## Spuštění

Modul se běžně spouští z kořene projektu přes hlavní `docker-compose.yml`:

```bash
docker compose up -d riot-time-series-store
```

Služba závisí na RabbitMQ, Backend Core a InfluxDB. Backend Core připravuje RabbitMQ infrastrukturu, RabbitMQ zprostředkovává požadavky a InfluxDB slouží jako vlastní úložiště časových řad. Modul nemá vlastní veřejný HTTP port; port `8086` patří přidružené InfluxDB službě.

## Konfigurace

Hlavní proměnné prostředí používané modulem:

- `BACKEND_CORE_URL`: adresa Backend Core používaná při čekání na dostupnost systému
- `RABBITMQ_URL`: připojovací řetězec k RabbitMQ
- `INFLUX_URL`: adresa InfluxDB
- `INFLUX_TOKEN`: token pro přístup do InfluxDB
- `INFLUX_ORGANIZATION`: organizace v InfluxDB
- `INFLUX_BUCKET`: bucket pro ukládání časových dat
- `INFLUX_READ_WORKERS`: počet workerů pro běžné čtecí požadavky
- `INFLUX_REPROCESS_WORKERS`: počet workerů pro čtení historických dat při přepočtu KPI
- `DEV_ENV`: vývojové nastavení předávané z hlavního compose

## RabbitMQ komunikace

Typické přijímané fronty a zprávy:

- `time-series-raw-data`: dávky surových bodů k uložení
- `time-series-kpi-results`: dávky KPI výsledků k uložení
- `time-series-read-request`: požadavky na čtení časových dat pro API a exporty
- `time-series-read-cancel-request`: zrušení běžícího čtecího požadavku
- `time-series-distinct-tag-values-request`: požadavky na dostupné hodnoty konkrétního tagu
- `time-series-reprocess-read-request`: požadavky na historická surová data pro přepočet KPI
- `time-series-delete-requests-tsdb`: požadavky na mazání KPI výsledků v časové databázi

Typické odesílané odpovědi:

- `time-series-read-response`: dávkové odpovědi na běžné čtecí požadavky
- `time-series-distinct-tag-values-response`: odpovědi s dostupnými hodnotami tagů
- `time-series-reprocess-read-response`: dávkové odpovědi se surovými body pro reprocessing

## Struktura modulu

- `Dockerfile`: build Go aplikace společně s lokálním modulem `commons`
- `go.mod`, `go.sum`: Go modul a závislosti
- `src/main.go`: inicializace modulu, kontrola závislostí a spuštění RabbitMQ konzumentů
- `src/internal/client.go`: vytvoření InfluxDB klienta
- `src/internal/write.go`: zápis surových a KPI bodů, mazání KPI výsledků
- `src/internal/read.go`: streamované čtení časových dat
- `src/internal/reprocess.go`: čtení historických surových dat pro přepočet KPI
- `src/internal/aggregation.go`: agregované čtení KPI v časových oknech
- `src/internal/filters.go`: převod filtrů do dotazů nad InfluxDB
- `src/internal/fluxBuilder.go`: skládání Flux dotazů
- `src/internal/planner.go`: příprava plánu čtení podle požadavku
- `src/internal/chunks.go`: dělení dlouhých časových rozsahů na menší okna
- `src/internal/readJobs.go`: evidence a rušení běžících čtecích požadavků
- `src/internal/types.go`: interní typy modulu
