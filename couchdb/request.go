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
	db        *couchdb
	path      string
	doc       interface{}
	docReader io.Reader
	query     *Query
	header    *http.Header
}

// newRequest creates a new request for the given location, method, and path. If needed
// query and header can be added like newRequest().setQuery().setHeader.do().
func newRequest(db *couchdb, path string, doc interface{}) *request {
	req := &request{
		db:   db,
		path: path,
		doc:  doc,
	}
	return req
}

// setQuery sets query values.
func (req *request) setQuery(query *Query) *request {
	req.query = query
	return req
}

// setHeader sets header values.
func (req *request) setHeader(header *http.Header) *request {
	req.header = header
	return req
}

// get performs a GET request.
func (req *request) get() *response {
	return req.do(methGet)
}

// put performs a PUT request.
func (req *request) put() *response {
	return req.do(methPut)
}

// post performs a POST request.
func (req *request) post() *response {
	return req.do(methPost)
}

// delete performs a DELETE request.
func (req *request) delete() *response {
	return req.do(methDelete)
}

// do performs a request.
func (req *request) do(method string) *response {
	// Prepare URL.
	u := &url.URL{
		Scheme: "http",
		Host:   req.db.host,
		Path:   req.path,
	}
	if req.query != nil {
		u.RawQuery = req.query.Encode()
	}
	// Marshal a potential document.
	if req.doc != nil {
		marshalled, err := json.Marshal(req.doc)
		if err != nil {
			return newResponse(nil, errors.Annotate(err, ErrMarshallingDoc, errorMessages))
		}
		req.docReader = bytes.NewBuffer(marshalled)
	}
	// Prepare HTTP request.
	httpReq, err := http.NewRequest(method, u.String(), req.docReader)
	if err != nil {
		return newResponse(nil, errors.Annotate(err, ErrPreparingRequest, errorMessages))
	}
	if req.header != nil {
		httpReq.Header = *req.header
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	// Perform HTTP request.
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return newResponse(nil, errors.Annotate(err, ErrPerformingRequest, errorMessages))
	}
	return newResponse(httpResp, nil)
}

// EOF
