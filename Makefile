GOPATH := ${PWD}
export GOPATH

# Opted to set up a bit differently and not use glide this time.
BASE="src/urlcheck"
SRC_FOLDERS=""

default: deps build

fmt:
	for FOLDER in ${SRC_FOLDERS} ; do \
	    cd ${BASE}/${FOLDER} ; \
	    go fmt ;\
	done

vet:
	for FOLDER in ${SRC_FOLDERS} ; do \
	    cd ${BASE}/${FOLDER} ; \
	    go vet ;\
	done

deps:
	for FOLDER in ${SRC_FOLDERS} ; do \
	    cd ${BASE}/${FOLDER} ; \
	    go get ;\
	done

build:
	mkdir -p bin/
	for FOLDER in ${SRC_FOLDERS} ; do \
		cd ${BASE}/${FOLDER} ; \
		go build ;\
	done

	mv ${BASE}/urlcheck bin/

run:
	docker-compose up --build
