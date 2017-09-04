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
	"encoding/json"
	"strconv"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

// Fixed values for some of the view parameters.
const (
	SinceNow = "now"

	StyleMainOnly = "main_only"
	StyleAllDocs  = "all_docs"
)

//--------------------
// PARAMETERS
//--------------------

// Limit sets the maximum number of result rows.
func Limit(limit int) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("limit", strconv.Itoa(limit))
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

// Style sets how many revisions are returned. Default is
// StyleMainOnly only returning the winning document revision.
// StyleAllDocs will return all revision including possible
// conflicts.
func Style(style string) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("style", style)
	}
}

// FilterDocumentIDs sets a filtering of the changes to the
// given document identifiers.
func FilterDocumentIDs(documentIDs ...string) couchdb.Parameter {
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

// FilterSelector sets the filter to the passed selector expression.
func FilterSelector(selector json.RawMessage) couchdb.Parameter {
	update := func(doc interface{}) interface{} {
		// TODO 2017-04-09 Mue Set selector expression.
		return doc
	}
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("filter", "_selector")
		pa.UpdateDocument(update)
	}
}

// FilterView sets the name of a view which map function acts as
// filter in case it emits at least one record.
func FilterView(view string) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("filter", "_view")
		pa.SetQuery("view", view)
	}
}

// EOF
