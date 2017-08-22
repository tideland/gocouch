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
	"encoding/json"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// API
//--------------------

// Request contains the parameters of a find request.
// TODO(mue) Think about selector as fixed argument and
// optional parameters.
type Request struct {
	Selector Selector `json:"selector"`
	Limit    int      `json:"limit,omitempty"`
	Skip     int      `json:"skip,omitempty"`
	Sort     string   `json:"sort,omitempty"`
	Fields   []string `json:"fields,omitempty"`
}

// Find returns access to the found results.
func Find(cdb couchdb.CouchDB, request Request) ResultSet {
	rs := cdb.Post(cdb.DatabasePath("_find"), request)
	return newResultSet(rs)
}

//--------------------
// FIND RESULT SET
//--------------------

// Processor is a function processing the content of a found document.
type Processor func(document couchdb.Unmarshable) error

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

// resultSet implements the ResultSet interface.
type resultSet struct {
	rs       couchdb.ResultSet
	response *response
}

// newResultSet returns a ResultSet.
func newResultSet(rs couchdb.ResultSet) ResultSet {
	frs := &resultSet{
		rs: rs,
	}
	return frs
}

// IsOK implements the ResultSet interface.
func (frs *resultSet) IsOK() bool {
	return frs.rs.IsOK()
}

// StatusCode implements the ResultSet interface.
func (frs *resultSet) StatusCode() int {
	return frs.rs.StatusCode()
}

// Error implements the ResultSet interface.
func (frs *resultSet) Error() error {
	return frs.rs.Error()
}

// Len implements ResultSet.
func (frs *resultSet) Len() int {
	if err := frs.readResponse(); err != nil {
		return -1
	}
	return len(frs.response.Documents)
}

// Do implements ResultSet.
func (frs *resultSet) Do(process Processor) error {
	if err := frs.readResponse(); err != nil {
		return err
	}
	for _, doc := range frs.response.Documents {
		unmarshableDoc := couchdb.NewUnmarshableJSON(doc)
		if err := process(unmarshableDoc); err != nil {
			return err
		}
	}
	return nil
}

// readResponse lazily reads and analyzes response.
func (frs *resultSet) readResponse() error {
	if !frs.IsOK() {
		return frs.Error()
	}
	if frs.response == nil {
		resp := response{}
		err := frs.rs.Document(&resp)
		if err != nil {
			return err
		}
		frs.response = &resp
	}
	return nil
}

type response struct {
	Warning   string            `json:"warning"`
	Documents []json.RawMessage `json:"docs"`
}

// EOF
