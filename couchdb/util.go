// Tideland Go CouchDB Client - CouchDB - Utilities
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
	"net/url"
)

//--------------------
// QUERY
//--------------------

// Query allows the easy creation of URL values for view queries.
type Query struct {
	values url.Values
}

// NewQuery creates an empty query.
func NewQuery() *Query {
	return &Query{
		values: url.Values{},
	}
}

// StartEndKey creates startkey and endkey.
func (q *Query) StartEndKey(start, end string) *Query {
	q.values.Set("startkey", "\""+start+"\"")
	q.values.Set("endkey", "\""+end+"\"")
	return q
}

// IncludeDocuments sets the flag for the including of the
// found documents.
func (q *Query) IncludeDocuments() *Query {
	q.values.Set("include_docs", "true")
	return q
}

// Encode the query for a URL.
func (q *Query) Encode() string {
	return q.values.Encode()
}

// EOF
