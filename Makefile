# Copyright(c) 2018 DSiSc Group. All Rights Reserved.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Version of dependencies
VER_DSiSc_craft=0
VER_DSiSc_blockchain=0
VER_DSiSc_statedb-NG=0
VER_DSiSc_evm-NG=0
VER_DSiSc_txpool=0

VERSION=$(shell grep "const Version" version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

.PHONY: default help all build test unit-test devenv gotools clean coverage

default: all

help:
	@echo 'Management commands for DSiSc/evm-NG:'
	@echo
	@echo 'Usage:'
	@echo '    make lint            Check code style.'
	@echo '    make spelling        Check code spelling.'
	@echo '    make fmt             Check code formatting.'
	@echo '    make static-check    Static code check: style & spelling & formatting.'
	@echo '    make build           Compile the project.'
	@echo '    make vet             Examine source code and reports suspicious constructs.'
	@echo '    make unit-test       Run unit tests with coverage report.'
	@echo '    make test            Run unit tests with coverage report.'
	@echo '    make devenv          Prepare devenv for test or build.'
	@echo '    make fetch-deps      Run govendor fetch for deps.'
	@echo '    make gotools         Prepare go tools depended.'
	@echo '    make clean           Clean the directory tree.'
	@echo

all: static-check build test

fmt:
	gofmt -d -l .

spelling:
	bash scripts/check_spelling.sh

lint:
	@echo "Check code style..."
	golint `go list ./...`

static-check: fmt spelling lint

build:
	@echo "building evm-NG ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -v -ldflags "-X github.com/DSiSc/evm-NG/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/DSiSc/evm-NG/version.BuildDate=${BUILD_DATE}" ./...

vet:
	@echo "Examine source code and reports suspicious constructs..."
	go vet `go list ./...`

unit-test:
	@echo "Run unit tests without coverage report..."
	go test -v -count=1 -race ./...

coverage:
	@echo "Run unit tests with coverage report..."
	bash scripts/unit_test_cov.sh

test: vet unit-test

get-tools:
        # official tools
	go get -u golang.org/x/lint/golint
	@# go get -u golang.org/x/tools/cmd/gotype
	@# go get -u golang.org/x/tools/cmd/goimports
	@# go get -u golang.org/x/tools/cmd/godoc
	@# go get -u golang.org/x/tools/cmd/gorename
	@# go get -u golang.org/x/tools/cmd/gomvpkg

        # thirdparty tools
	@# go get -u github.com/kardianos/govendor
	go get -u github.com/stretchr/testify
	@# go get -u github.com/axw/gocov/...
	@# go get -u github.com/client9/misspell/cmd/misspell

fetch-deps: get-tools
	@echo "Run go get to fetch dependencies ..."
	go get -u gopkg.in/DSiSc/craft.v${VER_DSiSc_craft}
	go get -u gopkg.in/DSiSc/blockchain.v${VER_DSiSc_blockchain}
	go get -u gopkg.in/DSiSc/txpool.v${VER_DSiSc_txpool}
	go get -u gopkg.in/DSiSc/evm-NG.v${VER_DSiSc_evm-NG}
	go get -u gopkg.in/DSiSc/statedb-NG.v${VER_DSiSc_statedb-NG}

## tools & deps
devenv: get-tools fetch-deps

