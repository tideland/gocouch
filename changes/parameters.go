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
// PARAMETERS
//--------------------

// Descending sets the flag for a descending order of changes.
func Descending() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetQuery("descending", "true")
	}
}

// EOF
