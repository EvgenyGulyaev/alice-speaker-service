# Alice Speaker Service Design

## Goal

Отдельный Go-микросервис, который подключает несколько Яндекс-аккаунтов через официальный API умного дома Яндекса, синхронизирует комнаты, колонки и сценарии и даёт `dashboard` безопасный HTTP API для вызова сценариев уведомлений на выбранные устройства.

## Why Separate Service

- Интеграция с Яндексом имеет собственные OAuth-токены, доступы `iot:view` и `iot:control`, модели аккаунтов и устройств.
- `dashboard` должен хранить только пользовательскую настройку выбора аккаунта, комнаты, колонки и сценария.
- Позже можно отдельно добавить ACL, retries, очереди и, при необходимости, неофициальный adapter для произвольного TTS, не загрязняя код админки.

## Constraints

- Используем только официальный Yandex Smart Home API в `v1`.
- Прямой произвольный TTS на Яндекс Станцию не обещаем: в `v1` работаем через сценарии и официальные сущности.
- Поддерживаем несколько Яндекс-аккаунтов.
- Любой админ в `dashboard` может вызывать сервис; ACL откладываем.
- Сервис пишем на Go в стиле существующего backend `dashboard`.

## Architecture

### Components

- `cmd/alice/main.go`
  - вход в приложение, инициализация config, db, http server.
- `internal/config`
  - загрузка `.env` и переменных окружения.
- `internal/http/server.go`
  - роутинг в стиле `dashboard`, middleware, health.
- `internal/http/middleware`
  - Bearer auth для внутренних вызовов от `dashboard`.
- `internal/http/routes`
  - CRUD аккаунтов, refresh, listing устройств/сценариев, announce.
- `internal/store`
  - BoltDB-репозитории аккаунтов, устройств, сценариев, журналов доставок.
- `internal/model`
  - `Account`, `Room`, `Device`, `Scenario`, `Delivery`.
- `internal/yandex`
  - клиент официального API Яндекса (`/v1.0/user/info`, `/v1.0/scenarios/{id}/actions`).

### Data Ownership

- `alice-speaker-service` хранит:
  - подключённые Яндекс-аккаунты и токены;
  - кэш синхронизированных комнат, колонок и сценариев;
  - журнал попыток отправки.
- `dashboard` хранит:
  - `alice_account_id`
  - `alice_room_id`
  - `alice_device_id`
  - `alice_scenario_id`
  - эти поля лежат в профиле пользователя.

## External Integration

### Yandex Smart Home API

- OAuth-приложение Яндекса с доступами:
  - `iot:view`
  - `iot:control`
- Сервис получает и хранит пользовательский OAuth token.
- Основные запросы:
  - `GET https://api.iot.yandex.net/v1.0/user/info`
  - `POST https://api.iot.yandex.net/v1.0/scenarios/{scenario_id}/actions`

### Dashboard Integration

- `dashboard` ходит в `alice-speaker-service` по внутреннему Bearer token.
- Для `v1` нужны эндпоинты:
  - `GET /health`
  - `GET /api/accounts`
  - `GET /api/accounts/:id/resources`
  - `POST /api/accounts/:id/refresh`
  - `POST /api/announce/scenario`

## User Flow

### Settings Flow

1. Админ открывает настройки.
2. Видит список подключённых Яндекс-аккаунтов.
3. Выбирает аккаунт.
4. После выбора получает список комнат, колонок и сценариев из микросервиса.
5. Сохраняет одну персональную связку:
   - аккаунт
   - комната
   - колонка
   - сценарий

### Manual Send Flow

1. В чате админ нажимает кнопку `На Алису`.
2. `dashboard` берёт настройки получателя.
3. Если настройки не заполнены, отправка не выполняется.
4. Если всё заполнено, `dashboard` вызывает `alice-speaker-service`.
5. Микросервис запускает выбранный сценарий через API Яндекса.
6. Результат пишется в журнал доставок и возвращается в `dashboard`.

## API Contract v1

### GET /health

Возвращает состояние сервиса.

### GET /api/accounts

Возвращает список подключённых Яндекс-аккаунтов:

- `id`
- `title`
- `is_active`
- `last_synced_at`

### GET /api/accounts/:id/resources

Возвращает агрегированные данные:

- `rooms`
- `devices`
- `scenarios`

### POST /api/accounts/:id/refresh

Пересинхронизирует `user/info` по аккаунту и обновляет локальный кэш.

### POST /api/announce/scenario

Тело запроса:

- `account_id`
- `scenario_id`
- `initiator_email`
- `recipient_email`
- `conversation_id`
- `message_id`

Ответ:

- `status`
- `request_id`
- `delivery_id`

## Security

- Внутренний Bearer token между `dashboard` и `alice-speaker-service`.
- OAuth tokens Яндекса храним в сервисе, не в `dashboard`.
- Логируем `request_id` от Яндекса для диагностики.
- В `v1` без ACL, но все сущности уже имеют `account_id`, чтобы ограничения можно было добавить позже.

## Error Handling

- Неподключённый аккаунт, пустые ресурсы или битый OAuth token возвращают читаемую ошибку.
- Если `dashboard` вызывает announce без сохранённых настроек получателя, это валидируем на стороне `dashboard` и дополнительно дублируем в сервисе.
- Ошибки Яндекса сохраняем в журнале доставок и прокидываем вверх с `request_id`.

## Testing

- unit tests для:
  - config
  - auth middleware
  - store
  - yandex client payload normalization
  - route handlers
- integration-style tests для:
  - refresh account resources
  - announce scenario

## Future v2

- ACL по аккаунтам/комнатам/колонкам
- авто-напоминания в `dashboard`
- retries / delivery queue
- неофициальный adapter для произвольного TTS на колонку, если решим идти в reverse-engineering

