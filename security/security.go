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

// HasAdministrator checks if a given administrator account exists.
func HasAdministrator(cdb couchdb.CouchDB, session Session, userID string) (bool, error) {
	params := []couchdb.Parameter{}
	if session != nil {
		params = append(params, session.Cookie())
	}
	rs := cdb.Get("/_config/admins/"+userID, nil, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusNotFound {
			return false, nil
		}
		return false, rs.Error()
	}
	return true, nil
}

// WriteAdministrator writes in administrator to the given database.
func WriteAdministrator(cdb couchdb.CouchDB, session Session, userID, password string) error {
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

// ReadSecurity returns the security for the given database.
func ReadSecurity(cdb couchdb.CouchDB, session Session) (*Security, error) {
	path := "/" + cdb.Database() + "/_security"
	rs := cdb.ReadDocument(path, session.Cookie())
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	var security Security
	err := rs.Document(&security)
	if err != nil {
		return nil, err
	}
	return &security, nil
}

// WriteSecurity writes new or changed security data to
// the given database.
func WriteSecurity(cdb couchdb.CouchDB, session Session, security Security) error {
	params := []couchdb.Parameter{}
	if session != nil {
		params = append(params, session.Cookie())
	}
	path := "/" + cdb.Database() + "/_security"
	rs := cdb.Put(path, security, params...)
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
