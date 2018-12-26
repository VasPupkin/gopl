# variables ###################################################################

BIN := poster
RELEASE := 0.2
ALT_NAME := "pukan_v_ogne"
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
API_KEY := ""
LDFLAGS := -ldflags="-X main.version=${RELEASE} -X main.buildtime=${BUILD_TIME} -X main.altname=${ALT_NAME} -X gopl/ch4/poster/omdbapi.apiKey=${API_KEY} -s -w"

# commands ####################################################################

.DEFAULT_GOAL: build

.PHONY: build
build: clean
	@go build -v ${LDFLAGS} -o ${BIN} 

.PHONY: clean
clean:
	@if [ -e ${BIN} ]; then rm -f ${BIN}; fi
