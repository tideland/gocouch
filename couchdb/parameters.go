// Tideland Go CouchDB Client - CouchDB - Parameters
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

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

	// UpdateDocument allows to modify or exchange the document.
	UpdateDocument(update func(interface{}) interface{})
}

//--------------------
// PARAMETERS
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

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("rev", revision)
	}
}

// EOF
