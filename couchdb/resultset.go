// Tideland GoCouch - CouchDB - RessultSet
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
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

// Status codes after database requests.
const (
	StatusOK       = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted

	StatusFound = http.StatusFound

	StatusBadRequest         = http.StatusBadRequest
	StatusUnauthorized       = http.StatusUnauthorized
	StatusForbidden          = http.StatusForbidden
	StatusNotFound           = http.StatusNotFound
	StatusMethodNotAllowed   = http.StatusMethodNotAllowed
	StatusNotAcceptable      = http.StatusNotAcceptable
	StatusConflict           = http.StatusConflict
	StatusPreconditionFailed = http.StatusPreconditionFailed
	StatusTooManyRequests    = http.StatusTooManyRequests

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

	// IsDeleted returns true if a returned document is already deleted.
	IsDeleted() bool

	// Document returns the received document of a client
	// request and unmorshals it.
	Document(value interface{}) error

	// Raw returns the received raw data of a client request.
	Raw() ([]byte, error)

	// Header provides access to header variables.
	Header(key string) string
}

// resultSet implements the ResultSet interface.
type resultSet struct {
	statusCode  int
	body        []byte
	headers     map[string]string
	document    map[string]interface{}
	id          string
	revision    string
	deleted     bool
	errorText   string
	errorReason string
	err         error
}

// newResultSet analyzes the HTTP response and creates a the
// client ResultSet type out of it.
func newResultSet(resp *http.Response, err error) *resultSet {
	rs := &resultSet{
		statusCode: 200,
		err:        err,
	}
	if err != nil && errors.IsError(err, ErrNotFound) {
		rs.statusCode = StatusNotFound
	}
	if resp != nil {
		// Get status code.
		rs.statusCode = resp.StatusCode
		// Read body.
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			rs.err = errors.Annotate(err, ErrReadingResponseBody, errorMessages)
		}
		rs.body = body
		// Read headers.
		rs.headers = make(map[string]string)
		for key, values := range resp.Header {
			if len(values) > 0 {
				rs.headers[key] = values[0]
			}
		}
	}
	return rs
}

// IsOK implements the ResultSet interface.
func (rs *resultSet) IsOK() bool {
	return rs.err == nil && (rs.statusCode >= 200 && rs.statusCode <= 299)
}

// StatusCode implements the ResultSet interface.
func (rs *resultSet) StatusCode() int {
	return rs.statusCode
}

// Error implements the ResultSet interface.
func (rs *resultSet) Error() error {
	if rs.IsOK() {
		return nil
	}
	if rs.err != nil {
		return rs.err
	}
	if err := rs.readDocument(); err != nil {
		return err
	}
	return errors.New(ErrClientRequest, errorMessages, rs.statusCode, rs.errorText, rs.errorReason)
}

// ID implements the ResultSet interface.
func (rs *resultSet) ID() string {
	if !rs.IsOK() {
		return ""
	}
	if err := rs.readDocument(); err != nil {
		return ""
	}
	return rs.id
}

// Revision implements the ResultSet interface.
func (rs *resultSet) Revision() string {
	if !rs.IsOK() {
		return ""
	}
	if err := rs.readDocument(); err != nil {
		return ""
	}
	return rs.revision
}

// IsDeleted implements the ResultSet interface.
func (rs *resultSet) IsDeleted() bool {
	if !rs.IsOK() {
		return false
	}
	if err := rs.readDocument(); err != nil {
		return false
	}
	return rs.deleted
}

// Document implements the ResultSet interface.
func (rs *resultSet) Document(value interface{}) error {
	if rs.err != nil {
		return rs.err
	}
	err := json.Unmarshal(rs.body, value)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return nil
}

// Raw implements the ResultSet interface.
func (rs *resultSet) Raw() ([]byte, error) {
	return rs.body, rs.err
}

// Header implements the ResultSet interface.
func (rs *resultSet) Header(key string) string {
	value, ok := rs.headers[key]
	if !ok {
		return ""
	}
	return value
}

// readDocument lazily loads and analyzis a generic document.
func (rs *resultSet) readDocument() error {
	if rs.document == nil {
		rs.document = make(map[string]interface{})
		if err := rs.Document(&rs.document); err != nil {
			return err
		}
		if id, ok := rs.document["_id"]; ok {
			rs.id = id.(string)
		} else if id, ok := rs.document["id"]; ok {
			rs.id = id.(string)
		}
		if revision, ok := rs.document["_rev"]; ok {
			rs.revision = revision.(string)
		} else if revision, ok := rs.document["rev"]; ok {
			rs.revision = revision.(string)
		}
		if deleted, ok := rs.document["_deleted"]; ok {
			rs.deleted = deleted.(bool)
		}
		if errorText, ok := rs.document["error"]; ok {
			rs.errorText = errorText.(string)
		}
		if errorReason, ok := rs.document["reason"]; ok {
			rs.errorReason = errorReason.(string)
		}
	}
	return nil
}

// EOF
