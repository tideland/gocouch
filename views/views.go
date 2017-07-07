// Tideland Go CouchDB Client - Views
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package views

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// API
//--------------------

// View performs a view request.
func View(cdb couchdb.CouchDB, design, view string, params ...couchdb.Parameter) ViewResultSet {
	rs := cdb.GetOrPost(cdb.DatabasePath("_design", design, "_view", view), nil, params...)
	return newViewResultSet(rs)
}

//--------------------
// VIEW RESULT SET
//--------------------

// RowProcessingFunc is a function processing the content
// of a view row.
type RowProcessingFunc func(id string, key, value, document couchdb.Unmarshable) error

// ViewResultSet contains the result set of a view.
type ViewResultSet interface {
	// IsOK checks the status code if the result is okay.
	IsOK() bool

	// StatusCode returns the status code of the request.
	StatusCode() int

	// Error returns a possible error of a request.
	Error() error

	// TotalRows returns the number of ViewResultSet rows.
	TotalRows() int

	// ReturnedRows returns the nnumber of returned ViewResultSet rows.
	ReturnedRows() int

	// Offset returns the starting offset of the ViewResultSet rows.
	Offset() int

	// RowsDo iterates over the rows of a ViewResultSet and
	// processes the content.
	RowsDo(rpf RowProcessingFunc) error
}

// viewResultSet implements the ViewResultSet interface.
type viewResultSet struct {
	rs couchdb.ResultSet
	vr *couchdbViewResult
}

// newViewResultSet returns a ChangesResultSet.
func newViewResultSet(rs couchdb.ResultSet) ViewResultSet {
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

// ReturnedRows implements the ViewResultSet interface.
func (vrs *viewResultSet) ReturnedRows() int {
	if err := vrs.readViewResult(); err != nil {
		return -1
	}
	return len(vrs.vr.Rows)
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
		key := couchdb.NewUnmarshableJSON(row.Key)
		value := couchdb.NewUnmarshableJSON(row.Value)
		doc := couchdb.NewUnmarshableJSON(row.Document)
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
