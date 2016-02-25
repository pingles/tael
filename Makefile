GOOPTS=GOPATH=$(shell pwd) GO15VENDOREXPERIMENT=1
GO=${GOOPTS} go
OSX=GOOS=darwin GOARCH=amd64 ${GO}
LINUX=GOOS=linux GOARCH=amd64 ${GO}
RELEASE_DIR?=./release

.PHONY: release

./bin/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} install tael

${RELEASE_DIR}/osx/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	mkdir -p ${RELEASE_DIR}/osx
	${OSX} build -o ${RELEASE_DIR}/osx/tael tael
	tar czf ${RELEASE_DIR}/osx/tael.tar.gz ${RELEASE_DIR}/osx/tael

${RELEASE_DIR}/linux/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	mkdir -p ${RELEASE_DIR}/linux
	${LINUX} build -o ${RELEASE_DIR}/linux/tael tael
	tar czf ${RELEASE_DIR}/linux/tael.tar.gz ${RELEASE_DIR}/linux/tael

release: ${RELEASE_DIR}/linux/tael ${RELEASE_DIR}/osx/tael

test: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} test tael

clean:
	rm -rf ./bin ./pkg ${RELEASE_DIR}
