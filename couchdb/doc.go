// Tideland Go CouchDB Client - CouchDB
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// The Tideland Go CouchDB Client provides a very powerful as well as
// convenient client for the CouchDB database.
//
// A connection to the database or at least a server can be established
// by calling
//
//     cdb := couchdb.Open(cfg)
//
// The expected configuration is
//
//    {etc
//        {hostname <hostname||localhost>}
//        {port <port||5984>}
//        {database <database||default>}
//    }
//
// If any of the values isn't defined the default values above are taken.
//
// The currently supported operations are the listing, creation, or deleting
// of databases, the listing of all design document and document, and the
// creation, reading, updating, and deleting of document.
package couchdb

// EOF
