PROJECT="cpds-detector"
OUTPUT_DIR=./out

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

ifeq (${GOFLAGS},)
	# go build with vendor by default.
	export GOFLAGS=-mod=vendor
endif

.PHONY:  all cpds-detector clean help
default: all
all: cpds-detector

cpds-detector: 
	@echo "Building $(PROJECT)"
	@go build -ldflags=$(BUILD_LDFLAGS) -o $(OUTPUT_DIR)/$(PROJECT) ./cmd/$(PROJECT)/main.go

clean:
	@rm -rf $(OUTPUT_DIR)

help:
	@echo "make help: get help"
	@echo "make cpds-detector: compile cpds-detector binaries"
	@echo "make clean: clean up binaries"