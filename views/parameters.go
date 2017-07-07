// Tideland Go CouchDB Client - Views - couchdb.Parameters
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package views

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"strconv"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// PARAMETERS
//--------------------

// Keys sets a number of keys wanted for a view request.
func Keys(keys ...interface{}) couchdb.Parameter {
	update := func(doc interface{}) interface{} {
		if doc == nil {
			doc = &couchdbKeys{}
		}
		kdoc, ok := doc.(*couchdbKeys)
		if ok {
			kdoc.Keys = append(kdoc.Keys, keys...)
			return kdoc
		}
		return doc
	}
	return func(pa couchdb.Parameterizable) {
		pa.UpdateDocument(update)
	}
}

// StringKeys sets a number of keys of type string wanted for a view request.
func StringKeys(keys ...string) couchdb.Parameter {
	var ikeys []interface{}
	for _, key := range keys {
		ikeys = append(ikeys, key)
	}
	return Keys(ikeys...)
}

// StartKey sets the startkey for a view request.
func StartKey(start interface{}) couchdb.Parameter {
	jstart, _ := json.Marshal(start)
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("startkey", string(jstart))
	}
}

// EndKey sets the endkey for a view request.
func EndKey(end interface{}) couchdb.Parameter {
	jend, _ := json.Marshal(end)
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("endkey", string(jend))
	}
}

// StartEndKey sets the startkey and endkey for a view request.
func StartEndKey(start, end interface{}) couchdb.Parameter {
	jstart, _ := json.Marshal(start)
	jend, _ := json.Marshal(end)
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("startkey", string(jstart))
		pa.SetQuery("endkey", string(jend))
	}
}

// OneKey reduces a view result to only one emitted key.
func OneKey(key interface{}) couchdb.Parameter {
	jkey, _ := json.Marshal(key)
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("key", string(jkey))
	}
}

// Skip sets the number to skip for view requests.
func Skip(skip int) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		if skip > 0 {
			pa.SetQuery("skip", strconv.Itoa(skip))
		}
	}
}

// Limit sets the limit for view requests.
func Limit(limit int) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		if limit > 0 {
			pa.SetQuery("limit", strconv.Itoa(limit))
		}
	}
}

// SkipLimit sets the number to skip and the limit for
// view requests.
func SkipLimit(skip, limit int) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		if skip > 0 {
			pa.SetQuery("skip", strconv.Itoa(skip))
		}
		if limit > 0 {
			pa.SetQuery("limit", strconv.Itoa(limit))
		}
	}
}

// Descending sets the flag for a descending order of found view documents.
func Descending() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("descending", "true")
	}
}

// NoReduce sets the flag for usage of a reduce function to false.
func NoReduce() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("reduce", "false")
	}
}

// Group sets the flag for grouping including the level for the
// reduce function.
func Group(level int) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("group", "true")
		if level > 0 {
			pa.SetQuery("group_level", strconv.Itoa(level))
		}
	}
}

// IncludeDocuments sets the flag for the including of found view documents.
func IncludeDocuments() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("include_docs", "true")
	}
}

//--------------------
// HELPER FUNCTIONS
//--------------------

// ComplexKey simply combines individual values to a combined
// key as slice of those values. It's only for making the code
// more readable.
func ComplexKey(values ...interface{}) interface{} {
	return values
}

// EOF
