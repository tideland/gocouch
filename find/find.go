// Tideland Go CouchDB Client - Find
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// API
//--------------------

// Find returns access to the found results.
func Find(cdb couchdb.CouchDB, sel Selector) ResultSet {
	return nil
}

//--------------------
// FIND RESULT SET
//--------------------

// Processor is a function processing the content of a found document.
type Processor func(id, sequence string, deleted bool, revisions []string, document couchdb.Unmarshable) error

// ResultSet contains the result set of a find call.
type ResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// Len returns the number of changes.
	Len() int

	// Do iterates over the results of a ResultSet and
	// processes the content.
	Do(process Processor) error
}

// EOF
