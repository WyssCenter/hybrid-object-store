COMMIT=$(shell git rev-list HEAD --max-count=1)
export COMMIT

ifndef STAGE
STAGE = dev
endif

local:
	GOOS=linux go build

image:
	docker build --build-arg STAGE=$(STAGE) -t hoss-core:$(COMMIT) -f ./Dockerfile ../

test:
	echo tests
