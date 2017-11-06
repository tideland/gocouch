# Tideland GoCouch

## Description

*Tideland GoCouch* provides a convenient and powerful access to
CouchDB databases.

I hope you like it. ;)

[![GitHub release](https://img.shields.io/github/release/tideland/gocouch.svg)](https://github.com/tideland/gocouch)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/gocouch/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/tideland/gocouch?status.svg)](https://godoc.org/github.com/tideland/gocouch)
[![Sourcegraph](https://sourcegraph.com/github.com/tideland/gocouch/-/badge.svg)](https://sourcegraph.com/github.com/tideland/gocouch?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/gocouch)](https://goreportcard.com/report/github.com/tideland/gocouch)

## Version

Version 0.7.1

## Packages

### CouchDB

Package `couchdb` is the client for the access of the CouchDB. It provides the
standard functionality to create databases as well as read, write, and delete
documents.

### Views

Package `views` allows to request CouchDB views. Right now these have to be
created using the design documents in package `couchdb`. Future releases will
be able to create, modify, and delete them direct from this package too.

### Find

Package `find` helps to create *Mango* queries the Go way. Typically they have
a very special JSON notation. Searches will then be executed using the `Find()`
function. Addtional parameters help to restrict the result set to individual
fields, to filter the result, or to paginate it.

### Changes

Package `changes` allow to retrieve the changes made in a datebase in time order.

### Security

Package `security` helps with user administration and authentication for CouchDB.

### Startup

Package `startup` provides a simple mechanism for a clean startup and maintenance
of CouchDB databases including database versioning.

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland)

## License

*Tideland Go CouchDB Client* is distributed under the terms of the BSD 3-Clause license.
