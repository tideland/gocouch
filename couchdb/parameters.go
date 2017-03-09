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

import (
	"encoding/json"
	"strconv"
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

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("rev", revision)
	}
}

// Keys sets a number of keys wanted for a view request.
func Keys(keys ...interface{}) Parameter {
	return func(pa Parameterizable) {
		pa.AddKeys(keys...)
	}
}

// StringKeys sets a number of keys of type string wanted for a view request.
func StringKeys(keys ...string) Parameter {
	var ikeys []interface{}
	for _, key := range keys {
		ikeys = append(ikeys, key)
	}
	return Keys(ikeys...)
}

// StartKey sets the startkey for a view request.
func StartKey(start interface{}) Parameter {
	jstart, _ := json.Marshal(start)
	return func(pa Parameterizable) {
		pa.SetQuery("startkey", string(jstart))
	}
}

// EndKey sets the endkey for a view request.
func EndKey(end interface{}) Parameter {
	jend, _ := json.Marshal(end)
	return func(pa Parameterizable) {
		pa.SetQuery("endkey", string(jend))
	}
}

// StartEndKey sets the startkey and endkey for a view request.
func StartEndKey(start, end interface{}) Parameter {
	jstart, _ := json.Marshal(start)
	jend, _ := json.Marshal(end)
	return func(pa Parameterizable) {
		pa.SetQuery("startkey", string(jstart))
		pa.SetQuery("endkey", string(jend))
	}
}

// OneKey reduces a view result to only one emitted key.
func OneKey(key interface{}) Parameter {
	jkey, _ := json.Marshal(key)
	return func(pa Parameterizable) {
		pa.SetQuery("key", string(jkey))
	}
}

// Skip sets the number to skip for view requests.
func Skip(skip int) Parameter {
	return func(pa Parameterizable) {
		if skip > 0 {
			pa.SetQuery("skip", strconv.Itoa(skip))
		}
	}
}

// lIMIT sets the limit for view requests.
func Limit(limit int) Parameter {
	return func(pa Parameterizable) {
		if limit > 0 {
			pa.SetQuery("limit", strconv.Itoa(limit))
		}
	}
}

// SkipLimit sets the number to skip and the limit for
// view requests.
func SkipLimit(skip, limit int) Parameter {
	return func(pa Parameterizable) {
		if skip > 0 {
			pa.SetQuery("skip", strconv.Itoa(skip))
		}
		if limit > 0 {
			pa.SetQuery("limit", strconv.Itoa(limit))
		}
	}
}

// Descending sets the flag for a descending order of found view documents.
func Descending() Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("descending", "true")
	}
}

// NoReduce sets the flag for usage of a reduce function to false.
func NoReduce() Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("reduce", "false")
	}
}

// Group sets the flag for grouping including the level for the
// reduce function.
func Group(level int) Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("group", "true")
		if level > 0 {
			pa.SetQuery("group_level", strconv.Itoa(level))
		}
	}
}

// IncludeDocuments sets the flag for the including of found view documents.
func IncludeDocuments() Parameter {
	return func(pa Parameterizable) {
		pa.SetQuery("include_docs", "true")
	}
}

// EOF
