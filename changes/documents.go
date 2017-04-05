// Tideland Go CouchDB Client - Changes - Document Types
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package changes

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

// couchdbDocumentIDs contains document identifiers as body
// for the according changes filter.
type couchdbDocumentIDs struct {
	DocumentIDs []string `json:"doc_ids"`
}

// couchdbChanges is a generic result of a CouchDB changes feed.
type couchdbChanges struct {
	LastSequence interface{}           `json:"last_seq"`
	Pending      int                   `json:"pending"`
	Results      couchdbChangesResults `json:"results"`
}

// couchdbChangesResult contains one result of a changes feed.
type couchdbChangesResult struct {
	ID       string                      `json:"id"`
	Sequence interface{}                 `json:"seq"`
	Changes  couchdbChangesResultChanges `json:"changes"`
	Document json.RawMessage             `json:"doc,omitempty"`
	Deleted  bool                        `json:"deleted,omitempty"`
}

type couchdbChangesResults []couchdbChangesResult

// couchdbChangesResultChange contains the revision number of one
// change of one document.
type couchdbChangesResultChange struct {
	Revision string `json:"rev"`
}

type couchdbChangesResultChanges []couchdbChangesResultChange

// EOF
