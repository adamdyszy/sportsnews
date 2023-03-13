IMG ?= adamdyszy/sportsnews
$(eval GIT_COMMIT := $(shell git rev-parse --short HEAD))
TAG ?= $(GIT_COMMIT)

.PHONY: all
all: test build docker-build

.PHONY: test
test:
	go test ./...

.PHONY: build
build: test
	go build -o=bin/sportsnews -mod=vendor cmd/sportsnews/main.go

.PHONY: docker-build
docker-build: test
	docker build -t ${IMG}:${TAG} .

.PHONY: docker-run
docker-run: docker-build
	docker run -it --rm -v ${PWD}/config:/config -p 8080:8080 ${IMG}:${TAG}

.PHONY: docker-run-background
docker-run-background: docker-build
	docker run -d --rm -v ${PWD}/config:/config -p 8080:8080 ${IMG}:${TAG}

.PHONY: quickstart quickstart-mongo
quickstart quickstart-mongo: docker-build quickstart-config-mongo
	docker run -d --rm --name mongocontainer${TAG} \
    	-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
    	-e MONGO_INITDB_ROOT_PASSWORD=secret \
    	mongo
	docker run -d --rm -v ${PWD}/config:/config \
        --link mongocontainer${TAG}:mongo -p 8080:8080 \
        --name news${TAG} ${IMG}:${TAG}

.PHONY: quickstart-mongo-kill quickstart-kill
quickstart-mongo-kill quickstart-kill:
	docker rm -f mongocontainer${TAG}
	docker rm -f news${TAG}

.PHONY: quickstart-config-mongo
quickstart-config-mongo:
	echo "storageKind: \"mongo\"" > config/custom.yaml
	echo "mongoStorage:" >> config/custom.yaml
	echo "  uri: \"mongodb://mongo:27017\"" >> config/custom.yaml
	echo "  user: \"mongoadmin\"" >> config/custom.yaml
	echo "  password: \"secret\"" >> config/custom.yaml

.PHONY: quickstart-mem
quickstart-mem: docker-build
	docker run -d --rm -p 8080:8080 --name news${TAG} ${IMG}:${TAG}

.PHONY: quickstart-mem-kill
quickstart-mem-kill:
	docker rm -f news${TAG}