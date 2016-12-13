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

// couchdbUser contains the data of one user.
type couchdbUser struct {
	ID       string   `json:"_id"`
	UserID   string   `json:"name"`
	Password string   `json:"password"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
}

// couchdRoles contains the roles of a user if the
// authentication succeeded.
type couchdbRoles struct {
	OK       bool     `json:"ok"`
	UserID   string   `json:"name"`
	Password string   `json:"password_sha,omitempty"`
	Salt     string   `json:"salt,omitempty"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
}

// EOF
