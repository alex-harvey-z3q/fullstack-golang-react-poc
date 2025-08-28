.PHONY: test
test:
	cd services/tasks && go test ./... -v
