# Backend Core

Backend Core je hlavní serverový modul platformy RIoT. Drží konfiguraci sledovaných zdrojů dat, instancí, KPI definic, uživatelských rolí a API klíčů. Zároveň poskytuje veřejná API pro analytické aplikace, předává konfiguraci dalším backendovým částem přes RabbitMQ a zprostředkovává dotazy nad časovou vrstvou. Modul sám neukládá časové řady, ale koordinuje jejich čtení, exporty, mazání a přepočty KPI.

## Spuštění

Modul se běžně spouští z kořene projektu přes hlavní `docker-compose.yml`:

```bash
docker compose up -d riot-backend-core
```

Služba je v compose publikována na portu `9090` a závisí hlavně na PostgreSQL, RabbitMQ a časové vrstvě systému. Při lokálním vývoji je praktické spouštět ji společně s navazujícími službami, protože část funkcionality REST, GraphQL a WebSocket rozhraní očekává dostupnou databázi, fronty a odpovědi z Time Series Store. V plné Docker sestavě je stejné API dostupné také přes frontendovou proxy na portu `8080`.

## Konfigurace

Hlavní proměnné prostředí používané modulem:

- `DEV_ENV`: režim běhu aplikace, developmet hodnota `true` nebo produkční hodnota `false`
- `POSTGRES_URL`: připojovací řetězec k PostgreSQL
- `RABBITMQ_URL`: připojovací řetězec k RabbitMQ
- `ALLOWED_ORIGINS`: povolené originy pro CORS
- `JWT_SECRET`: tajemství pro podepisování JWT
- `JWT_AUTHENTICATION_MIDDLEWARE_ENABLED`: zapnutí nebo vypnutí JWT middleware
- `GOOGLE_OAUTH2_CLIENT_ID`: klientský identifikátor pro Google OAuth
- `GOOGLE_OAUTH2_CLIENT_SECRET`: klientské tajemství pro Google OAuth
- `AUTH_REDIRECT_URL`: návratová adresa po přihlášení
- `ROOT_ADMIN_EMAIL`: e-mail uživatele s výchozí administrátorskou rolí
- `SECURE_COOKIES`: nastavení zabezpečených cookies pro HTTPS prostředí

## Rozhraní

Základní adresy při lokálním spuštění:

- REST API: `http://localhost:9090/rest`
- GraphQL HTTP API: `http://localhost:9090/graphql`
- GraphQL WebSocket: `ws://localhost:9090/graphql`
- WebSocket API: `ws://localhost:9090/ws`
- OAuth přihlášení: `http://localhost:9090/auth/login`
- OAuth odhlášení: `http://localhost:9090/auth/logout`

REST, GraphQL a WebSocket rozhraní z velké části zpřístupňují stejnou doménovou funkcionalitu systému. REST je vhodné pro přímé dotazy, správu entit a stahování exportů. GraphQL je vhodné pro klienty, kteří potřebují vybrat jen konkrétní část odpovědi a tím omezit velikost přenášených dat. WebSocket rozhraní se používá pro postupné streamování dat a dlouhé operace, kde není vhodné čekat na jednu velkou odpověď.

Autorizace probíhá přes přihlašovací část v `auth/`, typicky pomocí session nebo JWT, případně pomocí hlavičky `X-API-Key` pro externí aplikace. Stejnou hlavičku používají REST, GraphQL i WebSocket klienti.

Příklad REST požadavku s API klíčem:

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

Příklad připojení k WebSocket API:

```bash
wscat -c ws://localhost:9090/ws \
  -H "X-API-Key: <API_KEY>"
```

Podrobnější příklady volání REST, GraphQL a WebSocket rozhraní jsou součástí dokumentace API ve frontendu, hlavně v `frontend/src/modules/apiKeys/data/apiDocs.ts` a `frontend/src/modules/apiKeys/utils/apiDocsExamples.ts`. Frontendové adresy jednotlivých rozhraní se skládají v `frontend/src/app/apiEndpoints.ts`.

## RabbitMQ komunikace

Backend Core používá RabbitMQ jako vnitřní komunikační vrstvu mezi konfiguračním backendem, preprocesory, Message Processing Unit a časovou vrstvou.

Typické přijímané fronty:

- `sd-type-registration-requests`: registrace typu zdroje dat
- `sd-instance-registration-requests`: registrace konkrétní instance zdroje dat
- `raw-data-point`: příjem surových bodů z preprocesoru nebo MPU
- `kpi-fulfillment-check-results`: výsledky ověření splnění KPI
- `message-processing-unit-connection-notifications`: oznámení o připojení MPU
- `time-series-read-response`: odpovědi na čtení z časové vrstvy
- `time-series-distinct-tag-values-response`: odpovědi s dostupnými hodnotami tagů

Typické odesílané fronty a zprávy:

- `set-of-sd-types-updates`: aktuální konfigurace typů zdrojů dat
- `set-of-sd-instances-updates`: aktuální konfigurace instancí zdrojů dat
- `kpi-config-update`: aktuální konfigurace KPI
- `raw-data-point-cache-bootstrap`: inicializace cache surových bodů
- `kpi-fulfillment-cache-bootstrap`: inicializace cache KPI vyhodnocení
- `kpi-reprocess-requests`: požadavky na přepočtení KPI
- `time-series-read-request`: požadavky na čtení z časové vrstvy
- `time-series-distinct-tag-values-request`: požadavky na dostupné hodnoty tagů
- `time-series-read-cancel-request`: zrušení dlouhého čtení nebo exportu
- `time-series-delete-requests-tsdb`: mazání dat v časové databázi
- `time-series-delete-requests-mpu`: mazání navazujících dat v MPU

## Struktura modulu

- `Dockerfile`: build Go aplikace a výsledného runtime image
- `go.mod`, `go.sum`: Go modul a závislosti
- `gqlgen.yml`: konfigurace generování GraphQL serveru
- `gqlgen-generate.bat`: pomocný skript pro GraphQL codegen
- `tools.go`: pomocné build-time závislosti
- `src/main.go`: inicializace modulu, připojení k databázi a RabbitMQ, spuštění API serveru
- `src/api/`: HTTP server, REST, GraphQL a WebSocket vrstva
- `src/auth/`: OAuth, JWT, cookies, API klíče, role a oprávnění
- `src/db/`: relační databázová vrstva
- `src/domainLogicLayer/`: aplikační logika nad entitami, KPI a historií
- `src/events/`: interní sběrnice událostí pro distribuci změn klientům
- `src/isc/`: interní servisní komunikace přes RabbitMQ
- `src/model/`: databázové a doménové modely
- `src/modelMapping/`: převody mezi interními modely a API modely

## Poznámky k vývoji

GraphQL kód je generovaný podle `gqlgen.yml`. Při změně GraphQL schématu je potřeba přegenerovat odpovídající resolverovou a modelovou vrstvu:

```bash
go run github.com/99designs/gqlgen generate
```

Ve vývojovém prostředí `DEV_ENV=true` lze GraphQL schéma použít pro generování klientských typů a dotazů externích aplikací přímo z běžícího API. Modul povoluje introspekční GraphQL dotazy, které k tomu klienti a codegen nástroje typicky používají.

Při změně autentizace nebo CORS je nutné kontrolovat navazující frontendové nastavení, hlavně `ALLOWED_ORIGINS`, cookies, JWT a použití hlavičky `X-API-Key`.
