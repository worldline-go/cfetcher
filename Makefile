BINARY_NAME = cfetcher
MAIN_PATH = ./cmd/cfetcher/main.go

BIN_PATH = bin
LOCAL_BIN_DIR := $(PWD)/$(BIN_PATH)

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_VERSION := $(or $(IMAGE_TAG),$(shell git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0))

.DEFAULT_GOAL := help

.PHONY: build
build: build-linux ## Build CLI for Linux

.PHONY: lint
lint: ## Run linter
	GOPATH=$(shell dirname ${PWD}) golangci-lint run ./...

.PHONY: build-windows
build-windows: ## Build CLI for Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.date=$(BUILD_DATE)" -o $(BIN_PATH)/$(BINARY_NAME).exe $(MAIN_PATH)

.PHONY: build-linux
build-linux: ## Build CLI for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(BUILD_COMMIT) -X main.date=$(BUILD_DATE)" -o $(BIN_PATH)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run-load
run-load: export VAULT_ADDR = http://localhost:8080
run-load: export CONSUL_HTTP_ADDR = http://localhost:8080
run-load: ## Run load test to struct
	VAULT_ADDR=$(VAULT_ADDR) CONSUL_HTTP_ADDR=$(CONSUL_HTTP_ADDR) go run ./cmd/load/main.go

.PHONY: run-consul
run-consul: export CONFIG_ENV ?= local
run-consul: ## Run consul save
	go run ./cmd/config/main.go --consul-save=./out/$(CONFIG_ENV)/finops-consul

.PHONY: run-vault
run-vault: export CONFIG_ENV ?= local
run-vault: ## Run vault save
	go run ./cmd/config/main.go --vault-save=./out/$(CONFIG_ENV)/finops-vault

.PHONY: vault
vault: ## Run vault server
	docker run -d --name vault -p 8200:8200 --cap-add IPC_LOCK vault:1.13.3

.PHONY: vault-log
vault-log: ## Show vault server logs
	docker logs -f vault

.PHONY: vault-destroy
vault-destroy: ## Destroy vault server
	docker rm -f vault

.PHONY: vault-login
vault-login: TOKEN ?= root
vault-login: ## Login to vault
	vault login $(TOKEN)

.PHONY: vault-role-enable
vault-role-enable: ## Enable vault role
	vault auth enable approle

.PHONY: vault-role
.ONESHELL:
vault-role:
	cat <<-EOF | vault policy write secret-read -
		path "finops/*" {
		capabilities = ["read", "list"]
		}
	EOF
	# vault write auth/approle/role/my-role bind_secret_id=false secret_id_bound_cidrs="0.0.0.0/0" policies="default,secret-read"
	vault write auth/approle/role/my-role policies="default,secret-read"
	vault read -field=role_id auth/approle/role/my-role/role-id

.PHONY: vault-role-destroy
vault-role-destroy:
	vault delete auth/approle/role/my-role

.PHONY: vault-secret
vault-secret:
	vault write -f auth/approle/role/my-role/secret-id

.PHONY: test
test: ## Run unit tests
	@go test -v -race -cover ./...

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
