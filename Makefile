.PHONY: test up up-d down logs health api-tasks web-dev web-build web-preview test-backend migrate migrate-test reset-db sqlc

test:
	cd services/tasks && go test ./... -v

up:
	docker compose up --build

up-d:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

health:
	@curl -sS http://localhost:8081/healthz || true

api-tasks:
	@curl -sS http://localhost:8081/api/tasks || true

migrate:
	$(MAKE) -C services/tasks migrate

migrate-test:
	$(MAKE) -C services/tasks migrate-test

sqlc:
	$(MAKE) -C services/tasks sqlc

test-backend:
	$(MAKE) -C services/tasks test

reset-db:
	$(MAKE) -C services/tasks reset-db

web-dev:
	cd services/web/react && npm install && npm run dev

web-build:
	cd services/web/react && npm install && npm run build

web-preview:
	cd services/web/react && npm install && npm run preview

paste:
	find * -type d -name node_modules -prune -o -type f \
		-not -path "services/web/react/package-lock.json" \
		-not -path "services/tasks/graph/generated.go" \
		-exec echo {} \; \
		-exec echo "===" \; \
		-exec cat {} \; \
		-exec echo "===" \;
