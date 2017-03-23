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

// changesResultSet implements the ChangesResultSet interface.
type changesResultSet struct {
	rs ResultSet
}

// newChangesResultSet returns a ChangesResultSet.
func newChangesResultSet(rs ResultSet) ChangesResultSet {
	crs := &changesResultSet{
		rs: rs,
	}
	return crs
}

// IsOK implements the ChangesResultSet interface.
func (crs *changesResultSet) IsOK() bool {
	return crs.rs.IsOK()
}

// StatusCode implements the ChangesResultSet interface.
func (crs *changesResultSet) StatusCode() int {
	return crs.rs.StatusCode()
}

// Error implements the ChangesResultSet interface.
func (crs *changesResultSet) Error() error {
	return crs.rs.Error()
}

// EOF
