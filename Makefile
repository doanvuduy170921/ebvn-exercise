run:
	go run cmd/api/main.go

COVERAGE_EXCLUDE=mocks|main.go|integration_test
COVERAGE_THRESHOLD = 50


test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

test-integration:
	@echo "==> Running Integration Tests..."
	go test -v ./internal/integration_test/...

swagger:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

dev-run: swagger run

IMAGE_NAME = duydoanvu1709/go-api
IMAGE_TAG = $(shell git rev-parse --abbrev-ref HEAD | sed 's/main/dev/')

docker-test:
	@mkdir -p coverage
	docker buildx build --target test \
		--build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" \
		--output type=local,dest=./coverage .
	@echo "✅ Test completed. Reports are in ./coverage"

docker-login:
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin


docker-build:
	docker build --target final -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-release: docker-build
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: run test test-integration swagger dev-run docker-test docker-login docker-build docker-release