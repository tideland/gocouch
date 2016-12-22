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
// EXTERNAL DOCUMENT TYPES
//--------------------

// User contains user ID and password
// for user management and authentication.
type User struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	UserID   string   `json:"name"`
	Password string   `json:"password"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
}

// UserIDs contains user IDs and roles for
// administrators and users.
type UserIDsRoles struct {
	UserIDs []string `json:"names,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

// Security contains administrators and
// members for one database.
type Security struct {
	Admins  UserIDsRoles `json:"admins,omitempty"`
	Members UserIDsRoles `json:"members,omitempty"`
}

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

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
