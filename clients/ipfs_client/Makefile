# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

TOP = ../..
include $(TOP)/build.mk

build:
	$(GO) build $(GOTAGS) -o ipfsclient

DOCKER_BUILD_OPTS += --build-arg FPC_VERSION=$(FPC_VERSION)
DOCKER_BUILD_OPTS += --no-cache

docker:
	$(DOCKER) build $(DOCKER_BUILD_OPTS) \
    	-t fpc/ipfsclient:$(FPC_VERSION) \
    	-f Dockerfile \
    	.
