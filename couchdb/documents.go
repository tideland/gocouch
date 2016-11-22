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

// DesignView defines a view inside a design document.
type DesignView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type DesignViews map[string]DesignView

// DesignAttachment defines an attachment inside a design document.
type DesignAttachment struct {
	Stub        bool   `json:"stub,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Length      int    `json:"length,omitempty"`
}

type DesignAttachments map[string]DesignAttachment

// DesignDocument contains the data of view design documents.
type DesignDocument struct {
	ID                     string            `json:"_id"`
	Revision               string            `json:"_rev,omitempty"`
	Language               string            `json:"language,omitempty"`
	ValidateDocumentUpdate string            `json:"validate_doc_update,omitempty"`
	Views                  DesignViews       `json:"views,omitempty"`
	Shows                  map[string]string `json:"shows,omitempty"`
	Attachments            DesignAttachments `json:"_attachments,omitempty"`
	Signatures             map[string]string `json:"signatures,omitempty"`
	Libraries              interface{}       `json:"libs,omitempty"`
}

// ViewRow contains one row of a view result.
type ViewRow struct {
	ID       string      `json:"id"`
	Key      interface{} `json:"key"`
	Value    interface{} `json:"value"`
	Document interface{} `json:"doc"`
}

type ViewRows []ViewRow

// ViewResult is a generic result of a CouchDB view.
type ViewResult struct {
	TotalRows int      `json:"total_rows"`
	Offset    int      `json:"offset"`
	Rows      ViewRows `json:"rows"`
}

// Result contains internal status information CouchDB returns.
type Result struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
	Error    string `json:"error"`
	Reason   string `json:"reason"`
}

// BulkResults is the list of results after a bulk writing.
type BulkResults []Result

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
