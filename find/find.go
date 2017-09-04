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

// Find returns access to the found results.
func Find(cdb couchdb.CouchDB, selector Selector, parameters ...Parameter) ResultSet {
	// Create request object.
	req := request{}
	req.SetParameter("selector", selector)
	req.apply(parameters...)
	// Perform find command.
	rs := cdb.Post(cdb.DatabasePath("_find"), req)
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
	rs          couchdb.ResultSet
	response    *response
	responseErr error
}

// newResultSet returns a ResultSet.
func newResultSet(rs couchdb.ResultSet) ResultSet {
	frs := &resultSet{
		rs: rs,
	}
	resp := response{}
	err := frs.rs.Document(&resp)
	if err != nil {
		frs.responseErr = err
	} else {
		frs.response = &resp
	}
	return frs
}

// IsOK implements the ResultSet interface.
func (frs *resultSet) IsOK() bool {
	return frs.rs.IsOK() && frs.responseErr == nil
}

// StatusCode implements the ResultSet interface.
func (frs *resultSet) StatusCode() int {
	return frs.rs.StatusCode()
}

// Error implements the ResultSet interface.
func (frs *resultSet) Error() error {
	if frs.rs.Error() != nil {
		return frs.rs.Error()
	}
	return frs.responseErr
}

// Len implements ResultSet.
func (frs *resultSet) Len() int {
	if !frs.IsOK() {
		return -1
	}
	return len(frs.response.Documents)
}

// Do implements ResultSet.
func (frs *resultSet) Do(process Processor) error {
	for _, doc := range frs.response.Documents {
		unmarshableDoc := couchdb.NewUnmarshableJSON(doc)
		if err := process(unmarshableDoc); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// REQUEST AND RESPONSE
//--------------------

// request contains all request object fields.
type request map[string]interface{}

// SetParameter implements Parameterizable.
func (req request) SetParameter(key string, parameter interface{}) {
	req[key] = parameter
}

// apply applies a list of parameters to the request.
func (req request) apply(parameters ...Parameter) {
	for _, applyParameterTo := range parameters {
		applyParameterTo(req)
	}
}

// response describes the document returned by CouchDB.
type response struct {
	Warning   string            `json:"warning"`
	Documents []json.RawMessage `json:"docs"`
}

// EOF
