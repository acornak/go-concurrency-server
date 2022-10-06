PORT=4001
ENV=develop
EXPONEA_URL="https://exponea-engineering-assignment.appspot.com/api/work"
BINARY_NAME=go-concurrency-server

## building binaries
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/api
	@echo "Built!"

## building and running binaries
run: build
	@echo "Starting..."
	@env PORT=${PORT} ENV=${ENV} EXPONEA_URL=${EXPONEA_URL} TIMEOUT=$(TIMEOUT) ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...