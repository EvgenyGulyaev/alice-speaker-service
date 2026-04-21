# Alice Speaker Service Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Построить отдельный Go-сервис, который подключает Яндекс-аккаунты, синхронизирует комнаты/колонки/сценарии и даёт `dashboard` внутренний API для запуска сценариев уведомлений.

**Architecture:** Сервис следует структуре backend `dashboard`: `cmd`, `internal/http`, `internal/store`, `internal/model`. Данные хранятся в BoltDB, внутренние вызовы защищены Bearer token, интеграция с Яндексом изолирована в `internal/yandex`.

**Tech Stack:** Go, BoltDB, net/http, dotenv-style config, internal Bearer auth, Yandex Smart Home API.

---

### Task 1: Repository bootstrap

**Files:**
- Create: `cmd/alice/main.go`
- Create: `internal/config/config.go`
- Create: `internal/http/server.go`
- Create: `internal/store/buckets.go`
- Create: `pkg/db/db.go`
- Modify: `.gitignore`

- [ ] **Step 1: Add repo hygiene**

Create `.gitignore` entries for `.idea`, `.env`, local binary, db files and logs.

- [ ] **Step 2: Add config loader**

Implement config loading for:
- `PORT`
- `ALICE_SERVICE_TOKEN`
- `ALICE_DB_PATH`

- [ ] **Step 3: Add db bootstrap**

Create BoltDB opener and bucket initialization helper.

- [ ] **Step 4: Add HTTP server skeleton**

Expose:
- `GET /health`

- [ ] **Step 5: Run smoke test**

Run:
```bash
go test ./...
```

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "Bootstrap Alice speaker service"
```

### Task 2: Models and repositories

**Files:**
- Create: `internal/model/account.go`
- Create: `internal/model/resource.go`
- Create: `internal/model/delivery.go`
- Create: `internal/store/account.go`
- Create: `internal/store/resource.go`
- Create: `internal/store/delivery.go`
- Test: `internal/store/store_test.go`

- [ ] **Step 1: Write failing repository tests**

Cover account create/list, resource upsert/list, delivery append.

- [ ] **Step 2: Run targeted tests**

Run:
```bash
go test ./internal/store -count=1
```

- [ ] **Step 3: Implement minimal models and repositories**

Support:
- accounts
- rooms
- devices
- scenarios
- deliveries

- [ ] **Step 4: Re-run store tests**

Run:
```bash
go test ./internal/store -count=1
```

- [ ] **Step 5: Commit**

```bash
git add internal/model internal/store
git commit -m "Add Alice service storage models"
```

### Task 3: Internal auth and route skeleton

**Files:**
- Create: `internal/http/middleware/auth.go`
- Create: `internal/http/routes/getHealth.go`
- Create: `internal/http/routes/getAccounts.go`
- Create: `internal/http/routes/getAccountResources.go`
- Create: `internal/http/routes/postAccountRefresh.go`
- Create: `internal/http/routes/postAnnounceScenario.go`
- Test: `internal/http/routes/routes_test.go`
- Modify: `internal/http/server.go`

- [ ] **Step 1: Write failing HTTP tests**

Cover:
- unauthorized request rejected
- health success
- accounts list success

- [ ] **Step 2: Run route tests to confirm failure**

Run:
```bash
go test ./internal/http/routes -count=1
```

- [ ] **Step 3: Implement Bearer auth middleware**

Compare `Authorization: Bearer <token>` with `ALICE_SERVICE_TOKEN`.

- [ ] **Step 4: Implement route skeletons**

Return placeholder but typed JSON for:
- accounts
- resources
- refresh
- announce

- [ ] **Step 5: Re-run route tests**

Run:
```bash
go test ./internal/http/routes -count=1
```

- [ ] **Step 6: Commit**

```bash
git add internal/http
git commit -m "Add Alice service auth and route skeleton"
```

### Task 4: Yandex Smart Home client

**Files:**
- Create: `internal/yandex/client.go`
- Create: `internal/yandex/types.go`
- Create: `internal/yandex/client_test.go`

- [ ] **Step 1: Write failing client tests**

Cover normalization of:
- rooms
- smart speakers
- scenarios
- request headers

- [ ] **Step 2: Run tests and confirm failure**

Run:
```bash
go test ./internal/yandex -count=1
```

- [ ] **Step 3: Implement `GET /v1.0/user/info` client**

Parse:
- rooms
- devices
- scenarios

Filter supported devices to Yandex speakers / TV Station device types.

- [ ] **Step 4: Implement `POST /v1.0/scenarios/{id}/actions` client**

Return Yandex `request_id` in result.

- [ ] **Step 5: Re-run yandex tests**

Run:
```bash
go test ./internal/yandex -count=1
```

- [ ] **Step 6: Commit**

```bash
git add internal/yandex
git commit -m "Add Yandex smart home client"
```

### Task 5: Refresh and announce handlers

**Files:**
- Modify: `internal/http/routes/postAccountRefresh.go`
- Modify: `internal/http/routes/postAnnounceScenario.go`
- Modify: `internal/http/routes/getAccountResources.go`
- Modify: `internal/store/resource.go`
- Modify: `internal/store/delivery.go`
- Test: `internal/http/routes/routes_test.go`

- [ ] **Step 1: Extend tests for refresh flow**

Assert that refresh stores synchronized rooms, devices and scenarios.

- [ ] **Step 2: Extend tests for announce flow**

Assert announce:
- validates account/scenario
- writes delivery log
- returns status + request id

- [ ] **Step 3: Implement refresh logic**

`POST /api/accounts/:id/refresh`:
- load account token
- call Yandex user info
- upsert resources

- [ ] **Step 4: Implement announce logic**

`POST /api/announce/scenario`:
- validate payload
- call Yandex scenario action
- store delivery record

- [ ] **Step 5: Re-run route tests**

Run:
```bash
go test ./internal/http/routes -count=1
```

- [ ] **Step 6: Commit**

```bash
git add internal/http/routes internal/store
git commit -m "Implement Alice account refresh and announce"
```

### Task 6: Account management endpoints

**Files:**
- Create: `internal/http/routes/postAccount.go`
- Create: `internal/http/routes/patchAccount.go`
- Modify: `internal/http/server.go`
- Test: `internal/http/routes/routes_test.go`

- [ ] **Step 1: Write failing tests for account CRUD-lite**

Cover:
- create account
- update title/token/active flag
- list accounts

- [ ] **Step 2: Run tests**

Run:
```bash
go test ./internal/http/routes -count=1
```

- [ ] **Step 3: Implement handlers**

Support:
- storing OAuth token per account
- active flag
- display title

- [ ] **Step 4: Re-run tests**

Run:
```bash
go test ./internal/http/routes -count=1
```

- [ ] **Step 5: Commit**

```bash
git add internal/http/routes internal/http/server.go
git commit -m "Add Alice account management endpoints"
```

### Task 7: End-to-end verification and docs

**Files:**
- Modify: `README.md`
- Modify: `docs/superpowers/specs/2026-04-21-alice-speaker-service-design.md`
- Modify: `docs/superpowers/plans/2026-04-21-alice-speaker-service.md`

- [ ] **Step 1: Document env and local run**

Explain:
- OAuth scopes
- service token
- db path
- example curl

- [ ] **Step 2: Run full test suite**

Run:
```bash
go test ./... -count=1
```

- [ ] **Step 3: Manual smoke check**

Run service locally and hit:
```bash
curl http://localhost:8090/health
```

- [ ] **Step 4: Commit**

```bash
git add README.md docs
git commit -m "Document Alice speaker service setup"
```

