
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifneq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

install-deepcopy-gen:
ifeq (, $(shell which deepcopy-gen))
	@{ \
	set -e ;\
	LICENSE_TMP_DIR=$$(mktemp -d) ;\
	cd $$LICENSE_TMP_DIR ;\
	go mod init tmp ;\
	go get -v k8s.io/code-generator/cmd/deepcopy-gen ;\
	rm -rf $$LICENSE_TMP_DIR ;\
	}
DEEPCOPY_BIN=$(GOBIN)/deepcopy-gen
else
DEEPCOPY_BIN=$(shell which deepcopy-gen)
endif



HEAD_FILE := hack/boilerplate.go.txt
INPUT_DIR := github.com/labring/sealvm/types/api
deepcopy:install-deepcopy-gen
	$(DEEPCOPY_BIN) \
      --input-dirs="$(INPUT_DIR)/v1" \
      -O zz_generated.deepcopy   \
      --go-header-file "$(HEAD_FILE)" \
      --output-base "${GOPATH}/src"
