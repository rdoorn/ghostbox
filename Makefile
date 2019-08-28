#
# Makefile
# @author Ronald Doorn <rdoorn@schubergphilis.com>
#

.PHONY: update clean build build-all run package deploy test authors dist

export PATH := $(PATH):$(GOPATH)/bin

NAME := ixxid
VERSION := $(shell [ -f .version ] && cat .version || echo "pipeline-test")
LASTCOMMIT := $(shell git rev-parse --verify HEAD)
#BUILD := $(shell cat tools/rpm/BUILDNR)
BUILD := "1"
LDFLAGS := "-X main.version=$(VERSION) -X main.versionBuild=$(BUILD) -X main.versionSha=$(LASTCOMMIT)"
PENDINGCOMMIT := $(shell git diff-files --quiet --ignore-submodules && echo 0 || echo 1)
# config points to oudside the repo
TESTPARAMS := serve --config-file ~/ixxi.yaml 
GODIRS := $(shell go list -f '{{.Dir}}' ./...)

default: build

clean:
	@echo Cleaning up...
	@rm -f build
	@echo Done.

rice:
	@echo Merging static content...
	#@which rice >/dev/null; if [ $$? -eq 1 ]; then \
		#go get github.com/GeertJohan/go.rice/rice; \
	#fi;
	#cd internal/core && rice embed-go
	@echo Done

builddir:
	@mkdir -p ./build/osx/
	@mkdir -p ./build/linux/
	@mkdir -p ./build/packages/$(NAME)

osx: builddir rice
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(NAME) -ldflags $(LDFLAGS) ./cmd/$(NAME)
	@echo Done.

osx-fast: builddir
	@echo Building OSX skipping rice...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(NAME) -ldflags $(LDFLAGS) ./cmd/$(NAME)
	@echo Done.

osx-race: builddir rice
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -race -v -o ./build/osx/$(NAME) -ldflags $(LDFLAGS) ./cmd/$(NAME)
	@echo Done.

osx-static:
	@echo Building OSX...
	GOOS=darwin GOARCH=amd64 go build -v -o ./build/osx/$(NAME) -ldflags '-s -w --extldflags "-static”  $(LDFLAGS)' ./cmd/$(NAME)
	@echo Done.

linux: builddir rice
	@echo Building Linux...
	GOOS=linux GOARCH=amd64 go build -v -o ./build/linux/$(NAME) -ldflags '-s -w --extldflags "-static”  $(LDFLAGS)' ./cmd/$(NAME)
	@echo Done.

build: osx linux

makeconfig:
	#@echo Making config...
	#cat ./test/$(NAME)-template.toml | sed -e 's/%LOCALIP%/$(LOCALIP)/g' > ./test/$(NAME).toml
	#cat ./test/$(NAME).toml | sed -e 's/port = 9/port = 10/g' -e 's/localhost1/localhost3/' -e 's/localhost2/localhost1/' -e 's/localhost3/localhost2/' -e 's/127.0.0.1:9000/127.0.0.1:8000/' -e 's/127.0.0.1:10000/127.0.0.1:9000/' -e 's/127.0.0.1:8000/127.0.0.1:10000/' -e 's/preference = 0#1/preference = 1/' -e 's/preference = 1#0/preference = 0/' -e 's/15353/25353/' -e 's/ip = "127.0.0.1"/ip = "127.0.0.2"/' > ./test/$(NAME)-secondary.toml

run: osx makeconfig
	./build/osx/$(NAME) $(TESTPARAMS)

run-linux: linux makeconfig
	./build/linux/$(NAME) $(TESTPARAMS)

run-race: osx-race makeconfig
	./build/osx/$(NAME) $(TESTPARAMS)

sudo-run: osx
	sudo ./build/osx/$(NAME) $(TESTPARAMS)

test:
	go test -v ./...
	go test -v ./... --race --short
	go vet ./...

coverage: ## Shows coverage
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
	    go get -u golang.org/x/tools/cmd/cover; \
	fi
	for pkg in $(GODIRS); do go test -coverprofile=go$${pkg//\//-}.cover $$pkg ;done
	echo "mode: set" > c.out
	grep -h -v "^mode:" ./*.cover >> c.out
	rm -f *.cover

coverage-upload:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
	chmod +x ./cc-test-reporter
	./cc-test-reporter after-build
	rm -f ./cc-test-reporter

prep_package:
	gem install fpm

committed:
ifeq ($(PENDINGCOMMIT), 1)
	   $(error You have a pending commit, please commit your code before making a package $(PENDINGCOMMIT))
endif

linux-package: builddir linux # committed
	cp -a ./tools/rpm/$(NAME)/* ./build/packages/$(NAME)/
	cp ./build/linux/$(NAME) ./build/packages/$(NAME)/usr/sbin/
	fpm -s dir -t rpm -C ./build/packages/$(NAME) --name $(NAME) --rpm-os linux --version $(VERSION) --iteration $(BUILD) --exclude "*/.keepme"
	rm -rf ./build/packages/$(NAME)/
	mv $(NAME)-$(VERSION)*.rpm build/packages/

docker-scratch:
	if [ -a /System/Library/Keychains/SystemRootCertificates.keychain ] ; \
	then \
		security find-certificate /System/Library/Keychains/SystemRootCertificates.keychain > build/docker/ca-certificates.crt; \
	fi;
	if [ -a /etc/ssl/certs/ca-certificates.crt ] ; \
	then \
		cp /etc/ssl/certs/ca-certificates.crt build/docker/ca-certificates.crt; \
	fi;
	docker build -t $(NAME)-scratch -f build/docker/Dockerfile.scratch .

deps: ## Updates the vendored Go dependencies
	@dep ensure -v

updatedeps: ## Updates the vendored Go dependencies
	@dep ensure -update

get-version:
	git describe --tags --always > .version
	echo "path: ${PWD}"
	cat .version

ci-package:
	./tools/ci-package.sh

#authors:
#	@git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort | uniq -c | sort -nr | sed "s/^ *[0-9]* //g" > AUTHORS
#	@cat AUTHORS
#
