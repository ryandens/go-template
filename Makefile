PHONY: clean deps test build

clean:
	rm -rf server
	rm -rf OPATH/

deps:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $GOPATH/bin latest

test:
	go test
	OPATH/bin/gosec ./...

build:
	go build -o server
