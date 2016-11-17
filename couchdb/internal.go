// Tideland Go CouchDB Client - CouchDB - Internal Document Types
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
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbResponse contains response information CouchDB returns.
type couchdbResponse struct {
	OK     bool   `json:"ok"`
	ID     string `json:"_id"`
	Rev    string `json:"_rev"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// couchdbViewResult is a generic result of a CouchDB view.
type couchdbViewResult struct {
	TotalRows int `json:"total_rows"`
	Offset    int `json:"offset"`
	Rows      []struct {
		ID       string      `json:"id"`
		Key      interface{} `json:"key"`
		Value    interface{} `json:"value"`
		Document interface{} `json:"doc"`
	} `json:"rows"`
}

// EOF
