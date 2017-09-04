// Tideland Go CouchDB Client - Changes
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package changes

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// API
//--------------------

// Changes returns access to the changes of the database.
func Changes(cdb couchdb.CouchDB, params ...couchdb.Parameter) ResultSet {
	rs := cdb.GetOrPost(cdb.DatabasePath("_changes"), nil, params...)
	return newResultSet(rs)
}

//--------------------
// CHANGES RESULT SET
//--------------------

// Processor is a function processing the content of a changed document.
type Processor func(id, sequence string, deleted bool, revisions []string, document couchdb.Unmarshable) error

// ResultSet contains the result set of a change.
type ResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// LastSequence returns the sequence ID of the last change.
	LastSequence() string

	// Pending returns the number of pending changes if the
	// query has been limited.
	Pending() int

	// Len returns the number of changes.
	Len() int

	// Do iterates over the results of a ResultSet and
	// processes the content.
	Do(process Processor) error
}

// resultSet implements the ResultSet interface.
type resultSet struct {
	rs      couchdb.ResultSet
	changes *couchdbChanges
}

// newResultSet returns a ResultSet.
func newResultSet(rs couchdb.ResultSet) ResultSet {
	newRS := &resultSet{
		rs: rs,
	}
	return newRS
}

// IsOK implements the ResultSet interface.
func (rs *resultSet) IsOK() bool {
	return rs.rs.IsOK()
}

// StatusCode implements the ResultSet interface.
func (rs *resultSet) StatusCode() int {
	return rs.rs.StatusCode()
}

// Error implements the ResultSet interface.
func (rs *resultSet) Error() error {
	return rs.rs.Error()
}

// LastSequence implements the ResultSet interface.
func (rs *resultSet) LastSequence() string {
	if err := rs.readChanges(); err != nil {
		return ""
	}
	return fmt.Sprintf("%v", rs.changes.LastSequence)
}

// Pending implements the ResultSet interface.
func (rs *resultSet) Pending() int {
	if err := rs.readChanges(); err != nil {
		return -1
	}
	return rs.changes.Pending
}

// Len implements the ResultSet interface.
func (rs *resultSet) Len() int {
	if err := rs.readChanges(); err != nil {
		return -1
	}
	return len(rs.changes.Results)
}

// Do implements the ResultSet interface.
func (rs *resultSet) Do(process Processor) error {
	if err := rs.readChanges(); err != nil {
		return err
	}
	for _, result := range rs.changes.Results {
		revisions := []string{}
		for _, change := range result.Changes {
			revisions = append(revisions, change.Revision)
		}
		seq := fmt.Sprintf("%v", result.Sequence)
		doc := couchdb.NewUnmarshableJSON(result.Document)
		if err := process(result.ID, seq, result.Deleted, revisions, doc); err != nil {
			return err
		}
	}
	return nil
}

// readChanges lazily reads the changes out of the CouchDB result set.
func (rs *resultSet) readChanges() error {
	if !rs.IsOK() {
		return rs.Error()
	}
	if rs.changes == nil {
		changes := couchdbChanges{}
		err := rs.rs.Document(&changes)
		if err != nil {
			return err
		}
		rs.changes = &changes
	}
	return nil
}

//--------------------
// HELPERS
//--------------------

// EOF
