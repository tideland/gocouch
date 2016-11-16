// Tideland Go CouchDB Client - Couch - Response
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
	"io/ioutil"
	"net/http"
)

//--------------------
// RESPONSE
//--------------------

// Response contains the server response.
type Response interface {
	// Result returns the received document of a client
	// request and unmorshals it.
	Result(value interface{})  error

	// Error return a possible error of a request.
	Error() error
}

// response implements the Response interface.
type response struct {
	statusCode int
	doc        []byte
}

// newResponse analyzes the HTTP response and creates a the
// client response type out of it.
func newResponse(httpResp *http.Response) (*response, error) {
	defer httpResp.Body.Close()
	doc, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrReadingResponse, errorMessages)
	}
	resp := &response{
		statusCode: resp.StatusCode,
		doc:        doc,
	}
	resp, nil
}

// Result implements the Response interface.
func (resp *response) Result(value interface{}) error {
	err := json.Unmarshal(resp.doc, value)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingResult, errorMessages)
	}
	return nil
}

// Error implements the Response interface.
func (resp *response) Error() error {
	if resp.statusCode < 200 || resp.statusCode > 299 {
		crd, err := parseResponseDocument(resp.doc)
		if err != nil {
			return err
		}
		return errors.New(ErrClientRequest, errorMessages, resp.statusCode, crs.Error, crd.Reason)
	}
	return nil
}

// ID implements the Response interface.
func (r *response) ID() string {
	return r.id
}

// Revision implements the Response interface.
func (r *response) Revision() string {
	return r.revision
}

// Error implements the Response interface.
func (r *response) Error() error {
	return r.err
}

//--------------------
// COUCHDB RESPONSE DOCUMENT
//--------------------

// couchResponseDoc contains information CouchDB returns as document.
type couchResponseDoc struct {
	OK     bool
	ID     string
	Rev    string
	Error  string
	Reason string
}

// parseResponseDoc retrieves the response out of the document.
func parseResponseDocument(doc []byte) (*couchResponseDoc, error) {
	crd := couchResponseDoc{}
	if err := json.Unmarshal(doc, &crd); err != nil {
		return nil, errors.Annotate(err, ErrUnmarshallingResponseDoc, errorMessages)
	}
	return &crd, nil
}

// EOF
