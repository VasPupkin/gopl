# variables ###################################################################

BIN := xkcd
RELEASE := 0.6
ALT_NAME := "gatos_de_larga_vida"
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags="-X main.version=${RELEASE} -X main.buildtime=${BUILD_TIME} -X main.altname=${ALT_NAME} -s -w"

# commands ####################################################################

.DEFAULT_GOAL: build

.PHONY: build
build: clean
	@go build -v ${LDFLAGS} -o ./bin/${BIN} 

.PHONY: clean
clean:
	@if [ -e ./bin/${BIN} ]; then rm -f ./bin/${BIN}; fi
