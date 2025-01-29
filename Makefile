BINFILE ?= productctl
GENERATED_PYXIS_CODE ?= internal/genpyxis/generated.go
OUT_DIR = $(shell pwd)/out

default: build

# Generate and Build Everything.
.PHONY: build
build: generate bin

# Just build the binary.
.PHONY: bin
bin:
	go build -o $(BINFILE) ./internal/cmd/productctl

.PHONY: clean
clean:
	rm -vf $(BINFILE)
	rm -vf $(GENERATED_PYXIS_CODE)
	rm -vrf $(OUT_DIR)

### Generating everything.
.PHONY: generate
generate: generate.schema generate.graphql

### Generate client code. Assumes schema is present.
.PHONY: generate.graphql
generate.graphql: install.genqlient
	go generate ./...

### Generating Catalog API GraphQL Schema
.PHONY: generate.schema
generate.schema: schemagen-venv
	$(SCHEMAGEN_VENV_BIN_DIR)/python3 scripts/generate-schema/generate-schema.py > internal/genpyxis/schema.graphql

SCHEMAGEN_VENV_DIR = $(OUT_DIR)/generate-schema-venv
SCHEMAGEN_VENV_BIN_DIR = $(OUT_DIR)/generate-schema-venv/bin
.PHONY: schemagen-venv
schemagen-venv:
	python3 -m venv $(SCHEMAGEN_VENV_DIR)
	$(SCHEMAGEN_VENV_BIN_DIR)/pip3 install -r scripts/generate-schema/requirements.txt

### Enforcing Project Standards
.PHONY: lint
lint: install.golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt: install.gofumpt
	${GOFUMPT} -l -w .

.PHONY: tidy
tidy:
	go mod tidy

### Installing Developer Tools
# gofumpt
GOFUMPT = $(OUT_DIR)/gofumpt
install.gofumpt:
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt@latest)

# golangci-lint
GOLANGCI_LINT = $(OUT_DIR)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.63.4
install.golangci-lint:
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

# genqlient
GENQLIENT = $(OUT_DIR)/genqlient
GENQLIENT_VERSION ?= v0.7.0
install.genqlient:
	$(call go-install-tool,$(GENQLIENT),github.com/Khan/genqlient@$(GENQLIENT_VERSION))


# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
GOBIN=$(PROJECT_DIR)/out go install $(2) ;\
}
endef