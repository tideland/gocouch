// Tideland Go CouchDB Client - Security - Document Types
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package security

//--------------------
// IMPORTS
//--------------------

import ()

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbAuthentication contains user ID and password
// for authentication.
type couchdbAuthentication struct {
	UserID   string `json:"name"`
	Password string `json:"password"`
}

// couchdRoles contains the roles of a user if the
// authentication succeeded.
type couchdbRoles struct {
	OK     bool     `json:"ok"`
	UserID string   `json:"name"`
	Roles  []string `json:"roles"`
}

// EOF
