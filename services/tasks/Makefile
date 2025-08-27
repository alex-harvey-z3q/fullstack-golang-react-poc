PGURL=postgres://app:app@localhost:5432/tasks?sslmode=disable

.PHONY: dev sqlc gqlgen migrate run test
dev: ## run server
	go run ./cmd/server

sqlc: ## generate sqlc code
	sqlc generate -f internal/db/sqlc.yaml

gqlgen: ## generate gql code
	go run github.com/99designs/gqlgen generate

migrate: ## apply initial migration (requires psql)
	psql $(PGURL) -f internal/db/migrate/001_init.sql

run:
	PORT=8081 DATABASE_URL=$(PGURL) go run ./cmd/server

test:
	go test ./... -v
