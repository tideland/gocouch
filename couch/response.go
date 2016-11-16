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

//--------------------
// RESPONSE
//--------------------

// Response contains the server response.
type Response interface {
	// IsOK checks the status code if the response is okay.
	IsOK() bool

	// Error return a possible error of a request.
	Error() error

	// ResultValue returns the received document of a client
	// request and unmorshals it.
	ResultValue(value interface{})  error

	// ResultData returns the received data of a client
	// request.
	ResultData() ([]byte, error)
}

// response implements the Response interface.
type response struct {
	statusCode int
	body        []byte
	crd        *couchResponseDoc
}

// newResponse analyzes the HTTP response and creates a the
// client response type out of it.
func newResponse(httpResp *http.Response) (*response, error) {
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrReadingResponse, errorMessages)
	}
	resp := &response{
		statusCode: resp.StatusCode,
		body:        body,
	}
	return resp, nil
}

// IsOK implements the Response interface.
func (resp *response) IsOK() bool {
	return resp.statusCode >= 200 && resp.statusCode <= 299
}

// Error implements the Response interface.
func (resp *response) Error() error {
	if !resp.IsOK() {
		crd, err := parseResponseDocument(resp.doc)
		if err != nil {
			return err
		}
		return errors.New(ErrClientRequest, errorMessages, resp.statusCode, crs.Error, crd.Reason)
	}
	return nil
}

// ResultValue implements the Response interface.
func (resp *response) ResultValue(value interface{}) error {
	if !resp.IsOK() {
		return resp.Error()
	}
	err := json.Unmarshal(resp.body, value)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingResult, errorMessages)
	}
	return nil
}

// ResultData implements the Response interface.
func (resp *response) ResultData() ([]byte, error) {
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	return resp.body, nil
}

// EOF
