# alice-speaker-service

Go-микросервис для интеграции с официальным API умного дома Яндекса.

## Что умеет сейчас

- хранить несколько Яндекс-аккаунтов;
- синхронизировать комнаты, колонки и сценарии через `user/info`;
- запускать выбранный сценарий через API Яндекса;
- отдавать `dashboard` внутренний HTTP API.

## Env

```env
PORT=8090
ALICE_DB_PATH=./alice.db
ALICE_SERVICE_TOKEN=change-me
```

## Локальный запуск

```bash
export GOENV_VERSION=1.25.4
go run ./cmd/alice
```

## Основные эндпоинты

- `GET /health`
- `GET /api/accounts`
- `POST /api/accounts`
- `PATCH /api/accounts/:id`
- `GET /api/accounts/:id/resources`
- `POST /api/accounts/:id/refresh`
- `POST /api/announce/scenario`

## Пример добавления аккаунта

```bash
curl -X POST http://localhost:8090/api/accounts \
  -H "Authorization: Bearer change-me" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "home-main",
    "title": "Основной дом",
    "oauth_token": "YANDEX_OAUTH_TOKEN",
    "is_active": true
  }'
```

## Что нужно от Яндекса

OAuth-приложение со scope:

- `iot:view`
- `iot:control`
