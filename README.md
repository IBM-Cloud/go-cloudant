[![Build Status](https://travis-ci.org/IBM-Cloud/go-cloudant.svg?branch=master)](https://travis-ci.org/IBM-Cloud/go-cloudant)

# go-cloudant

go-cloudant is a Cloudant DB client written in Go. It takes advantage
of the go-couchdb client and add Index and Search into
it to ease the usage of the Cloudant DB. Also, it tries to simplify the
use of couchdb library by adding more native structs.

The go-couchdb credits go to `fjl/go-couchdb` and `timjacobi/go-couchdb`

## Usage

    import "github.com/IBM-Cloud/go-cloudant"

For detailed usage, check cloudant_test.go

## Test

    make test

All methods should be covered by tests, and the Makefile will also check
the format of the code, so try to use `make` before the commit.

## Contribution

To make contributions, please add tests to the methods or functionality
you've added.
