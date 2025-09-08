.PHONY: test
test:
	cd services/tasks && go test ./... -v

paste:
	find * -type d -name node_modules -prune -o -type f -exec echo {} \; -exec echo "===" \; -exec cat {} \; -exec echo "===" \;
