// Tideland Go CouchDB Client - CouchDB - Marshables
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
	"encoding/json"

	"github.com/tideland/golib/errors"
)

//--------------------
// UNMARSHABLE
//--------------------

// Unmarshable describes a not yet unmarshalled value that
// can be unmarshalled into a given variable. It is used to
// access key, value, or document of view result rows.
type Unmarshable interface {
	// String returns the original as string.
	String() string

	// Raw returns the unmarshable as raw byte slice.
	Raw() []byte

	// Unmarshal unmarshals the interface into the
	// passed variable.
	Unmarshal(doc interface{}) error
}

// unmarshable implements the Unmarshable interface.
type unmarshable struct {
	message json.RawMessage
}

// NewUnmarshableRaw creates a new Unmarshable out of
// the raw bytes.
func NewUnmarshableRaw(raw []byte) Unmarshable {
	return NewUnmarshableJSON(json.RawMessage(raw))
}

// NewUnmarshableJSON creates a new Unmarshable out of
// a json.RawMessage
func NewUnmarshableJSON(msg json.RawMessage) Unmarshable {
	return &unmarshable{
		message: msg,
	}
}

// String implements the Unmarshable interface.
func (u *unmarshable) String() string {
	if u.message == nil {
		return ""
	}
	return string(u.message)
}

// Raw implements the Unmarshable interface.
func (u *unmarshable) Raw() []byte {
	if u.message == nil {
		return nil
	}
	dest := make([]byte, len(u.message))
	copy(dest, u.message)
	return dest
}

// Unmarshal implements the Unmarshable interface.
func (u *unmarshable) Unmarshal(doc interface{}) error {
	err := json.Unmarshal(u.message, doc)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return nil
}

// EOF
