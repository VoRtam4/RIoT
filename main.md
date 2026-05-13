# RIoT

[![RIoT Platform: organization](https://img.shields.io/badge/RIoT_Platform-organization-blue?logo=github)](https://github.com/RIoT-Platform)

Autor práce: **Vojtěch Hubáček**

Tento soubor slouží jako úvodní README k odevzdávané složce bakalářské práce. Práce rozšiřuje platformu RIoT pro práci s dopravními daty v reálném čase, jejich historizaci, opětovné vyhodnocování KPI a zpřístupnění přes autorizovaná aplikační rozhraní. Výsledkem je obecně použitelný systém, který lze využít samostatně i jako serverovou část pro další aplikace.

Řešení navazuje na původní platformu RIoT **Michala Bureše** a na systém RTAlerts **Dominika Vondrušky**, z něhož vychází dopravní aplikační doména a práce s konkrétními zdroji otevřených dat. Podrobné technické popisy, návody ke spuštění a dokumentace jednotlivých modulů jsou uvedené přímo v README souborech odpovídajících částí.

## Obsah Odevzdávané Složky

Předpokládaná struktura předávané složky je:

```text
.
├── RIoT-Platform/
├── RIoT/
├── docs/
└── README.md
```

## RIoT

Složka [RIoT-Platform/RIoT/](RIoT-Platform/RIoT/) obsahuje hlavní jádro systému. Jde o samostatně spustitelnou a zdrojově nezávislou část platformy: webovou aplikaci, backendové služby, časovou vrstvu, KPI zpracování, API a testovací nástroj. RIoT přijímá normalizované vstupy od libovolných preprocesorů nebo integračních služeb, ale na konkrétním typu zdroje není přímo závislý.

Hlavní dokumentace je v [RIoT-Platform/RIoT/README.md](RIoT-Platform/RIoT/README.md). Obsah repozitáře zahrnuje zejména:

- frontend pro správu zdrojů dat, KPI, API klíčů, historie a exportů,
- Backend Core, Message Processing Unit a Time Series Store,
- sdílený modul `commons`,
- benchmarkový a experimentální nástroj ve složce `testing`.

RIoT lze spustit samostatně. Pro praktické ověření celé dopravní sestavy je však vhodnější použít spuštění přes repozitář preprocesorů, který z důvodu své závislosti na RIoT zajistí start jádra i navazujících integračních služeb v jednom kroku.

## RIoT-Preprocessors

Složka [RIoT-Platform/RIoT-Preprocessors/](RIoT-Platform/RIoT-Preprocessors/) obsahuje zdrojově specifickou integrační vrstvu pro reálná dopravní data. Na rozdíl od obecného jádra RIoT zde vznikají konkrétní převody vstupních zdrojů do interního datového modelu platformy.

Hlavní dokumentace je v [RIoT-Platform/RIoT-Preprocessors/README.md](RIoT-Platform/RIoT-Preprocessors/README.md). Repozitář obsahuje:

- Waze Jam preprocesor, NDIC preprocesor, MHD preprocesor,
- pomocnou službu DATEX Downloader pro příjem NDIC zpráv, jejíž hostovaná instance je na <https://riot-preprocessors.onrender.com/>, při nedostupnosti je potřeba službu zveřejnit jinde nebo kontaktovat autora,
- sdílený modul `commons`.

Tento repozitář je zároveň doporučeným výchozím bodem pro spuštění celé dopravní varianty systému, protože jeho `Makefile` pracuje s cestou k repozitáři RIoT a umí spustit oba stacky společně.

## RIoT-Apps

Repozitář RIoT-Apps slouží k integračnímu ověření RIoT nad dvěma existujícími aplikacemi: analytickou aplikací Analyticity **Magdalény Ondruškové** nad Waze daty a aplikací Lissy **Juraje Lazúra** pro analýzu veřejné dopravy. Upravené varianty využívají RIoT jako zdroj historických dat a KPI výsledků, aby bylo možné ověřit jeho použitelnost v běžném provozu s reálnými daty a s již existujícími analytickými nebo vizualizačními nástroji.

Tyto aplikace nejsou součástí odevzdávané složky. Po konzultaci a vyhodnocení rozsahu byly vedené odděleně, protože nejde o moduly jádra RIoT ani o preprocesory, ale o cizí aplikace s převážně původním kódem, do nichž byly doplněny cílené integrační úpravy. Samostatný repozitář s aplikacemi připravenými ke spuštění nad běžícím RIoT je pro ukázku dostupný zde:

> [https://github.com/VoRtam4/RIoT-Apps](https://github.com/VoRtam4/RIoT-Apps)

## docs

Složka [docs/](docs/) obsahuje dokumentační část odevzdání:

- výsledný text práce [docs/xhubacv00_BP.pdf](docs/xhubacv00_BP.pdf),
- zdrojové soubory práce [docs/xhubacv00_BP/](docs/xhubacv00_BP/).
