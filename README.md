# go-template [![Build Status](https://travis-ci.com/ryandens/go-template.svg?token=zrk2qgs8pLCF6veGNV1S&branch=master)](https://travis-ci.com/ryandens/go-template)
Template repository which has a simple HTTP server using TLS.

## Development
To build and test this project, simply run `make` in the root project directory. See the
[Makefile](Makefile) for specific details but this command generally
- cleans local artifacts not checked into the repository
- installs/creates dependencies
- runs the tests
- builds an executable named `server`. When run, the server responds to HTTP requests on
localhost port `8080`. Note that the default https port is typically `443`, however 
depending on user permissions, a user may not be able to use port `443` without assuming the
super-user role.

## TLS
The [Makefile](Makefile) in this repository uses the [OpenSSL](https://www.openssl.org/) 
command line utility to create a private key and a self-signed certificate for allowing secure 
transmission of data from the client to the server. In a "real" environment, using a widely
trusted Certificate Authority such as [Let's Encrypt](https://letsencrypt.org/) would be 
much more sensible. Most HTTP clients do not trust self-signed certificates, for good reason.
However, using a self-signed certificate for local development is perfectly sensible. 

I used `curl` during my local development of this project. There are many ways to bypass 
the HTTP client's checks on the certificate, many of which result in no benefits, from the 
clients perspective, of using TLS at all. However, since we have the server's certificate
locally, we can simply tell the HTTP client that if the server's certificate matches the one
we have locally, it can be trusted. This is exactly what is done when using a certificate from
a Certificate Authority, except that the trusted certificates are shipped with the operating
system. To make an HTTP request to the server using curl simply run
```shell script
$ curl --request GET --cacert public.crt https://localhost:8080/
```

## Application Security
This project uses [gosec](https://github.com/securego/gosec),
which scans the Go AST of the project in the Travis CI pipeline. This scan occurs in the 
`test` section of the [Makefile](Makefile) and runs with every `make` command in the 
root directory of this project.
