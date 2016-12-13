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
	"github.com/tideland/golib/errors"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// USERS
//--------------------

// Users provides a user and role management for
// a CouchDB.
type Users interface {
	// Create a new user.
	CreateUser(userID, password string) error
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
func (u *users) CreateUser(userID, password string) error {
	user := &couchdbUser{
		ID:       userDocumentID(userID),
		UserID:   userID,
		Password: password,
		Type:     "user",
	}
	rs := u.cdb.CreateDocument(user)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusConflict {
			return errors.New(ErrUserExists, errorMessages)
		}
		return rs.Error()
	}
	return nil
}

// userDocumentID builds the document ID based
// on the user ID.
func userDocumentID(userID string) string {
	return "org.couchdb.user:" + userID
}

// EOF
