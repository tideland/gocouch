// Tideland Go CouchDB Client - Views - Document Types
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
	"encoding/json"
)

//--------------------
// EXTERNAL DOCUMENT TYPES
//--------------------

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbKeys sets key constraints for view requests.
type couchdbKeys struct {
	Keys []interface{} `json:"keys"`
}

// couchdbViewResult is a generic result of a CouchDB view.
type couchdbViewResult struct {
	TotalRows int             `json:"total_rows"`
	Offset    int             `json:"offset"`
	Rows      couchdbViewRows `json:"rows"`
}

// couchdbViewRow contains one row of a view result.
type couchdbViewRow struct {
	ID       string          `json:"id"`
	Key      json.RawMessage `json:"key"`
	Value    json.RawMessage `json:"value"`
	Document json.RawMessage `json:"doc"`
}

type couchdbViewRows []couchdbViewRow

// EOF
