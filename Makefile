# @file Makefile
# @brief Hlavní Makefile pro lokální správu Docker Compose stacku platformy RIoT.
#
# @author Vojtěch Hubáček
#
# @par Autorský podíl
# - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
#
# @defgroup riot_root Root
# @ingroup riot
# @see README.md

COMPOSE ?= docker compose
DOCKER_DIR ?= docker
LOG_TAIL ?= 30
ARGS ?= -h

ifeq ($(OS),Windows_NT)
	RM_DOCKER_DIR = if exist "$(DOCKER_DIR)" rmdir /S /Q "$(DOCKER_DIR)"
else
	RM_DOCKER_DIR = rm -rf "$(DOCKER_DIR)"
endif

.PHONY: help build run stop clear restart reset prune clear-all reset-all status logs test

help:
	@echo "Stack:"
	@echo "  make build              Start stack with image rebuild"
	@echo "  make run                Start stack"
	@echo "  make stop               Stop stack"
	@echo "  make restart            Stop stack and start it again"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clear              Stop stack, remove volumes and delete ./$(DOCKER_DIR)"
	@echo "  make reset              Run clear and then build"
	@echo ""
	@echo "Aggressive cleanup:"
	@echo "  make prune              Stop stack, remove volumes, images and orphan containers"
	@echo "  make clear-all          Run prune and delete ./$(DOCKER_DIR)"
	@echo "  make reset-all          Run clear-all and then build"
	@echo ""
	@echo "Diagnostics:"
	@echo "  make status             Show stack containers"
	@echo "  make logs               Show stack logs"
	@echo "  make logs LOG_TAIL=100  Show last 100 log lines"

build:
	$(COMPOSE) up --build -d

run:
	$(COMPOSE) up -d

stop:
	$(COMPOSE) down

clear:
	-$(COMPOSE) down -v
	@$(RM_DOCKER_DIR)

restart: stop run

reset: clear build

prune:
	-$(COMPOSE) down -v --rmi local --remove-orphans

clear-all: prune
	@$(RM_DOCKER_DIR)

reset-all: clear-all build

status:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs --tail=$(LOG_TAIL)
