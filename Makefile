GOPATH := ${PWD}
export GOPATH

BASE=urlcheck

PACKAGES=${BASE}          \
         ${BASE}/data     \
         ${BASE}/models   \
		 ${BASE}/services \
		 ${BASE}/utils

default: deps build

test:
	@go test -v ${PACKAGES}

fmt:
	@go fmt ${PACKAGES}

vet:
	@go vet ${PACKAGES}

deps:
	@go get ${PACKAGES}

build:
	mkdir -p bin/
	cd src/${BASE} && go build
	mv src/${BASE}/${BASE} bin/

run: deps build
	docker-compose up --build
