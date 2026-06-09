run:
	go run cmd/api/main.go


test:
	go test ./... -coverprofile=coverage.out

test-integration:
	@echo "==> Running Integration Tests..."
	go test -v ./internal/integration_test/...


swagger:
	swag init -g cmd/api/main.go

dev-run: swagger run

.PHONY: run test test-integration swagger dev-run