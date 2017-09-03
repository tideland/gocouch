// Tideland Go CouchDB Client - Startup
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package startup of the Tideland Go CouchDB Client provides a simple
// mechanism for a clean startup and maintenance of CouchDB databases.
//
// The major function is
//
//    err := startup.Run(cdb, stepA, stepB, stepC)
//
// Here cdb is the opened CouchDB and the steps is a variadic number of
// functions having the signature
//
//    func(cdb couchdb.CouchDB, v version.Version) (version.Version, error)
//
// When Run() is called it first checks if the database already exists and
// if needed creates it. Then it adds a version document with the version
// 0.0.0. Now each step is called in order with the current version. It
// can check if it has to modify the database (e.g. add design and other
// document, add fields to existing documents, transform documents, etc.),
// perform those changes, and return the new version. So the version
// document will be updated and the next step performed.
package startup

// EOF
