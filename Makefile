.PHONY: all
all: build test


setup:
	mkdir -p $GOPATH/bin
	if which dep &> /dev/null ; then go get -u github.com/golang/dep/cmd/dep ; fi


compile:
	mkdir -p out

	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o out/ogi main.go
	echo "compiled ogi main"

	echo "done."

build: compile

build-test-plugins:
	export THIS_DIR=$(pwd)
	cd tests/consumer ; \
		GO111MODULE=on go build -o "../consumer.so" -buildmode=plugin . ; cd -
	cd tests/consumer-bad ; \
		GO111MODULE=on go build -o "../consumer-bad.so" -buildmode=plugin . ; cd -
	cd tests/transformer ; \
		GO111MODULE=on go build -o "../transformer.so" -buildmode=plugin . ; cd -
	cd tests/transformer-bad ; \
		GO111MODULE=on go build -o "../transformer-bad.so" -buildmode=plugin . ; cd -
	cd tests/producer ; \
		GO111MODULE=on go build -o "../producer.so" -buildmode=plugin . ; cd -
	cd tests/producer-bad ; \
		GO111MODULE=on go build -o "../producer-bad.so" -buildmode=plugin . ; cd -

test: build-test-plugins
	CONSUMER_PLUGIN_PATH=$(shell pwd)/tests/consumer.so \
	PRODUCER_PLUGIN_PATH=$(shell pwd)/tests/producer.so \
	TRANSFORMER_PLUGIN_PATH=$(shell pwd)/tests/transformer.so \
	GO111MODULE=on go test -gcflags=-l github.com/OpenChaos/ogi/consumer github.com/OpenChaos/ogi/transformer github.com/OpenChaos/ogi/producer
