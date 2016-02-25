GOOPTS=GOPATH=$(shell pwd) GO15VENDOREXPERIMENT=1
GO=${GOOPTS} go
OSX=GOOS=darwin GOARCH=amd64 ${GO}
LINUX=GOOS=linux GOARCH=amd64 ${GO}

.PHONY: release

./bin/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} install tael

release: $(wildcard src/tael/*.go src/tael/**/*.go)
	mkdir -p release/osx release/linux
	${OSX} build -o ./release/osx/tael tael
	${LINUX} build -o ./release/linux/tael tael

clean:
	rm -rf ./bin ./pkg ./release
