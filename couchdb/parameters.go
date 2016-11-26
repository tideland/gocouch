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
	"encoding/base64"
	"net/http"
	"net/url"
)

//--------------------
// PARAMETER
//--------------------

// KeyValue is used for the generic query and header parameters.
type KeyValue struct {
	Key   string
	Value string
}

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(ps *Parameters)

// Query is generic for setting request query parameters.
func Query(kvs ...KeyValue) Parameter {
	return func(ps *Parameters) {
		for _, kv := range kvs {
			ps.query.Add(kv.Key, kv.Value)
		}
	}
}

// Header is generic for setting request header parameters.
func Header(kvs ...KeyValue) Parameter {
	return func(ps *Parameters) {
		for _, kv := range kvs {
			ps.header.Set(kv.Key, kv.Value)
		}
	}
}

// BasicAuthentication is intended for basic authentication
// against the database.
func BasicAuthentication(userID, password string) Parameter {
	return func(ps *Parameters) {
		up := []byte(userID + ":" + password)
		auth := "Basic " + base64.StdEncoding.EncodeToString(up)

		ps.header.Set("Authorization", auth)
	}
}

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(ps *Parameters) {
		ps.query.Set("rev", revision)
	}
}

// Keys sets a number of keys wanted from a view request.
func Keys(keys ...interface{}) Parameter {
	return func(ps *Parameters) {
		ps.keys = append(ps.keys, keys...)
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
	cdbps  []Parameter
	keys   []interface{}
	query  url.Values
	header http.Header
}

// newParameters creates a new empty set of parameters.
func newParameters(cdbps []Parameter) *Parameters {
	return &Parameters{
		cdbps:  cdbps,
		query:  url.Values{},
		header: http.Header{},
	}
}

// apply passes possible parameters to a request.
func (ps *Parameters) apply(req *request, rps ...Parameter) {
	rps = append(ps.cdbps, rps...)
	for _, rp := range rps {
		rp(ps)
	}
	if len(ps.keys) > 0 {
		req.keys = ps.keys
	}
	if len(ps.query) > 0 {
		req.query = ps.query
	}
	if len(ps.header) > 0 {
		req.header = ps.header
	}
}

// EOF
