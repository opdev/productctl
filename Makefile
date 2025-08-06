GENERATED_PYXIS_CODE ?= internal/genpyxis/generated.go
OUT_DIR = $(shell pwd)/out

default: build

# Generate and Build Everything.
.PHONY: build
build: generate bin

BIN_FILE    ?= productctl
BIN_VERSION ?= unknown
BIN_COMMIT  = $(shell git rev-parse HEAD)

# Just build the binary.
.PHONY: bin
bin:
	CGO_ENABLED=0 go build \
		-o $(BIN_FILE) \
		-trimpath \
		-ldflags "\
			-s -w \
			-X github.com/opdev/productctl/internal/version.commit=$(BIN_COMMIT) \
			-X github.com/opdev/productctl/internal/version.version=$(BIN_VERSION)" \
		./internal/cmd/productctl

.PHONY: clean
clean:
	rm -vf $(BIN_FILE)
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

.PHONY: lint.fix
lint.fix: install.golangci-lint
	$(GOLANGCI_LINT) run --fix

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt: install.gofumpt
	${GOFUMPT} -l -w .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v ./...

.PHONY: cover
cover:
	go test \
	 -race \
	 -cover -coverprofile=coverage.out \
	 $$(go list ./... | grep -vP 'genpyxis|testutils')

### Fail if git diff detects a change. Useful for CI.
.PHONY: diff-check
diff-check:
	git diff --exit-code

.PHONY: ci.fmt
ci.fmt: fmt diff-check
	echo "=> ci.fmt done"

.PHONY: ci.generate
ci.generate: generate diff-check
	echo "=> ci.generate done"

.PHONY: ci.tidy
ci.tidy: tidy diff-check
	echo "=> ci.tidy done"

### Installing Developer Tools
# gofumpt
GOFUMPT = $(OUT_DIR)/gofumpt
install.gofumpt:
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt@latest)

# golangci-lint
GOLANGCI_LINT = $(OUT_DIR)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.1.6
install.golangci-lint:
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

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