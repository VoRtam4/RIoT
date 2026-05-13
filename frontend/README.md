# Frontend

Frontend je webová administrační a analytická aplikace RIoT postavená nad Reactem, Vite, TypeScriptem, Apollo Clientem a Material UI. Umožňuje správu API klíčů, typů a instancí zdrojů dat, KPI definic, časové historie a exportů. Součástí aplikace je také dokumentace API pro externí klienty.

Aplikace komunikuje s Backend Core přes GraphQL HTTP, GraphQL WebSocket a REST. V produkčním Docker buildu jsou tyto cesty proxované přes Nginx na backend, v lokálním vývoji je možné použít přímé Vite proměnné s adresami API.

## Spuštění

V hlavní sestavě systému se frontend spouští z kořene projektu:

```bash
docker compose up -d riot-frontend
```

Služba je publikovaná na portu `8080`. Z hlediska kontejnerů závisí na `riot-backend-core`, protože Nginx proxy směruje backendové cesty na adresu předanou v `BACKEND_CORE_URL`. Docker build používá Node.js 20 a výslednou statickou aplikaci obsluhuje Nginx.

Lokální vývoj bez Dockeru:

```bash
npm install
npm run dev
```

Vite vývojový server běží typicky na portu `5173` a pro komunikaci s backendem používá hodnoty z proměnných `VITE_*`, případně relativní cesty proxované až v produkční Docker sestavě.

Produkční build:

```bash
npm run build
```

## Konfigurace

Proměnné používané při Docker buildu a běhu:

- `BACKEND_CORE_URL`: adresa Backend Core pro Nginx proxy

Proměnné používané ve Vite vývojovém režimu:

- `VITE_API_ORIGIN`: společný základ pro GraphQL, REST a WebSocket endpointy
- `VITE_GRAPHQL_URL`: explicitní adresa GraphQL HTTP endpointu
- `VITE_GRAPHQL_WS_URL`: explicitní adresa GraphQL WebSocket endpointu
- `VITE_REST_URL`: explicitní adresa REST API
- `VITE_WS_URL`: explicitní adresa WebSocket API

Pokud nejsou Vite proměnné nastavené, frontend používá relativní cesty `/graphql`, `/rest` a `/ws`. V Docker buildu je obslouží Nginx proxy.

## Rozhraní na backend

Frontend používá tato backendová rozhraní:

- GraphQL HTTP: běžné dotazy a mutace
- GraphQL WebSocket: subscription operace
- REST: exporty časových dat

Adresy endpointů se skládají v `src/app/apiEndpoints.ts`. Nginx proxy pro produkční build je definovaná v `nginx.conf`.

## Hlavní části aplikace

- dashboard a základní přehled systému
- správa API klíčů a dokumentace API
- správa KPI definic a detail KPI
- přehled instancí zdrojů dat a jejich detail
- prohlížení časové historie, agregací a exportů
- autentizace přes Backend Core

## GraphQL codegen

GraphQL typy a dokumenty se generují podle `codegen.yml`:

```bash
npm run codegen
```

Codegen očekává dostupné GraphQL schéma na `http://localhost:9090/graphql`. Výstupy se zapisují do `src/generated/` a introspekční schema do `graphql.schema.json`.

## Struktura modulu

- `Dockerfile`: build statické aplikace a výsledný Nginx image
- `nginx.conf`: proxy konfigurace pro Backend Core a fallback na SPA routy
- `package.json`, `package-lock.json`: Node skripty a závislosti
- `vite.config.ts`: konfigurace Vite
- `codegen.yml`: konfigurace GraphQL codegenu
- `graphql.schema.json`: introspekční GraphQL schéma
- `src/main.tsx`: inicializace React aplikace, providerů a routeru
- `src/app/`: router, Apollo Client, chráněné routy a skládání API endpointů
- `src/components/`: sdílené UI komponenty a layout
- `src/modules/`: doménové moduly aplikace
- `src/pages/`: routované stránky aplikace
- `src/generated/`: generované GraphQL typy a dokumenty
- `src/theme/`, `src/styles/`: téma, globální styly a inicializace vzhledu
