compose = docker-compose -f build/docker-compose.yml

.PHONY: config
config: ##@development Initial configuration files.
	cp config-example.yaml config.yaml
	cp example.env .env

.PHONY: build
build: ##@development Build development environment.
	$(compose) build

.PHONY: up
up: build ##@development Starts development environment in background.
	$(compose) up -d

.PHONY: test
test: build ##@development Runs the tests.
	$(compose) run --rm app go test ./...

.PHONY: logs
logs: ##@development Follows development logs.
	$(compose) logs -f --tail=100

lint_version ?= v1.40-alpine
.PHONY: lint
lint: ##@development Runs static analysis.
	docker run --rm \
		-w /app \
		-v $(shell pwd):/app \
		golangci/golangci-lint:$(lint_version) \
		golangci-lint run --timeout 3m

.PHONY: clean
clean: ##@development Stop development environment.
	$(compose) down --volumes --remove-orphans
