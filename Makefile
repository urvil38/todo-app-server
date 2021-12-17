VERSION := 0.0.1
GO := go
GOOS := $(go env GOOS)

NAME := todo-app
BIN := todo-app-server
PKG := github.com/urvil38/todo-app
ECR_REGISTRY := 910481930765.dkr.ecr.ap-south-1.amazonaws.com

BUILD_COMMIT := $(shell ./build/get-build-commit.sh)

LDFLAGS := "-X $(PKG)/internal/version.Version=$(VERSION) \
	-X $(PKG)/internal/version.Commit=$(BUILD_COMMIT) \
	-w"

GO111MODULE := on

.PHONY: all
all: test build push clean

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: build
build: go-build

.PHONY: docker-build
docker-build:
	@docker build -t urvil38/$(NAME):$(VERSION) -t $(ECR_REGISTRY)/$(NAME):$(VERSION) \
		-f Dockerfile \
		.

	@echo Successfully built $(BIN)

push:
	docker push urvil38/$(NAME):$(VERSION)

aws-push:
	aws ecr get-login-password --region ap-south-1 | docker login --username AWS --password-stdin 910481930765.dkr.ecr.ap-south-1.amazonaws.com
	docker push $(ECR_REGISTRY)/$(NAME):$(VERSION)
	docker logout

clean:
	rm -f $(BIN) -y

go-build:
	GOOS=$(GOOS) $(GO) build -o $(BIN) -ldflags $(LDFLAGS)

.PHONY: test
test: go-test

.PHONY: go-test
go-test:
	$(GO) test -v -count=1 ./...