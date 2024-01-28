PROGRAM_NAME = zeroward

COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags | cut -d- -f1)

LDFLAGS = -ldflags "-a -installsuffix cgo -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"

.PHONY: help dep build clean

.DEFAULT_GOAL := help

help:
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dep:
	go mod tidy
	go mod download

build: dep
    mkdir -p ./bin
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${PROGRAM_NAME} ./main.go

clean:
	rm -rf ./bin

docker-build:
	docker build -t abdiaoo/zeroward:${TAG}
	docker image prune --force --filter label=stage=intermediate

docker-push:
	docker push abdiaoo/zeroward:${TAG}
