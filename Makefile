PHONY: clean devdeps deps test build

clean:
	rm -rf server
	rm -rf OPATH/

deps:
	openssl genrsa -out private.key 2048
	openssl req -new -x509 -sha256 -key private.key -out public.crt -subj "/C=US/ST=Ryan Dens/L=Baltimore/O=Ryan Dens Enterprises/OU=Dev/CN=localhost"

devdeps:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $GOPATH/bin latest

test:
	go test
	OPATH/bin/gosec ./...

build:
	go build -o server
