.PHONY: build kubectl plugin 

GOENV  := GO15VENDOREXPERIMENT="1" GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO := $(GOENV) go

default: build

build: view-node-resource view-node-taints

plugin:
	GO111MODULE=on CGO_ENABLED=0 go build -o kubectl-debug cmd/plugin/main.go

view-node-resource:
	$(GO) build -o kubectl-view-node-resource cmd/view-node-resource/main.go

view-node-taints:
	$(GO) build -o kubectl-view-node-taints cmd/view-node-taints/main.go
