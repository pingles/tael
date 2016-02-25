GOOPTS=GOPATH=$(shell pwd) GO15VENDOREXPERIMENT=1
GO=${GOOPTS} go
OSX=GOOS=darwin GOARCH=amd64 ${GO}
LINUX=GOOS=linux GOARCH=amd64 ${GO}
RELEASE_DIR?=./release

.PHONY: release

./bin/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} install tael

release: $(wildcard src/tael/*.go src/tael/**/*.go)
	mkdir -p ${RELEASE_DIR}/osx ${RELEASE_DIR}/linux
	${OSX} build -o ${RELEASE_DIR}/osx/tael tael
	${LINUX} build -o ${RELEASE_DIR}/linux/tael tael

test: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} test tael

clean:
	rm -rf ./bin ./pkg ${RELEASE_DIR}
