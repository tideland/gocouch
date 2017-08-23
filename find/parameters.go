// Tideland Go CouchDB Client - Find - Parameters
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find

//--------------------
// PARAMETERIZABLE
//--------------------

// Parameterizable defines the methods needed to apply the parameters.
type Parameterizable interface {
	// SetParameter sets one of the request parameters.
	SetParameter(key string, parameter interface{})
}

//--------------------
// PARAMETERS
//--------------------

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(pa Parameterizable)

// Fields specifies which fields of each object should be returned.
// if it's omitted, the entire object is returned.
func Fields(fields ...string) Parameter {
	return func(pa Parameterizable) {
		pa.SetParameter("fields", fields)
	}
}

// Limit sets the maximum number of results returned. Default
// by database is 25.
func Limit(limit int) Parameter {
	return func(pa Parameterizable) {
		pa.SetParameter("limit", limit)
	}
}

// EOF
