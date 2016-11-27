// Tideland Go CouchDB Client - CouchDB - Response
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
// RESPONSE
//--------------------

// Response contains the server response.
type Response interface {
	// IsOK checks the status code if the response is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// ID returns a potentially returned document identifier.
	ID() string

	// Revision returns a potentially returned document revision.
	Revision() string

	// ResultValue returns the received document of a client
	// request and unmorshals it.
	ResultValue(value interface{}) error

	// ResultData returns the received data of a client
	// request.
	ResultData() ([]byte, error)
}

// response implements the Response interface.
type response struct {
	httpResp *http.Response
	result   *Result
	err      error
}

// newResponse analyzes the HTTP response and creates a the
// client response type out of it.
func newResponse(httpResp *http.Response, err error) *response {
	resp := &response{
		httpResp: httpResp,
		err:      err,
	}
	return resp
}

// IsOK implements the Response interface.
func (resp *response) IsOK() bool {
	return resp.err == nil && (resp.httpResp.StatusCode >= 200 && resp.httpResp.StatusCode <= 299)
}

// StatusCode implements the Response interface.
func (resp *response) StatusCode() int {
	if resp.httpResp == nil {
		return -1
	}
	return resp.httpResp.StatusCode
}

// Error implements the Response interface.
func (resp *response) Error() error {
	if resp.IsOK() {
		return nil
	}
	if resp.err != nil {
		return resp.err
	}
	if err := resp.readResult(); err != nil {
		return err
	}
	return errors.New(ErrClientRequest, errorMessages, resp.httpResp.StatusCode, resp.result.Error, resp.result.Reason)
}

// ID implements the Response interface.
func (resp *response) ID() string {
	if !resp.IsOK() {
		return ""
	}
	if err := resp.readResult(); err != nil {
		return ""
	}
	return resp.result.ID
}

// Revision implements the Response interface.
func (resp *response) Revision() string {
	if !resp.IsOK() {
		return ""
	}
	if err := resp.readResult(); err != nil {
		return ""
	}
	return resp.result.Revision
}

// ResultValue implements the Response interface.
func (resp *response) ResultValue(value interface{}) error {
	data, err := resp.ResultData()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return nil
}

// ResultData implements the Response interface.
func (resp *response) ResultData() ([]byte, error) {
	if resp.err != nil {
		return nil, resp.err
	}
	defer resp.httpResp.Body.Close()
	body, err := ioutil.ReadAll(resp.httpResp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrReadingResponseBody, errorMessages)
	}
	return body, nil
}

// readResult lazily loads the internal response of
// the CouchDB.
func (resp *response) readResult() error {
	if resp.result == nil {
		if err := resp.ResultValue(&resp.result); err != nil {
			return err
		}
	}
	return nil
}

// EOF
