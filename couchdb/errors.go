// Tideland Go CouchDB Client - CouchDB - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrNoConfiguration = iota + 1
	ErrMarshallingDoc
	ErrPreparingRequest
	ErrPerformingRequest
	ErrClientRequest
	ErrUnmarshallingDoc
	ErrReadingResponseBody
)

var errorMessages = errors.Messages{
	ErrNoConfiguration:     "cannot open database without configuration",
	ErrMarshallingDoc:      "cannot marshal into database document",
	ErrPreparingRequest:    "cannot prepare request",
	ErrPerformingRequest:   "cannot perform request",
	ErrClientRequest:       "client request failed; status code %d, error '%s', reason '%s'",
	ErrUnmarshallingDoc:    "cannot unmarshal database document",
	ErrReadingResponseBody: "cannot read response body",
}

// EOF
