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
func Changes(cdb couchdb.CouchDB, params ...couchdb.Parameter) ChangesResultSet {
	rs := cdb.GetOrPost(cdb.DatabasePath("_changes"), nil, params...)
	return newChangesResultSet(rs)
}

//--------------------
// CHANGES RESULT SET
//--------------------

// ChangesProcessingFunc is a function processing the content
// of a changes row.
type ChangesProcessingFunc func(id, sequence string, deleted bool, revisions []string, document couchdb.Unmarshable) error

// ChangesResultSet contains the result set of a change.
type ChangesResultSet interface {
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

	// ResultsLen returns the number of changes.
	ResultsLen() int

	// ResultsDo iterates over the results of a ChangesResultSet and
	// processes the content.
	ResultsDo(cpf ChangesProcessingFunc) error
}

// changesResultSet implements the ChangesResultSet interface.
type changesResultSet struct {
	rs      couchdb.ResultSet
	changes *couchdbChanges
}

// newChangesResultSet returns a ChangesResultSet.
func newChangesResultSet(rs couchdb.ResultSet) ChangesResultSet {
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

// LastSequence implements the ChangesResultSet interface.
func (crs *changesResultSet) LastSequence() string {
	if err := crs.readChangesResult(); err != nil {
		return ""
	}
	return fmt.Sprintf("%v", crs.changes.LastSequence)
}

// Pending implements the ChangesResultSet interface.
func (crs *changesResultSet) Pending() int {
	if err := crs.readChangesResult(); err != nil {
		return -1
	}
	return crs.changes.Pending
}

// ResultsLen implements the ChangesResultSet interface.
func (crs *changesResultSet) ResultsLen() int {
	if err := crs.readChangesResult(); err != nil {
		return -1
	}
	return len(crs.changes.Results)
}

// ResultsDo implements the ChangesResultSet interface.
func (crs *changesResultSet) ResultsDo(cpf ChangesProcessingFunc) error {
	if err := crs.readChangesResult(); err != nil {
		return err
	}
	for _, result := range crs.changes.Results {
		revisions := []string{}
		for _, change := range result.Changes {
			revisions = append(revisions, change.Revision)
		}
		seq := fmt.Sprintf("%v", result.Sequence)
		doc := couchdb.NewUnmarshableJSON(result.Document)
		if err := cpf(result.ID, seq, result.Deleted, revisions, doc); err != nil {
			return err
		}
	}
	return nil
}

// readChangesResult lazily reads the viewResultSet result.
func (crs *changesResultSet) readChangesResult() error {
	if !crs.IsOK() {
		return crs.Error()
	}
	if crs.changes == nil {
		changes := couchdbChanges{}
		err := crs.rs.Document(&changes)
		if err != nil {
			return err
		}
		crs.changes = &changes
	}
	return nil
}

//--------------------
// HELPERS
//--------------------

// EOF
