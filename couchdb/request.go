// Tideland Go CouchDB Client - CouchDB - Request
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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	methHead   = "HEAD"
	methGet    = "GET"
	methPut    = "PUT"
	methPost   = "POST"
	methDelete = "DELETE"
)

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
	keys      []interface{}
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
		keys:   []interface{}{},
	}
	req.apply(cdb.parameters...)
	return req
}

// SetQuery implements the Parametrizable interface.
func (req *request) SetQuery(key, value string) {
	req.query.Set(key, value)
}

// AddQuery implements the Parametrizable interface.
func (req *request) AddQuery(key, value string) {
	req.query.Add(key, value)
}

// SetHeader implements the Parametrizable interface.
func (req *request) SetHeader(key, value string) {
	req.header.Set(key, value)
}

// AddKeys implements the Parametrizable interface.
func (req *request) AddKeys(keys ...interface{}) {
	req.keys = append(req.keys, keys...)
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
	return req.do(methHead)
}

// get performs a GET request.
func (req *request) get() *resultSet {
	return req.do(methGet)
}

// put performs a PUT request.
func (req *request) put() *resultSet {
	return req.do(methPut)
}

// post performs a POST request.
func (req *request) post() *resultSet {
	return req.do(methPost)
}

// delete performs a DELETE request.
func (req *request) delete() *resultSet {
	return req.do(methDelete)
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
	// Check if keys shall be used for the body.
	if len(req.keys) > 0 {
		req.doc = &couchdbViewKeys{Keys: req.keys}
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
	// Perform HTTP request.
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return newResultSet(nil, errors.Annotate(err, ErrPerformingRequest, errorMessages))
	}
	return newResultSet(httpResp, nil)
}

// EOF
