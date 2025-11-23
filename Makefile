APP_NAME = pr-service
PROJECT = pr-service
COMPOSE_FILE = docker-compose.yml

# Detect docker compose command for Windows/Linux/Mac
DOCKER_COMPOSE = $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; else echo "docker-compose"; fi)

.PHONY: up up-d stop down clean restart migrate rollback

## Запуск всего проекта (с пересборкой)
up:
	$(DOCKER_COMPOSE) -p $(PROJECT) -f $(COMPOSE_FILE) up --build

## Запуск в фоне
up-d:
	$(DOCKER_COMPOSE) -p $(PROJECT) -f $(COMPOSE_FILE) up -d --build

## Остановка
stop:
	$(DOCKER_COMPOSE) -p $(PROJECT) stop

## Полная остановка + удаление контейнеров
down:
	$(DOCKER_COMPOSE) -p $(PROJECT) down

## Полная очистка контейнеров + volumes
clean:
	$(DOCKER_COMPOSE) -p $(PROJECT) down -v

## Перезапуск
restart:
	$(DOCKER_COMPOSE) -p $(PROJECT) restart

## Применить миграции
migrate:
	$(DOCKER_COMPOSE) -p $(PROJECT) run --rm liquibase \
	  --defaultsFile=/liquibase/liquibase.properties update

## Откатить одну миграцию
rollback:
	$(DOCKER_COMPOSE) -p $(PROJECT) run --rm liquibase \
	  --defaultsFile=/liquibase/liquibase.properties rollbackCount 1
