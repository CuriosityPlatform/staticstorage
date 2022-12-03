# enable buildkit
export DOCKER_BUILDKIT=1

IMAGE_TAG?=master
export IMAGE=vadimmakerov/staticstorage:${IMAGE_TAG}

# Some commands can be duplicated locally for better local development
# Set ENV to local to force DUPLICATE COMMANDS for local environment
ENV?=
ifeq (${ENV},)
	# Automatically ENV set up to local only if there is GO compiler
	ifeq ($(shell which go > /dev/null && echo "1" || echo "0"), 1)
	ENV=local
	endif
endif

# Builder cache ttl in hours
CACHE_TTL=1

all: build test check

.PHONY: build
build: modules build-image

.PHONY: build-image
build-image:
	@docker buildx build . \
	--target make-staticstorage-image \
	--output type=image,name=${IMAGE}

.PHONY: modules
modules:
	@docker buildx build . \
 	--target go-mod-tidy \
	--output .

ifeq (${ENV}, local)
	go mod download
endif

.PHONY: test
test:
	@docker buildx build . \
	--target test

.PHONY: check
check:
	@docker buildx build . --target lint

.PHONY: cache-clear
cache-clear: ## Clear the builder cache
	@docker builder prune --force --filter type=exec.cachemount --filter=unused-for=${CACHE_TTL}h
