# Message Processing Unit

Message Processing Unit zpracovává příchozí data ze zdrojů a převádí je na interní surové body a výsledky KPI. U každé instance sleduje poslední známý stav, vyhodnocuje změny hodnot a podle aktuální konfigurace KPI rozhoduje, které metriky jsou splněné. Výsledky předává dále Backend Core a časové vrstvě.

Modul zároveň podporuje přepočet KPI nad historickými daty. V takovém případě obdrží požadavek od Backend Core, vyžádá si potřebné surové body z časové vrstvy, znovu nad nimi vyhodnotí příslušnou KPI definici a odešle nově vypočtené výsledky.

S okolím komunikuje přes RabbitMQ. Backend Core modulu předává konfiguraci typů zdrojů, KPI definic a požadavky na přepočet, zatímco časová vrstva přijímá surové a KPI body a vrací data potřebná pro reprocessing.

## Spuštění

Modul se běžně spouští z kořene projektu přes hlavní `docker-compose.yml`:

```bash
docker compose up -d riot-message-processing-unit
```

Služba závisí na RabbitMQ a Backend Core. Backend Core připravuje RabbitMQ infrastrukturu a posílá počáteční konfiguraci, podle které MPU rozlišuje typy zdrojů, parametry a KPI definice. Modul nemá vlastní veřejný HTTP port, komunikuje pouze přes RabbitMQ. Počet instancí lze při Docker Compose spuštění navyšovat například přes `--scale riot-message-processing-unit=N`.

## Konfigurace

Hlavní proměnné prostředí používané modulem:

- `BACKEND_CORE_URL`: adresa Backend Core používaná při čekání na dostupnost konfiguračního backendu
- `RABBITMQ_URL`: připojovací řetězec k RabbitMQ
- `MPU_INPUT_WORKERS`: počet workerů pro běžné vstupní zpracování
- `MPU_INTERNAL_INPUT_BUFFER`: kapacita interní vstupní fronty mezi RabbitMQ konzumenty a processing workery
- `MPU_INPUT_PAUSE_HIGH_WATERMARK_PERCENT`: celkové zaplnění interní fronty, od kterého MPU dočasně pozastaví příjem ze všech vstupních front
- `MPU_INPUT_RESUME_LOW_WATERMARK_PERCENT`: celkové zaplnění interní fronty, pod kterým MPU znovu obnoví globálně pozastavený příjem
- `MPU_INPUT_SOURCE_PAUSE_HIGH_WATERMARK_PERCENT`: podíl kapacity interní fronty obsazený jedním zdrojem, od kterého MPU dočasně pozastaví příjem jen z dané fronty
- `MPU_INPUT_SOURCE_RESUME_LOW_WATERMARK_PERCENT`: podíl kapacity interní fronty obsazený jedním zdrojem, pod kterým MPU obnoví příjem z dané fronty
- `MPU_INPUT_PAUSE_CHECK_INTERVAL_MS`: interval kontroly zaplnění interní fronty během pauzy příjmu
- `MPU_INPUT_TREND_SAMPLE_INTERVAL_MS`: interval výpočtu trendů příjmu, zpracování a backlogu pro budoucí scheduler rozhodování
- `MPU_INPUT_WATCH_PRESSURE_THRESHOLD`: krátkodobý rozdíl `recv/s - proc/s`, od kterého se zdroj označí jako sledovaný
- `MPU_INPUT_DROP_CANDIDATE_PRESSURE_THRESHOLD`: dlouhodobý rozdíl `recv/s - proc/s`, který zdroj po hysteresi označí jako kandidáta pro budoucí dropování
- `MPU_INPUT_DROP_CANDIDATE_BACKLOG_DELTA_THRESHOLD`: dlouhodobý růst backlogu vyžadovaný pro stav kandidáta dropování
- `MPU_INPUT_STATE_ESCALATION_SAMPLES`: počet po sobě jdoucích trend vzorků potřebných pro eskalaci scheduler stavu
- `MPU_INPUT_STATE_RECOVERY_SAMPLES`: počet po sobě jdoucích zotavovacích vzorků potřebných pro návrat do normálu
- `MPU_INPUT_DROPPER_ENABLED`: zapnutí skutečného zahazování zpráv ze zdrojů ve stavu `dropCandidate`; výchozí hodnota je `false`, tedy pouze dry-run
- `MPU_INPUT_DROPPER_GRACE_SAMPLES`: počet dalších trend vzorků ve stavu `dropCandidate` před tím, než je zdroj považovaný za připravený k dropování
- `MPU_REPROCESS_WORKERS`: počet workerů pro přepočet KPI
- `DEV_ENV`: vývojové nastavení předávané z hlavního compose

## RabbitMQ komunikace

Typické přijímané fronty a zprávy:

- `kpi-fulfillment-check-requests`: dávky nových dat určených ke zpracování a vyhodnocení KPI
- `ingest.<sdTypeUID>`: per-type fronty pro nový ingest tok, navázané na direct exchange `ingest` routing key `<sdTypeUID>`
- `set-of-sd-types-updates`: aktuální konfigurace typů zdrojů dat
- `kpi-config-update`: konfigurace KPI přijímaná přes fanout exchange
- `raw-data-point-cache-bootstrap`: počáteční naplnění cache posledních surových hodnot
- `kpi-fulfillment-cache-bootstrap`: počáteční naplnění cache posledních KPI hodnot
- `kpi-reprocess-requests`: požadavky na přepočet KPI nad historií
- `time-series-delete-requests-mpu`: požadavky související s mazáním KPI výsledků a rušením aktivního přepočtu
- `time-series-reprocess-read-response`: odpovědi časové vrstvy se surovými body pro reprocessing

Typické odesílané fronty a zprávy:

- `message-processing-unit-connection-notifications`: oznámení Backend Core, že MPU běží
- `sd-type-registration-requests`: registrace nově zjištěných parametrů typu zdroje dat
- `raw-data-point`: změněné surové body pro Backend Core
- `kpi-fulfillment-check-results`: výsledky KPI pro Backend Core
- `time-series-raw-data`: surové body pro uložení do časové vrstvy
- `time-series-kpi-results`: KPI body pro uložení do časové vrstvy
- `time-series-reprocess-read-request`: požadavky na čtení historických dat pro přepočet KPI

## Struktura modulu

- `Dockerfile`: build Go aplikace společně s lokálním modulem `commons`
- `go.mod`, `go.sum`: Go modul a závislosti
- `src/main.go`: inicializace modulu, čekání na Backend Core a spuštění konzumentů RabbitMQ
- `src/dispatch/`: příjem vstupních zpráv, per-type ingest fronty, interní input buffer, statistiky a scheduler stav
- `src/processing/bootstrap.go`: bootstrap runtime cache a oznámení dostupnosti modulu
- `src/processing/processingState.go`: sdílený stav konfigurace KPI, typů zdrojů a aktivních jobů
- `src/processing/tupleProcessing.go`: dávkové zpracování vstupních zpráv
- `src/processing/rawProcessing.go`: normalizace a detekce změn surových hodnot
- `src/processing/kpiProcessing.go`: vyhodnocení KPI nad příchozími daty
- `src/processing/kpiFulfillmentCheck.go`: implementace vyhodnocování stromu KPI podmínek
- `src/processing/reprocess.go`: přepočet KPI nad historickými daty
- `src/processing/reprocessSync.go`: synchronizace přepočtu s aktuální konfigurací a mazáním
- `src/processing/publishing.go`: publikování výsledků do navazujících front
- `src/processing/processingCache.go`: runtime cache posledních surových a KPI hodnot
