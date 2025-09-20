# Quick Start

This guide is for developers who want to clone the repo and get the stack running quickly.

## Prerequisites

- **Docker & Docker Compose**
- **Go 1.23+**
- **golangci-lint v2**
- **Node.js 18+** with npm

## Clone the repository

```bash
git clone https://github.com/alex-harvey-z3q/fullstack-golang-react-poc.git
cd fullstack-golang-react-poc
```

## Run the backend stack (Postgres + tasks service)

### Option A — foreground
```bash
make up
```

### Option B — background (detached)
```bash
make up-d
```

Check health:
```bash
make health
```

List tasks:
```bash
make api-tasks
```

Stop the stack:
```bash
make down
```

## Run the frontend (React client)

In a separate terminal:
```bash
make web-dev
```
- React dev server runs at: <http://localhost:5173>
- It proxies `/api` requests to the backend at `http://localhost:8081`.

## Database helpers

Apply latest migration to **dev** DB:
```bash
make migrate
```

Apply migration to **test** DB:
```bash
make migrate-test
```

Reset dev DB (delete all rows and reset IDs):
```bash
make reset-db
```

## Backend development

Run tests:
```bash
make test-backend
```

Generate sqlc code:
```bash
make sqlc
```

Generate gqlgen code (if/when GraphQL is wired):
```bash
make gqlgen
```

## Build the frontend for production

```bash
make web-build
# optionally preview the built app locally
make web-preview
```
