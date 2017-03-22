// Tideland Go CouchDB Client - CouchDB - Changes
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

//--------------------
// CHANGES RESULT SET
//--------------------

// ChangesResultSet contains the result set of a change.
type ChangesResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error
}

// EOF
