.PHONY: test
test:
	cd services/tasks && go test ./... -v

paste:
	find * -type d -name node_modules -prune -o -type f -not -path "services/web/react/package-lock.json" -not -path "services/tasks/graph/generated.go" -exec wc -l {} \; -exec echo {} \; -exec echo "===" \; -exec cat {} \; -exec echo "===" \;
