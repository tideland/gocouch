// Tideland Go CouchDB Client - Security - Users
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
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// USERS
//--------------------

// Users provides a user and role management for
// a CouchDB.
type Users interface {
	// Create a new user.
	CreateUser() error
}

// users implements the Users interface.
type users struct {
	cdb couchdb.CouchDB
}

func NewUsers(cdb couchdb.CouchDB) (Users, error) {
	u := &users{
		cdb: cdb,
	}
	return u, nil
}

// CreateUser implements the Users interface.
func (u *users) CreateUser() error {
	return nil
}

// EOF
