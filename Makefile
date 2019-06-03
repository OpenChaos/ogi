.PHONY: all
all: build test


setup:
	mkdir -p $GOPATH/bin
	if which dep &> /dev/null ; then go get -u github.com/golang/dep/cmd/dep ; fi


compile:
	mkdir -p out

	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o out/ogi main.go
	echo "compiled ogi main"

	cd plugin-examples/transformer/message_logs ; \
		GO111MODULE=on go build -o "../../../out/transformer-message-log.so" -buildmode=plugin . ; \
		cd - ; echo "compiled transformer.message_logs plugin"
	cd plugin-examples/transformer/os_path_exists ; \
		GO111MODULE=on go build -o "../../../out/transformer-os-path-exists.so" -buildmode=plugin . ; \
		cd - ; echo "compiled transformer.os_path_exists plugin"

	cd plugin-examples/producer/echo; \
		GO111MODULE=on go build -o "../../../out/producer-echo.so" -buildmode=plugin . ; \
		cd - ; echo "compiled producer.echo plugin"
	cd plugin-examples/producer/filedump; \
		GO111MODULE=on go build -o "../../../out/producer-filedump.so" -buildmode=plugin . ; \
		cd - ; echo "compiled producer.filedump plugin"

	cd plugin-examples/consumer/gcp_stackdriver_logs; \
		GO111MODULE=on go build -o "../../../out/consumer-gcp-stackdriver-logs.so" -buildmode=plugin . ; \
		cd - ; echo "compiled consumer.gcp_stackdriver_logs plugin"
	cd plugin-examples/consumer/file_line_by_line; \
		GO111MODULE=on go build -o "../../../out/consumer-file-line-by-line.so" -buildmode=plugin . ; \
		cd - ; echo "compiled consumer.file_line_by_line plugin"

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
	go test -gcflags=-l github.com/OpenChaos/ogi/consumer github.com/OpenChaos/ogi/transformer github.com/OpenChaos/ogi/producer
