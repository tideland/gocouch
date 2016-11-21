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

// View defines a view inside a design document.
type View struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type Views map[string]View

// Attachment defines an attachment inside a design document.
type Attachment struct {
	Stub        bool   `json:"stub,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Length      int    `json:"length,omitempty"`
}

type Attachments map[string]Attachment

// DesignDocument contains the data of view design documents.
type DesignDocument struct {
	ID                     string            `json:"_id"`
	Revision               string            `json:"_rev,omitempty"`
	Language               string            `json:"language,omitempty"`
	ValidateDocumentUpdate string            `json:"validate_doc_update,omitempty"`
	Views                  Views             `json:"views,omitempty"`
	Shows                  map[string]string `json:"shows,omitempty"`
	Attachments            Attachments       `json:"_attachments,omitempty"`
	Signatures             map[string]string `json:"signatures,omitempty"`
	Libraries              interface{}       `json:"libs,omitempty"`
}

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbResponse contains response information CouchDB returns.
type couchdbResponse struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
	Error    string `json:"error"`
	Reason   string `json:"reason"`
}

// couchdbViewKeys sets key constraints for view requests.
type couchdbViewKeys struct {
	Keys []interface{} `json:"keys"`
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

// idAndRevision is used to simply retrieve ID and revision of
// a document.
type idAndRevision struct {
	ID       string `json:"_id"`
	Revision string `json:"_rev"`
}

// EOF
