#! /usr/bin/make

# Name of this app, used as basename for almost everything.
NAME=mjbot

# Dependencies that are used by the build & test processes.
# These need to be installed in the global Go env and not in the vendor sub-tree.
DEPEND=golang.org/x/tools/cmd/cover github.com/onsi/ginkgo/ginkgo \
       github.com/onsi/gomega github.com/rlmcpherson/s3gof3r/gof3r \
       github.com/Masterminds/glide github.com/golang/lint/golint

# Eg: 2026-05-21 19:17:44
DATE=$(shell date '+%F %T')

# We might not need this anymore.
GO15VENDOREXPERIMENT=1
export GO15VENDOREXPERIMENT

# govvv allows version embedding at compile time.
VFLAG=$(shell govvv -flags -pkg $(shell go list ./src/security))

.PHONY: clean default depend lint release

# The default target builds a binary in the top-level dir for whatever the local OS is.
# It does not depend on 'depend' 'cause it's a pain to have that run every time we hit 'make'.
# Instead we need to 'make depend' manually once during the initial setup.
default: $(NAME)
$(NAME): $(shell find . -name \*.go)
	go build -ldflags "$(VFLAG)" -o $(NAME) src/main.go

# the standard build produces a "local" executable, a linux tgz, and a darwin (macos) tgz
# uncomment and join the windows zip if you need it
build: $(NAME)
# build/$(NAME)-linux-amd64.tgz build/$(NAME)-darwin-amd64.tgz
# build/$(NAME)-linux-arm.tgz build/$(NAME)-windows-amd64.zip

# create a tgz with the binary and any artifacts that are necessary
# note the hack to allow for various GOOS & GOARCH combos, sigh
build/$(NAME)-%.tgz: $(NAME)
	rm -rf build/$(NAME)
	mkdir -p build/$(NAME)
	tgt=$*; GOOS=$${tgt%-*} GOARCH=$${tgt#*-} go build -ldflags "$(VFLAG)" -o build/$(NAME)/$(NAME) .
	chmod +x build/$(NAME)/$(NAME)
	cp README.md build/$(NAME)/
	tar -zcf $@ -C build ./$(NAME)
	rm -r build/$(NAME)

build/$(NAME)-%.zip: $(NAME)
	touch $@

release: $(NAME)
	strip "./$(NAME)"
	upx --ultra-brute "./$(NAME)"

# Installing build dependencies. You will need to run this once manually when you clone the repo.
depend:
	go get -v $(DEPEND)
	glide install

clean:
	rm -rf build .vendor/pkg

# Run gofmt and complain if a file is out of compliance
# Run go vet and similarly complain if there are issues
# Run go lint and complain if there are issues
lint:
	@if gofmt -l . | egrep -v ^vendor/ | grep .go; then \
	  echo "^- Repo contains improperly formatted go files; run gofmt -w *.go" && exit 1; \
	  else echo "All .go files formatted correctly"; fi
	#go tool vet -v -composites=false *.go
	#go tool vet -v -composites=false **/*.go
	for pkg in $$(go list ./... |grep -v /vendor/); do golint $$pkg; done

# upload assumes you have AWS_ACCESS_KEY_ID and AWS_SECRET_KEY env variables set,
# which happens in the .travis.yml for CI. Yup, that means you can't run it from your laptop,
# which is a good thing!
#upload:
#	@which gof3r >/dev/null || (echo 'Please "go get github.com/rlmcpherson/s3gof3r/gof3r"'; false)
#	(cd build; set -ex; \
#	  for f in *.tgz; do \
#	    gof3r put --no-md5 --acl=$(ACL) -b ${BUCKET} -k rsbin/$(NAME)/$(TRAVIS_COMMIT)/$$f <$$f; \
#	    if [ "$(TRAVIS_PULL_REQUEST)" = "false" ]; then \
#	      gof3r put --no-md5 --acl=$(ACL) -b ${BUCKET} -k rsbin/$(NAME)/$(TRAVIS_BRANCH)/$$f <$$f; \
#	    fi; \
#	  done)
