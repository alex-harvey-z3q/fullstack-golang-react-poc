.PHONY: lint lint-fix fmt test up up-d down logs health api-tasks web-dev web-build web-preview test-backend migrate migrate-test reset-db sqlc web-ng-dev web-ng-build

lint:
	cd services/tasks && golangci-lint run ./...

lint-fix:
	cd services/tasks && golangci-lint run --fix ./...

fmt:
	gofmt -s -w services/tasks

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

# React frontend helpers
web-dev:
	cd services/web/react && npm install && npm run dev

web-build:
	cd services/web/react && npm install && npm run build

web-preview:
	cd services/web/react && npm install && npm run preview

# Angular frontend helpers
web-ng-dev:
	cd services/web/angular && npm install && npm run start

web-ng-build:

paste:
	find * -type d -name node_modules -prune -o -type f \
		-not -path "services/web/react/package-lock.json" \
		-not -path "services/tasks/graph/generated.go" \
		-not -path "scripts/interview.go" \
		-exec echo "===" \; \
		-exec echo {} \; \
		-exec echo "===" \; \
		-exec cat {} \;
