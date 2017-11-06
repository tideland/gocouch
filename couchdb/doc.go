// Tideland GoCouch - CouchDB
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package couchdb provides the powerful as well as convenient
// Tideland Go CouchDB Client as client for the CouchDB database.
//
// A connection to the database or at least a server can be established
// by calling
//
//     cdb, err := couchdb.Open(cfg)
//
// The expected configuration is
//
//    {etc
//        {hostname <hostname||localhost>}
//        {port <port||5984>}
//        {database <database||default>}
//        {debug-logging <true/false||false>}
//    }
//
// If any of the values isn't defined the default values above are taken.
// Instead of splitting a larger configuration it's also possible to use
//
//    cdb, err := couchdb.OpenPath(cfg, "path/to/couchdb/config")
//
// In case of not using the etc configuration package there's the
// little helper to create a configuration by calling
//
//    cfg := couchdb.Configure(hostname, port, database)
//
// The supported operations are the listing, creation, and deleting of
// databases, the listing of all design documents and data documents, and
// the creation, reading, updating, and deleting of documents.
package couchdb

// EOF
