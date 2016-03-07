GOOPTS=GO15VENDOREXPERIMENT=1
GO=${GOOPTS} go

.PHONY: release

tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} build -o tael

test: $(wildcard src/tael/*.go src/tael/**/*.go)
	${GO} test tael

clean:
	rm -rf ./bin ./pkg ${RELEASE_DIR}
