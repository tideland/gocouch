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

import ()

//--------------------
// EXTERNAL DOCUMENT TYPES
//--------------------

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
