IMG ?= adamdyszy/sportsnews
$(eval GIT_COMMIT := $(shell git rev-parse --short HEAD))
TAG ?= $(GIT_COMMIT)
ifeq ($(strip $(TAG)),)
	TAG = $(shell find . -path './vendor' -prune -type f -o -name "*.go" -type f -print0 | sort -z | xargs -0 cat | md5sum | cut -d ' ' -f1)
endif
NEWS_CONTAINER_NAME ?= news${TAG}
MONGO_CONTAINER_NAME ?= mongocontainer${TAG}
NEWS_MEM_CONTAINER_NAME ?= newsmem${TAG}

CONFIG_FILE_MEMORY ?= $(CONFIG_FILE)
ifeq ($(strip $(CONFIG_FILE_MEMORY)),)
    CONFIG_FILE_MEMORY = config/quickstart/memory.yaml
endif
CONFIG_FILE_MONGO ?= $(CONFIG_FILE)
ifeq ($(strip $(CONFIG_FILE_MONGO)),)
    CONFIG_FILE_MONGO = config/quickstart/mongo.yaml
endif
CONFIG_FILE ?= config/custom.yaml

.PHONY: all
all: test build docker-build

.PHONY: test
test:
	go test ./...

.PHONY: build
build: test
	go build -o=bin/sportsnews -mod=vendor cmd/sportsnews/main.go

.PHONY: run run-memory
run run-memory:
	go run -mod=vendor cmd/sportsnews/main.go --customConfigFile=$(CONFIG_FILE_MEMORY)

.PHONY: run-mongo
run-mongo:
	go run -mod=vendor cmd/sportsnews/main.go --customConfigFile=$(CONFIG_FILE_MONGO)

.PHONY: docker-build
docker-build: test
	docker build -t ${IMG}:${TAG} .

.PHONY: docker-run
docker-run: docker-build
	docker run -it --rm -v "${PWD}/config:/config" -p 8080:8080 "${IMG}:${TAG}" --customConfigFile=$(CONFIG_FILE)

.PHONY: docker-run-background
docker-run-background: docker-build
	docker run -d --rm -v "${PWD}/config:/config" -p 8080:8080 "${IMG}:${TAG}" --customConfigFile=$(CONFIG_FILE)

.PHONY: quickstart quickstart-mongo
quickstart quickstart-mongo: docker-build
	@if docker ps --filter "name=${MONGO_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${MONGO_CONTAINER_NAME}" > /dev/null; then \
		echo "Mongo container (${MONGO_CONTAINER_NAME}) is already running."; \
	else \
		docker run -d --rm --name ${MONGO_CONTAINER_NAME} \
		-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
		-e MONGO_INITDB_ROOT_PASSWORD=secret \
		mongo; \
		echo "started mongodb with container name: ${MONGO_CONTAINER_NAME}"; \
	fi
	@if docker ps --filter "name=${NEWS_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${NEWS_CONTAINER_NAME}" > /dev/null; then \
		echo "News container (${NEWS_CONTAINER_NAME}) is already running."; \
	else \
		docker run -d --rm -v "${PWD}/config:/config" \
		--link ${MONGO_CONTAINER_NAME}:mongo -p 8080:8080 \
		--name ${NEWS_CONTAINER_NAME} "${IMG}:${TAG}" --customConfigFile=$(CONFIG_FILE_MONGO); \
		echo "started server with container name: ${NEWS_CONTAINER_NAME}"; \
	fi
	@if docker ps --filter "name=${MONGO_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${MONGO_CONTAINER_NAME}" > /dev/null; then \
		echo "Mongo container (${MONGO_CONTAINER_NAME}) is running."; \
	else \
		echo "Mongo container (${MONGO_CONTAINER_NAME}) failed to start."; \
		exit 1; \
	fi
	@if docker ps --filter "name=${NEWS_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${NEWS_CONTAINER_NAME}" > /dev/null; then \
		echo "News container (${NEWS_CONTAINER_NAME}) is running."; \
	else \
		echo "News container (${NEWS_CONTAINER_NAME}) failed to start."; \
		exit 1; \
	fi
	@echo "to kill created pods run"
	@echo "make quickstart-kill"
	@echo "or:"
	@echo "docker rm -f ${MONGO_CONTAINER_NAME}"
	@echo "docker rm -f ${NEWS_CONTAINER_NAME}"
	@echo "to get logs of the server and db:"
	@echo "docker logs ${NEWS_CONTAINER_NAME}"
	@echo "docker logs ${MONGO_CONTAINER_NAME}"
	@echo "to curl the server use:"
	@echo "curl localhost:8080/articles"
	@echo "to curl the server with article id use:"
	@echo "curl localhost:8080/articles/{ID}"
	@echo "restart the server with (this will trigger getting details for articles):"
	@echo "docker restart ${NEWS_CONTAINER_NAME}"
	@echo "to show these options again simply run make quickstart again"

.PHONY: quickstart-restart-server quickstart-mongo-restart-server
quickstart-restart-server quickstart-mongo-restart-server:
	docker restart ${NEWS_CONTAINER_NAME}

.PHONY: quickstart-mongo-kill quickstart-kill
quickstart-mongo-kill quickstart-kill:
	docker rm -f ${MONGO_CONTAINER_NAME}
	docker rm -f ${NEWS_CONTAINER_NAME}

.PHONY: quickstart-mem
quickstart-mem: docker-build
	@if docker ps --filter "name=${NEWS_MEM_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${NEWS_MEM_CONTAINER_NAME}" > /dev/null; then \
		echo "News container (${NEWS_MEM_CONTAINER_NAME}) is already running."; \
	else \
		docker run -d --rm -p 8080:8080 --name ${NEWS_MEM_CONTAINER_NAME} "${IMG}:${TAG}"; \
		echo "started server with container name: ${NEWS_MEM_CONTAINER_NAME}"; \
	fi
	@if docker ps --filter "name=${NEWS_MEM_CONTAINER_NAME}" --format "{{.Names}}" | grep -w "${NEWS_MEM_CONTAINER_NAME}" > /dev/null; then \
		echo "News container (${NEWS_MEM_CONTAINER_NAME}) started OK."; \
	else \
		echo "News container (${NEWS_MEM_CONTAINER_NAME}) failed to start."; \
		exit 1; \
	fi
	@echo "to kill created pods run"
	@echo "make quickstart-mem-kill"
	@echo "or:"
	@echo "docker rm -f ${NEWS_MEM_CONTAINER_NAME}"
	@echo "to get logs of the server:"
	@echo "docker logs ${NEWS_MEM_CONTAINER_NAME}"
	@echo "to curl the server use:"
	@echo "curl localhost:8080/articles"
	@echo "to curl the server with article id use:"
	@echo "curl localhost:8080/articles/{ID}"
	@echo "to show these options again simply run make quickstart-mem again"

.PHONY: quickstart-mem-kill
quickstart-mem-kill:
	docker rm -f ${NEWS_MEM_CONTAINER_NAME}