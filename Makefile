PHONY: clean devdeps deps test build

clean:
	rm -rf server
	rm -rf OPATH/
	cp test_users_tmp.csv test_users.csv


deps:
	openssl genrsa -out private.key 2048
	openssl req -new -x509 -sha256 -key private.key -out public.crt -subj "/C=US/ST=Ryan Dens/L=Baltimore/O=Ryan Dens Enterprises/OU=Dev/CN=localhost"
	go get golang.org/x/crypto/bcrypt

devdeps:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $GOPATH/bin latest

test:
	go test
	OPATH/bin/gosec ./...

build:
	go build -o server


curlcommands:
	curl --user ryandens:123ABD --request POST --cacert public.crt https://localhost:8080/signup
	curl --user ryandens:123ABD --request GET --cacert public.crt https://localhost:8080/hello
	curl --user ryandens:123ABD --request POST --cacert public.crt --data 'ryancdens' https://localhost:8080/update-name

badcurlcommands:
	curl --user ryandensaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa:123ABD --request POST --cacert public.crt https://localhost:8080/signup
	curl --user ryandens:123ABDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa --request POST --cacert public.crt https://localhost:8080/signup
	curl --user alert%281%29:123ABD --request POST --cacert public.crt https://localhost:8080/signup
