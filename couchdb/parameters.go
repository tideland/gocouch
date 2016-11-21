// Tideland Go CouchDB Client - CouchDB - Parameters
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
	"net/http"
	"net/url"
)

//--------------------
// PARAMETER
//--------------------

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(ps *Parameters)

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(ps *Parameters) {
		ps.query.Set("rev", revision)
	}
}

// StartEndKey sets the startkey and endkey for view requests.
func StartEndKey(start, end string) Parameter {
	return func(ps *Parameters) {
		ps.query.Set("startkey", "\""+start+"\"")
		ps.query.Set("endkey", "\""+end+"\"")
	}
}

// IncludeDocuments sets the flag for the including of found view documents.
func SetIncludeDocuments() Parameter {
	return func(ps *Parameters) {
		ps.query.Set("include_docs", "true")
	}
}

//--------------------
// PARAMETERS
//--------------------

// Prameters contains different parameters for the requests to a CouchDB.
type Parameters struct {
	query  url.Values
	header http.Header
}

// newParameters creates a new empty set of parameters.
func newParameters() *Parameters {
	return &Parameters{
		query:  url.Values{},
		header: http.Header{},
	}
}

// apply passes possible parameters to a request.
func (ps *Parameters) apply(req *request, rps ...Parameter) {
	for _, rp := range rps {
		rp(ps)
	}
	if len(ps.query) > 0 {
		req.setQuery(&ps.query)
	}
	if len(ps.header) > 0 {
		req.setHeader(&ps.header)
	}
}

// EOF
