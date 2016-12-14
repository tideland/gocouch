// Tideland Go CouchDB Client - Security - Administration
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
// ADMINISTRATION
//--------------------

// Administration provides a user and role management for
// a CouchDB.
type Administration interface {
	// Create a new user.
	CreateUser(userID, password string) error
}

// administration implements the Administration interface.
type administration struct {
	cdb      couchdb.CouchDB
	userID   string
	password string
}

// NewAdministration creates a user administration. The
// passed user ID and password combination is an administrator.
// If no administrator exists so far it will be created.
func NewAdministration(cdb couchdb.CouchDB, userID, password string) (Administration, error) {
	a := &administration{
		cdb:      cdb,
		userID:   userID,
		password: password,
	}
	// Check if the administrator already exists.
	config := map[string]interface{}{}
	rs := a.cdb.Get("/_config", config)
	if rs.IsOK() {
		// No administrator so far.
		rs = a.cdb.Put("/_config/admins/"+a.userID, a.password)
		if !rs.IsOK() {
			return nil, rs.Error()
		}
	}
	return a, nil
}

// CreateUser implements the Administration interface.
func (a *administration) CreateUser(userID, password string) error {
	user := &couchdbUser{
		ID:       userDocumentID(userID),
		UserID:   userID,
		Password: password,
		Type:     "user",
	}
	rs := a.cdb.CreateDocument(user, BasicAuthentication(a.userID, a.password))
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
