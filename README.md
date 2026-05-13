# RIoT

[![RIoT Platform: organization](https://img.shields.io/badge/RIoT_Platform-organization-blue?logo=github)](https://github.com/RIoT-Platform)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow)](https://opensource.org/licenses/MIT)

RIoT je informační systém pro zpracování dat v reálném čase, sledování stavů zdrojů dat a vyhodnocování uživatelsky definovaných KPI. Platforma přijímá normalizované vstupy z externích služeb, ukládá historická data do časové databáze a zpřístupňuje aktuální i historické informace přes webovou aplikaci a otevřená API.

Tento repozitář obsahuje jádro systému rozšířené v rámci bakalářské práce **Vojtěcha Hubáčka**. Práce navazuje na původní systém RIoT **Michala Bureše** a na řešení RTAlerts **Dominika Vondrušky**, vůči kterému rozšiřuje práci s historickými daty, API, autorizací, API klíči, exporty a opětovným vyhodnocováním KPI. Výsledkem je obecně použitelná serverová a klientská část, ke které se mohou připojovat různé preprocesory nebo jiné integrační služby posílající data v očekávaném interním formátu. Preprocesory a pomocné služby pro stahování dat jsou vedené odděleně.

![RIoT schema](docs/riot-analytics-scheme.svg)

## Moduly

- [Frontend](frontend/) ([README](frontend/README.md)) - webová administrační a analytická aplikace pro správu zdrojů dat, KPI, API klíčů, historie a exportů.
- [Backend Core](backend/backend-core/) ([README](backend/backend-core/README.md)) - hlavní serverová část, API vrstva, autentizace, správa konfigurace, oprávnění a koordinace ostatních služeb.
- [Message Processing Unit](backend/message-processing-unit/) ([README](backend/message-processing-unit/README.md)) - zpracování příchozích bodů, detekce změn, vyhodnocování KPI a přepočet historických výsledků.
- [Time Series Store](backend/time-series-store/) ([README](backend/time-series-store/README.md)) - zápis a čtení historických surových a KPI dat v InfluxDB.
- [Commons](backend/commons/) ([README](backend/commons/README.md)) - sdílené Go modely, RabbitMQ kontrakty, konstanty a pomocné utility používané backendovými službami.
- [Testing](testing/) ([README](testing/README.md)) - benchmarkový a experimentální testovací nástroj pro ověření ukládání historie, ingestu, reprocessingu a exportů.

Součástí Docker sestavy jsou také podpůrné služby PostgreSQL, InfluxDB, RabbitMQ, Mosquitto, pgAdmin a volitelně monitoring přes Prometheus a Grafanu.

## Spuštění

Pro běžné spuštění stačí Docker, Docker Compose a hlavní `Makefile`.

```bash
make build
```

Příkaz sestaví image a spustí celý stack na pozadí. Pokud už jsou image sestavené, lze použít jen:

```bash
make run
```

Zastavení systému:

```bash
make stop
```

Další cíle, například logy, stav kontejnerů nebo vyčištění lokálních volume, vypíše:

```bash
make help
```

Výchozí konfigurace je v souboru [.env](.env). Pro produkční nebo veřejně dostupné nasazení je potřeba změnit přihlašovací údaje, tajné klíče, OAuth konfiguraci, CORS a nastavení cookies.

## Porty

Po lokálním spuštění jsou hlavní služby dostupné na těchto adresách:

| Služba | Adresa | Poznámka |
| --- | --- | --- |
| Frontend | <http://localhost:8080> | webová aplikace |
| Backend Core | `localhost:9090` | serverová API vrstva |
| RabbitMQ AMQP | `localhost:5672` | interní komunikace služeb |
| RabbitMQ Management | <http://localhost:15672> | správa RabbitMQ, přihlášení podle `.env` |
| PostgreSQL | `localhost:5432` | relační databáze |
| InfluxDB | <http://localhost:8086> | časová databáze |
| Mosquitto MQTT | `localhost:1883` | MQTT broker pro lokální integrace |
| pgAdmin | <http://localhost:8081> | administrační nástroj pro PostgreSQL |

Volitelný monitoring je v compose za profilem `dev`: Grafana na portu `3000`, Prometheus na portu `9091` a MQTT exporter na portu `9000`.

## API

Systém zpřístupňuje API přes GraphQL, REST a WebSocket. Externí klienti se typicky autentizují hlavičkou `X-API-Key`, zatímco webová aplikace používá přihlašovací vrstvu Backend Core.

Příklad REST požadavku:

```bash
curl http://localhost:9090/rest/sd-types \
  -H "X-API-Key: <API_KEY>"
```

Příklad GraphQL požadavku:

```bash
curl http://localhost:9090/graphql \
  -H "Content-Type: application/json" \
  -H "X-API-Key: <API_KEY>" \
  -d '{"query":"{ sdTypes { id uid label } }"}'
```

Příklad WebSocket požadavku:

```bash
wscat -c ws://localhost:9090/ws \
  -H "X-API-Key: <API_KEY>"

> {
>   "type": "request",
>   "id": "sd-types-list",
>   "action": "get_sd_types"
> }
```

Ve výchozím Docker spuštění lze stejné backendové cesty volat také přes frontendovou adresu, kde Nginx proxy směruje požadavky na Backend Core: `http://localhost:8080/graphql`, `http://localhost:8080/rest` a `ws://localhost:8080/ws`.

Další dokumentace API je dostupná přímo ve frontendu v části správy API klíčů `http://localhost:8080/api-keys/docs`, případně v raw podobě v [frontend/src/modules/apiKeys/data/apiDocs.ts](frontend/src/modules/apiKeys/data/apiDocs.ts). Implementační souvislosti backendových rozhraní popisuje také [Backend Core README](backend/backend-core/README.md).

## Napojení Preprocesorů

Preprocesory nejsou součástí tohoto repozitáře. Aby mohly posílat data do běžícího jádra RIoT, musí být připojené do stejné Docker sítě jako služby z `docker-compose.yml`. Při výchozím spuštění přes Docker Compose jde typicky o síť `riot_default`. Přesný název lze ověřit příkazem `docker network ls`.

Z pohledu preprocesoru je důležité, aby z této sítě viděl RabbitMQ pod názvem služby `rabbitmq` a používal AMQP adresu odpovídající konfiguraci:

```text
amqp://riot:secret@rabbitmq:5672
```

Hlavní komunikační kontrakty a názvy front jsou definované ve sdíleném modulu [Commons](backend/commons/README.md). Go preprocesory nebo jiné integrační služby by tento modul měly používat jako sdílenou závislost, aby publikovaly zprávy ve stejném formátu jako jádro systému. Podrobnější popis toku zpráv je v README modulů [Backend Core](backend/backend-core/README.md) a [Message Processing Unit](backend/message-processing-unit/README.md).

Pro běžné napojení zdroje dat jsou podstatné hlavně tyto fronty:

- `sd-type-registration-requests`: registrace typu zdroje dat a jeho parametrů. Zpracovává ji Backend Core.
- `sd-instance-registration-requests`: registrace konkrétních instancí zdroje dat. Zpracovává ji Backend Core.
- `kpi-fulfillment-check-requests`: dávky příchozích stavových zpráv určené k ingestu a vyhodnocení KPI. Zpracovává je Message Processing Unit.

Navazující interní fronty, například `raw-data-point`, `kpi-fulfillment-check-results`, `time-series-raw-data` a `time-series-kpi-results`, už typicky obsluhují backendové moduly RIoT po zpracování vstupu.

## Testování

Automatizované experimenty a benchmarky jsou ve složce [testing](testing/README.md). Nástroj ověřuje hlavní scénáře práce: přípravu dat, online příjem zpráv, ukládání historie, opětovné vyhodnocení KPI a exporty historických dat.

Základní kontrola prostředí:

```bash
python testing/runner.py validate
```

Podrobný přehled příkazů, konfigurace a výstupů je v [testing README](testing/README.md).

## Vývojové Závislosti

Pro běžné spuštění celé platformy stačí Docker a Docker Compose. Při lokálním vývoji jednotlivých částí se používají také:

- Go pro backendové moduly.
- Node.js a npm pro frontend.
- Python pro testovací nástroj.

Přesné příkazy pro lokální vývoj jsou uvedené v README konkrétních modulů.

## Licence

Projekt je licencovaný pod [MIT licencí](LICENSE.txt).
