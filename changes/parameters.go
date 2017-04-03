// Tideland Go CouchDB Client - Changes - Parameters
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
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

const (
	SinceNow = "now"
)

//--------------------
// PARAMETERS
//--------------------

// DecomentIDs sets a filtering of the changes to the
// given document identifiers.
func DocumentIDs(documentIDs ...string) couchdb.Parameter {
	update := func(doc interface{}) interface{} {
		if doc == nil {
			doc = &couchdbDocumentIDs{}
		}
		idsdoc, ok := doc.(*couchdbDocumentIDs)
		if ok {
			idsdoc.DocumentIDs = append(idsdoc.DocumentIDs, documentIDs...)
			return idsdoc
		}
		return doc
	}
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("filter", "_doc_ids")
		pa.UpdateDocument(update)
	}
}

// Since sets the start of the changes gathering, can also be "now".
func Since(sequence string) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("since", sequence)
	}
}

// Descending sets the flag for a descending order of changes.
func Descending() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("descending", "true")
	}
}

// EOF
