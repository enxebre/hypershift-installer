SRC_DIRS = cmd pkg


.PHONY: default
default: verify build

.PHONY: build
#build:  bindata control-plane-operator
build: bindata
	go build -mod=vendor -o bin/hypershift-installer github.com/openshift-hive/hypershift-installer/cmd

.PHONY: bindata
bindata:
	hack/update-bindata.sh

.PHONY: verify-bindata
verify-bindata:
	hack/verify-bindata.sh

.PHONY: verify-gofmt
verify-gofmt:
	@echo Verifying gofmt
	@gofmt -l -s $(SRC_DIRS)>.out 2>&1 || true
	@[ ! -s .out ] || \
	  (echo && echo "*** Please run 'make fmt' in order to fix the following:" && \
	  cat .out && echo && rm .out && false)
	@rm .out

.PHONY: verify
verify: verify-gofmt verify-bindata

# ---- CAPI

TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
BIN_DIR := bin
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/controller-gen)

$(CONTROLLER_GEN): $(TOOLS_DIR)/go.mod # Build controller-gen from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/controller-gen sigs.k8s.io/controller-tools/cmd/controller-gen

.PHONY: generate-go-capi
generate-capi: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) \
		object:headerFile=./hack/boilerplate.generatego.txt \
		paths=./pkg/capi/api/...
	$(CONTROLLER_GEN) \
		paths=./pkg/capi/api/... \
		crd:crdVersions=v1 \
		output:crd:dir=./pkg/capi/crds/

.PHONY: build-capi
build-capi:
	go build -mod=vendor -o bin/capiManager github.com/openshift-hive/hypershift-installer/cmd/capi