# hello-world-service [![Build Status](https://travis-ci.com/ryandens/hello-world-service.svg?token=zrk2qgs8pLCF6veGNV1S&branch=master)](https://travis-ci.com/ryandens/hello-world-service)
Simple GoLang HTTP server which supports TLS.

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


## Project post-mortem
I fell short on the requirements a bit. I used basic authentication, which is perfectly secure,
but I did not create the forms necessary for it to be usable in a browser. To use the app,
one must always provide the authorization header. Sample `curl` commands are in the [Makefile](Makefile)

I definitely bit off more than I could chew a bit, and got hung up a bit on some implementation
details. In retrospect, perhaps I should have not used TLS and left that out of scope of the
project. Regardless, I did enjoy toying with it

I did not have the time to read super thoroughly on Go best practices. The biggest security
pitfall in this project is me, because I am not super familiar with Go's APIs. In addition, this
project is not up to my style or testing standards. I was able to skim the "Effective Go" blog,
which reminded me of my favorite book "Effective Java". I found it hard to keep the design principles
I've learned in technical books in mind while trying to finish this small project in a reasonable
amount of time and learning to do things "the go way". I also wasn't familiar with the go approach
for integration testing, so I rely on using curl.

All in all, I had loads of fun doing this project. This is something to iterate on as I become
more familiar with Go, and I'm sure I'll have lots of opinions on how wrong I was to do
X, Y, and Z in not too long.
