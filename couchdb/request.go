// Tideland Go CouchDB Client - CouchDB - Request
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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/logger"
)

//--------------------
// CONSTANTS
//--------------------

//--------------------
// REQUEST
//--------------------

// request is responsible for an individual request to a CouchDB.
type request struct {
	cdb       *couchdb
	path      string
	doc       interface{}
	docReader io.Reader
	query     url.Values
	header    http.Header
}

// newRequest creates a new request for the given location, method, and path. If needed
// query and header can be added like newRequest().setQuery().setHeader.do().
func newRequest(cdb *couchdb, path string, doc interface{}) *request {
	req := &request{
		cdb:    cdb,
		path:   path,
		doc:    doc,
		query:  url.Values{},
		header: http.Header{},
	}
	req.apply(cdb.parameters...)
	return req
}

// SetQuery implements the Parameterizable interface.
func (req *request) SetQuery(key, value string) {
	req.query.Set(key, value)
}

// AddQuery implements the Parameterizable interface.
func (req *request) AddQuery(key, value string) {
	req.query.Add(key, value)
}

// SetHeader implements the Parameterizable interface.
func (req *request) SetHeader(key, value string) {
	req.header.Set(key, value)
}

// UpdateDocument implements the Parameterizable interface.
func (req *request) UpdateDocument(update func(interface{}) interface{}) {
	req.doc = update(req.doc)
}

// apply applies a list of parameters to the request.
func (req *request) apply(params ...Parameter) *request {
	for _, param := range params {
		param(req)
	}
	return req
}

// head performs a HEAD request.
func (req *request) head() *resultSet {
	return req.do(http.MethodHead)
}

// get performs a GET request.
func (req *request) get() *resultSet {
	return req.do(http.MethodGet)
}

// put performs a PUT request.
func (req *request) put() *resultSet {
	return req.do(http.MethodPut)
}

// post performs a POST request.
func (req *request) post() *resultSet {
	return req.do(http.MethodPost)
}

// delete performs a DELETE request.
func (req *request) delete() *resultSet {
	return req.do(http.MethodDelete)
}

// do performs a request.
func (req *request) do(method string) *resultSet {
	// Prepare URL.
	u := &url.URL{
		Scheme: "http",
		Host:   req.cdb.host,
		Path:   req.path,
	}
	if len(req.query) > 0 {
		u.RawQuery = req.query.Encode()
	}
	// Marshal a potential document.
	if req.doc != nil {
		marshalled, err := json.Marshal(req.doc)
		if err != nil {
			return newResultSet(nil, errors.Annotate(err, ErrMarshallingDoc, errorMessages))
		}
		req.docReader = bytes.NewBuffer(marshalled)
	}
	// Prepare HTTP request.
	httpReq, err := http.NewRequest(method, u.String(), req.docReader)
	if err != nil {
		return newResultSet(nil, errors.Annotate(err, ErrPreparingRequest, errorMessages))
	}
	httpReq.Close = true
	if len(req.header) > 0 {
		httpReq.Header = req.header
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	// Log if wanted.
	if req.cdb.debugLog {
		logger.Debugf("couchdb request '%s %s'", method, u)
	}
	// Perform HTTP request.
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return newResultSet(nil, errors.Annotate(err, ErrPerformingRequest, errorMessages))
	}
	return newResultSet(httpResp, nil)
}

// EOF
