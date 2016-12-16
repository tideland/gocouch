// Tideland Go CouchDB Client - Security
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
// SECURITY FUNCTIONS
//--------------------

// CreateAdministrator adds an administrator to the given database.
func CreateAdministrator(cdb couchdb.CouchDB, session Session, userID, password string) error {
	params := []couchdb.Parameter{}
	if session != nil {
		params = append(params, session.Cookie())
	}
	rs := cdb.Put("/_config/admins/"+userID, password, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteAdministrator deletes an administrator from the given database.
func DeleteAdministrator(cdb couchdb.CouchDB, session Session, userID string) error {
	rs := cdb.Delete("/_config/admins/"+userID, nil, session.Cookie())
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// CreateUser adds a user to the given database.
func CreateUser(cdb couchdb.CouchDB, session Session, userID, password string) error {
	user := &couchdbUser{
		ID:       userDocumentID(userID),
		UserID:   userID,
		Password: password,
		Type:     "user",
	}
	rs := cdb.CreateDocument(user, session.Cookie())
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
