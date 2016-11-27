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
)

//--------------------
// PARAMETERIZABLE
//--------------------

// KeyValue is used for generic query and header parameters.
type KeyValue struct {
	Key   string
	Value string
}

// Parameterizable defines the methods needed to apply the parameters.
type Parameterizable interface {
	// SetQuery sets a query parameter.
	SetQuery(key, value string)

	// AddQuery adds a query parameter to an existing one.
	AddQuery(key, value string)

	// SetHeader sets a header parameter.
	SetHeader(key, value string)

	// AddKeys adds view key parameters.
	AddKeys(keys ...interface{})
}

//--------------------
// PARAMETER
//--------------------

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(pa Parameterizable)

// Query is generic for setting request query parameters.
func Query(kvs ...KeyValue) Parameter {
	return func(pa Parameterizable) {
		for _, kv := range kvs {
			pa.AddQuery(kv.Key, kv.Value)
		}
	}
}

// Header is generic for setting request header parameters.
func Header(kvs ...KeyValue) Parameter {
	return func(pa Parameterizable) {
		for _, kv := range kvs {
			pa.SetHeader(kv.Key, kv.Value)
		}
	}
}

// BasicAuthentication is intended for basic authentication
// against the database.
func BasicAuthentication(userID, password string) Parameter {
	return func(pa Parameterizable) {
		up := []byte(userID + ":" + password)
		auth := "Basic " + base64.StdEncoding.EncodeToString(up)

		pa.SetHeader("Authorization", auth)
	}
}

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("rev", revision)
	}
}

// Keys sets a number of keys wanted from a view request.
func Keys(keys ...interface{}) Parameter {
	return func(pa Parameterizable) {
		pa.AddKeys(keys...)
	}
}

// StartEndKey sets the startkey and endkey for view requests.
func StartEndKey(start, end string) Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("startkey", "\""+start+"\"")
		pa.SetQuery("endkey", "\""+end+"\"")
	}
}

// IncludeDocuments sets the flag for the including of found view documents.
func SetIncludeDocuments() Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("include_docs", "true")
	}
}

// EOF
