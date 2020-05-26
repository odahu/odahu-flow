all: help

## test: Run unit tests
test:
	gotestsum --junitfile feedback-report.xml -- -race -coverprofile=coverage.txt ./pkg/...

## build-all: Buld all services
build-all: build-collector build-rq-catcher

## build-collector: Launch collector service
run-collector:
	go run ./cmd/collector/main.go

## build-collector: Build collector service
build-collector:
	go build -o collector cmd/collector/main.go

## build-rq-catcher: Launch rq-catcher service
run-rq-catcher:
	go run ./cmd/rq-catcher/main.go

## build-rq-catcher: Build rq-catcher service
build-rq-catcher:
	go build -o rq-catcher cmd/rq-catcher/main.go

## lint: Start golangci-lint linter
lint:
	golangci-lint run ${LINTER_ADDITIONAL_ARGS}

## help: Show the help message
help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
