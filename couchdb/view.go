// Tideland Go CouchDB Client - CouchDB - View
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

	"github.com/tideland/golib/errors"
)

//--------------------
// UNMARSHABLE
//--------------------

// Unmarshable describes an empty interface that can be
// unmarshalled into a given variable. It is used to
// access key, value, or document of view result rows.
type Unmarshable interface {
	// Raw returns the original empty interface.
	Raw() interface{}

	// JSON returns the JSON encoded value as it has
	// been part of a row.
	JSON() ([]byte, error)

	// Unmarshal unmarshals the interface into the
	// passed variable.
	Unmarshal(doc interface{}) error
}

// unmarshable implements the Unmarshable interface.
type unmarshable struct {
	value interface{}
}

// Raw implements the Unmarshable interface.
func (u unmarshable) Raw() interface{} {
	return u.value
}

// JSON implements the Unmarshable interface.
func (u unmarshable) JSON() ([]byte, error) {
	jb, err := json.Marshal(u.value)
	if err != nil {
		return nil, errors.Annotate(err, ErrRemarshalling, errorMessages)
	}
	return jb, nil
}

// Unmarshal implements the Unmarshable interface.
func (u unmarshable) Unmarshal(doc interface{}) error {
	jb, err := u.JSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(jb, doc)
	if err != nil {
		return errors.Annotate(err, ErrRemarshalling, errorMessages)
	}
	return nil
}

//--------------------
// VIEW RESULT SET
//--------------------

// RowProcessingFunc is a function processing the content
// of a viewResultSet row.
type RowProcessingFunc func(id string, key, value, document Unmarshable) error

// ViewResultSet contains the viewResultSet result set.
type ViewResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// TotalRows returns the number of viewResultSet rows.
	TotalRows() int

	// Offset returns the starting offset of the viewResultSet rows.
	Offset() int

	// RowsDo iterates over the rows of a viewResultSet and
	// processes the content.
	RowsDo(rpf RowProcessingFunc) error
}

// viewResultSet implements the ViewResultSet interface.
type viewResultSet struct {
	rs ResultSet
	vr *couchdbViewResult
}

// newView provides access to the viewResultSet data.
func newView(rs ResultSet) *viewResultSet {
	vrs := &viewResultSet{
		rs: rs,
	}
	return vrs
}

// IsOK implements the ViewResultSet interface.
func (vrs *viewResultSet) IsOK() bool {
	return vrs.rs.IsOK()
}

// StatusCode implements the ViewResultSet interface.
func (vrs *viewResultSet) StatusCode() int {
	return vrs.rs.StatusCode()
}

// Error implements the ViewResultSet interface.
func (vrs *viewResultSet) Error() error {
	return vrs.rs.Error()
}

// TotalRows implements the ViewResultSet interface.
func (vrs *viewResultSet) TotalRows() int {
	if err := vrs.readViewResult(); err != nil {
		return -1
	}
	return vrs.vr.TotalRows
}

// Offset implements the ViewResultSet interface.
func (vrs *viewResultSet) Offset() int {
	if err := vrs.readViewResult(); err != nil {
		return -1
	}
	return vrs.vr.Offset
}

// RowsDo implements the View interface.
func (vrs *viewResultSet) RowsDo(rpf RowProcessingFunc) error {
	if err := vrs.readViewResult(); err != nil {
		return err
	}
	for _, row := range vrs.vr.Rows {
		key := unmarshable{row.Key}
		value := unmarshable{row.Value}
		doc := unmarshable{row.Document}
		if err := rpf(row.ID, key, value, doc); err != nil {
			return err
		}
	}
	return nil
}

// readViewResult lazily reads the viewResultSet result.
func (vrs *viewResultSet) readViewResult() error {
	if !vrs.IsOK() {
		return vrs.Error()
	}
	if vrs.vr == nil {
		vr := couchdbViewResult{}
		err := vrs.rs.Document(&vr)
		if err != nil {
			return err
		}
		vrs.vr = &vr
	}
	return nil
}

// EOF
