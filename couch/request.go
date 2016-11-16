// Tideland Go CouchDB Client - Couch - Request
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couch

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

//--------------------
// CONSTANTS
//--------------------

const (
	methGet  = "GET"
	methPost = "POST"
)

//--------------------
// REQUEST
//--------------------

// request is responsible for an individual request to a CouchDB.
type request struct {
	url       *url.URL
	method    string
	doc       interface{}
	docReader io.Reader
	header    *http.Header
}

// newRequest creates a new request for the given location, method, and path. If needed
// query and header can be added like newRequest().setQuery().setHeader.do().
func newRequest(url *url.URL, method, path string, doc interface{}) *request {
	req := &request{
		url:    url,
		method: method,
		doc:    doc,
	}
	req.url.Path = path
	return req
}

// setQuery sets query values.
func (req *request) setQuery(values *url.Values) *request {
	req.url.RawQuery = values.Encode()
	return req
}

// setHeader sets header values.
func (req *request) setHeader(header *http.Header) *request {
	req.header = header
	return req
}

// do performs a request.
func (req *request) do() (*response, error) {
	// Marshal a potential document.
	if req.doc != nil {
		marshalled, err := json.Marshal(req.doc)
		if err != nil {
			return nil, errors.Annotate(err, ErrMarshallingDoc, errorMessages)
		}
		req.docReader = bytes.NewBuffer(marshalled)
	}
	// Prepare HTTP request.
	httpReq, err := http.NewRequest(req.method, req.url.String(), req.docReader)
	if err != nil {
		return nil, errors.Annotate(err, ErrPreparingRequest, errorMessages)
	}
	if req.header != nil {
		httpReq.Header = *req.header
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	// Perform HTTP request.
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, errors.Annotate(err, ErrPerformingRequest, errorMessages)
	}
	return newResponse(httpResp)
}

// EOF
