// Tideland Go CouchDB Client - CouchDB - Document Types
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
// EXTERNAL DOCUMENT TYPES
//--------------------

// ViewRow contains one row of a view result. The unmarshal
// methods allow to transform more complex results into according
// variables.
type ViewRow struct {
	ID       string      `json:"id"`
	Key      interface{} `json:"key"`
	Value    interface{} `json:"value"`
	Document interface{} `json:"doc"`
}

// UnmarshalKey converts the key field of the row into the
// passed variable.
func (vr ViewRow) UnmarshalKey(key interface{}) error {
	return vr.remarshal(vr.Key, key)
}

// UnmarshalValue converts the value field of the row into the
// passed variable.
func (vr ViewRow) UnmarshalValue(value interface{}) error {
	return vr.remarshal(vr.Value, value)
}

// UnmarshalDocument converts the document field of the row into the
// passed variable.
func (vr ViewRow) UnmarshalDocument(doc interface{}) error {
	return vr.remarshal(vr.Document, doc)
}

// remarshal marshals the in value to JSON again and unmarshals
// it to the out value.
func (vr ViewRow) remarshal(in, out interface{}) error {
	tmp, err := json.Marshal(in)
	if err != nil {
		return errors.Annotate(err, ErrRemarshalling, errorMessages)
	}
	err = json.Unmarshal(tmp, out)
	if err != nil {
		return errors.Annotate(err, ErrRemarshalling, errorMessages)
	}
	return nil
}

type ViewRows []ViewRow

// ViewResult is a generic result of a CouchDB view.
type ViewResult struct {
	TotalRows int      `json:"total_rows"`
	Offset    int      `json:"offset"`
	Rows      ViewRows `json:"rows"`
}

// Status contains internal status information CouchDB returns.
type Status struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
	Error    string `json:"error"`
	Reason   string `json:"reason"`
}

// Statuaess is the list of status information after a bulk writing.
type Statuses []Status

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbAuthentication contains user ID and password
// for authentication.
type couchdbAuthentication struct {
	UserID   string `json:"name"`
	Password string `json:"password"`
}

// couchdRoles contains the roles of a user if the
// authentication succeeded.
type couchdbRoles struct {
	OK     bool     `json:"ok"`
	UserID string   `json:"name"`
	Roles  []string `json:"roles"`
}

// couchdbBulkDocuments contains a number of documents added at once.
type couchdbBulkDocuments struct {
	Docs     []interface{} `json:"docs"`
	NewEdits bool          `json:"new_edits,omitempty"`
}

// couchdbViewKeys sets key constraints for view requests.
type couchdbViewKeys struct {
	Keys []interface{} `json:"keys"`
}

// idAndRevision is used to simply retrieve ID and revision of
// a document.
type idAndRevision struct {
	ID       string `json:"_id"`
	Revision string `json:"_rev"`
}

// EOF
