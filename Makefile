GOOPTS=GO15VENDOREXPERIMENT=1
GO=${GOOPTS} go

.PHONY: release

tael: $(shell find . -name '*.go')
	${GO} build -o tael ./cmd/tael

test: $(shell find . -name '*.go')
	${GO} test tael

clean:
	rm -rf ./bin ./pkg ${RELEASE_DIR}
