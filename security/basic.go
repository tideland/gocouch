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
func BasicAuthentication(userID, password string) couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		up := []byte(userID + ":" + password)
		auth := "Basic " + base64.StdEncoding.EncodeToString(up)

		pa.SetHeader("Authorization", auth)
	}
}

// EOF
