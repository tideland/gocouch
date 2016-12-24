// Tideland Go CouchDB Client - Security - Basic Authentication
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package security

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// PARAMETER
//--------------------

// BasicAuthentication is intended for basic authentication
// against the database.
func BasicAuthentication(name, password string) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		np := []byte(name + ":" + password)
		auth := "Basic " + base64.StdEncoding.EncodeToString(np)

		pa.SetHeader("Authorization", auth)
	}
}

// EOF
