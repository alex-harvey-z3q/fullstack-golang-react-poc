# Uses psql inside the docker-compose postgres container.
# Requires: `docker compose up` is running.

PSQL       := docker compose exec -T postgres psql -U app
DB_DEV     := tasks
DB_TEST    := tasks_test
MIG        := internal/db/migrate/001_init.sql

.PHONY: dev sqlc gqlgen migrate createdb-test migrate-test test tidy

## Run the server locally (PORT/DATABASE_URL via env or defaults)
dev:
	go run ./cmd/server

## Generate sqlc code
sqlc:
	sqlc generate -f internal/db/sqlc.yaml

## Generate GraphQL code (optional if/when you wire gqlgen)
gqlgen:
	go run github.com/99designs/gqlgen generate --config graph/gqlgen.yml

## Apply migration to the dev DB (inside container)
migrate: $(MIG)
	$(PSQL) -d $(DB_DEV) -v ON_ERROR_STOP=1 -f - < $(MIG)

## Create the test DB once (inside container)
createdb-test:
	@set -e; \
	EXISTS=$$($(PSQL) -d $(DB_DEV) -Atqc "SELECT 1 FROM pg_database WHERE datname = '$(DB_TEST)'"); \
	if [ "$$EXISTS" != "1" ]; then \
		echo "Creating database $(DB_TEST)"; \
		$(PSQL) -d $(DB_DEV) -v ON_ERROR_STOP=1 -c "CREATE DATABASE $(DB_TEST)"; \
	else \
		echo "Database $(DB_TEST) already exists"; \
	fi

## Apply migration to the test DB (run manually when schema changes)
migrate-test: createdb-test $(MIG)
	$(PSQL) -d $(DB_TEST) -v ON_ERROR_STOP=1 -f - < $(MIG)

## Run unit tests against the isolated test DB (no auto-migrations here)
test:
	TEST_DATABASE_URL=postgres://app:app@localhost:5432/$(DB_TEST)?sslmode=disable go test ./... -v

## Go module housekeeping
tidy:
	go mod tidy
