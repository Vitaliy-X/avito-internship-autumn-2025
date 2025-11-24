# PR Reviewer Assignment Service (Fall 2025)

Сервис для автоматического назначения ревьюверов на Pull Request’ы, управления командами и участниками.

---

### Запуск проекта

```bash
make up
```

### Управление проектом через Makefile

#### Полный запуск с пересборкой
```bash
make up
```

#### Запуск в фоне (detached mode)
```bash
make up-d
```

#### Остановка контейнеров
```bash
make stop
````

#### Полная остановка + удаление контейнеров
```bash
make down
```

#### Полная очистка + удаление volumes
```bash
make clean
```

#### Перезапуск контейнеров
```bash
make restart
```

#### Применить миграции
```bash
make migrate
```

#### Откатить последнюю миграцию
```bash
make rollback
```

## Реализованные эндпоинты

- `POST /team/add` — создать или обновить команду с участниками
- `GET /team/get` — получить команду с участниками
- `POST /users/setIsActive` — установить флаг активности пользователя
- `POST /pullRequest/create` — создать PR и автоматически назначить до 2 ревьюверов
- `POST /pullRequest/merge` — меняет статус PR на "merged" (идемпотентно), устанавливает время "mergedAt"
