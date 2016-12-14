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
// USER MANAGEMENT
//--------------------

// UserManagement provides a user and role management for
// a CouchDB.
type UserManagement interface {
	// Create a new user.
	CreateUser(userID, password string) error
}

// userManagement implements the UserManagement interface.
type userManagement struct {
	cdb      couchdb.CouchDB
	userID   string
	password string
}

// NewUserManagement create a user and role management. The
// passed user ID and password are those of an administrator.
func NewUserManagement(cdb couchdb.CouchDB, userID, password string) (UserManagement, error) {
	um := &userManagement{
		cdb:      cdb,
		userID:   userID,
		password: password,
	}
	// Check if the administrator already exists.
	config := map[string]interface{}{}
	rs := um.cdb.Get("/_config", config)
	if rs.IsOK() {
		// No administrator so far.
		rs = um.cdb.Put("/config/admins/"+um.userID, "\""+um.password+"\"")
		if !rs.IsOK() {
			return nil, rs.Error()
		}
	}
	return um, nil
}

// CreateUser implements the UserManagement interface.
func (um *userManagement) CreateUser(userID, password string) error {
	user := &couchdbUser{
		ID:       userDocumentID(userID),
		UserID:   userID,
		Password: password,
		Type:     "user",
	}
	rs := um.cdb.CreateDocument(user, BasicAuthentication(um.userID, um.password))
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
