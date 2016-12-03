// Tideland Go CouchDB Client - CouchDB - RessultSet
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tideland/golib/errors"
)

//--------------------
// STATUS CODES
//--------------------

const (
	StatusOK       = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted

	StatusFound = http.StatusFound

	StatusBadRequest         = http.StatusBadRequest
	StatusUnauthorized       = http.StatusUnauthorized
	StatusForbidden          = http.StatusForbidden
	StatusNotFound           = http.StatusNotFound
	StatusPreconditionFailed = http.StatusPreconditionFailed

	StatusInternalServerError = http.StatusInternalServerError
)

//--------------------
// RESULT SET
//--------------------

// ResultSet contains the server result set.
type ResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// ID returns a potentially returned document identifier.
	ID() string

	// Revision returns a potentially returned document revision.
	Revision() string

	// Document returns the received document of a client
	// request and unmorshals it.
	Document(value interface{}) error

	// Raw returns the received raw data of a client request.
	Raw() ([]byte, error)
}

// resultSet implements the ResultSet interface.
type resultSet struct {
	resp   *http.Response
	status *Status
	err    error
}

// newResultSet analyzes the HTTP response and creates a the
// client ResultSet type out of it.
func newResultSet(resp *http.Response, err error) *resultSet {
	rs := &resultSet{
		resp: resp,
		err:  err,
	}
	return rs
}

// IsOK implements the resultSet interface.
func (rs *resultSet) IsOK() bool {
	return rs.err == nil && (rs.resp.StatusCode >= 200 && rs.resp.StatusCode <= 299)
}

// StatusCode implements the resultSet interface.
func (rs *resultSet) StatusCode() int {
	if rs.resp == nil {
		return -1
	}
	return rs.resp.StatusCode
}

// Error implements the resultSet interface.
func (rs *resultSet) Error() error {
	if rs.IsOK() {
		return nil
	}
	if rs.err != nil {
		return rs.err
	}
	if err := rs.readStatus(); err != nil {
		return err
	}
	return errors.New(ErrClientRequest, errorMessages, rs.resp.StatusCode, rs.status.Error, rs.status.Reason)
}

// ID implements the resultSet interface.
func (rs *resultSet) ID() string {
	if !rs.IsOK() {
		return ""
	}
	if err := rs.readStatus(); err != nil {
		return ""
	}
	return rs.status.ID
}

// Revision implements the resultSet interface.
func (rs *resultSet) Revision() string {
	if !rs.IsOK() {
		return ""
	}
	if err := rs.readStatus(); err != nil {
		return ""
	}
	return rs.status.Revision
}

// Document implements the resultSet interface.
func (rs *resultSet) Document(value interface{}) error {
	data, err := rs.Raw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return nil
}

// Raw implements the resultSet interface.
func (rs *resultSet) Raw() ([]byte, error) {
	if rs.err != nil {
		return nil, rs.err
	}
	defer rs.resp.Body.Close()
	body, err := ioutil.ReadAll(rs.resp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrReadingResponseBody, errorMessages)
	}
	return body, nil
}

// readStatus lazily loads the internal status resultSet
// of CouchDB.
func (rs *resultSet) readStatus() error {
	if rs.status == nil {
		if err := rs.Document(&rs.status); err != nil {
			return err
		}
	}
	return nil
}

// header returns the value of a response header.
func (rs *resultSet) header(key string) string {
	return rs.resp.Header.Get(key)
}

// EOF
