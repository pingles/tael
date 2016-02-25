./bin/tael: $(wildcard src/tael/*.go src/tael/**/*.go)
	GOPATH=$(shell pwd) GO15VENDOREXPERIMENT=1 go install tael
