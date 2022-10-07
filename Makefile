CONTAINER=go-concurrency-server

## build docker image
build:
	@echo "Building..."
	docker compose build
	@echo "Built!"

## running docker container
run: build
	@echo "Starting..."
	@docker compose up &
	@echo "Started!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@docker stop $(CONTAINER) || true && docker rm ${CONTAINER} || true
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start
