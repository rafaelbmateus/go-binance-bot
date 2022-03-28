compose = docker-compose -f build/docker-compose.yml

.PHONY: up
up: ##@development Build and start development environment in background.
	$(compose) up --build

.PHONY: test
test: ##@development Runs the tests.
	$(compose) run --rm app go test ./...

.PHONY: logs
logs: ##@development Follows development logs [service="svc1 svc2..."].
	$(compose) logs -f --tail=100 $(service)

lint_version ?= v1.40-alpine
.PHONY: lint
lint: ##@development Runs static analysis.
	docker run --rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:$(lint_version) \
		golangci-lint run --timeout 3m

.PHONY: clean
clean: ##@development Stop development environment.
	$(compose) down -v --remove-orphans
