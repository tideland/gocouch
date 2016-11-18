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
// RESPONSE
//--------------------

// Response contains the server response.
type Response interface {
	// IsOK checks the status code if the response is okay.
	IsOK() bool

	// Error return a possible error of a request.
	Error() error

	// ID returns a potentially returned document identifier.
	ID() (string, error)

	// Revision returns a potentially returned document revision.
	Revision() (string, error)

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
	cdbResp  *couchdbResponse
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

// Error implements the Response interface.
func (resp *response) Error() error {
	if resp.IsOK() {
		return nil
	}
	if resp.err != nil {
		return resp.err
	}
	if err := resp.readCouchDBResponse(); err != nil {
		return err
	}
	return errors.New(ErrClientRequest, errorMessages, resp.httpResp.StatusCode, resp.cdbResp.Error, resp.cdbResp.Reason)
}

// ID implements the Response interface.
func (resp *response) ID() (string, error) {
	if !resp.IsOK() {
		return "", resp.Error()
	}
	if err := resp.readCouchDBResponse(); err != nil {
		return "", err
	}
	return resp.cdbResp.ID, nil
}

// Revision implements the Response interface.
func (resp *response) Revision() (string, error) {
	if !resp.IsOK() {
		return "", resp.Error()
	}
	if err := resp.readCouchDBResponse(); err != nil {
		return "", err
	}
	return resp.cdbResp.Revision, nil
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

// readCouchDBResponse lazily loads the internal response of
// the CouchDB.
func (resp *response) readCouchDBResponse() error {
	if resp.cdbResp == nil {
		if err := resp.ResultValue(&resp.cdbResp); err != nil {
			return err
		}
	}
	return nil
}

// EOF
