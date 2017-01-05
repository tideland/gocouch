// Tideland Go CouchDB Client - CouchDB - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes.
const (
	ErrNoConfiguration = iota + 1
	ErrInvalidDocument
	ErrNoIdentifier
	ErrNotFound
	ErrMarshallingDoc
	ErrPreparingRequest
	ErrPerformingRequest
	ErrClientRequest
	ErrUnmarshallingDoc
	ErrUnmarshallingField
	ErrReadingResponseBody
)

// Error messages.
var errorMessages = errors.Messages{
	ErrNoConfiguration:     "cannot open database without configuration",
	ErrInvalidDocument:     "document needs _id and _rev",
	ErrNoIdentifier:        "document contains no identifier",
	ErrNotFound:            "document with identifier '%s' not found",
	ErrMarshallingDoc:      "cannot marshal into database document",
	ErrPreparingRequest:    "cannot prepare request",
	ErrPerformingRequest:   "cannot perform request",
	ErrClientRequest:       "client request failed: status code %d, error '%s', reason '%s'",
	ErrUnmarshallingDoc:    "cannot unmarshal database document",
	ErrUnmarshallingField:  "cannot unmarshal the document field",
	ErrReadingResponseBody: "cannot read response body",
}

// EOF
