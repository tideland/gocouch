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
func HasAdministrator(cdb couchdb.CouchDB, userID string, params ...couchdb.Parameter) (bool, error) {
	rs := cdb.Get("/_config/admins/"+userID, nil, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusNotFound {
			return false, nil
		}
		return false, rs.Error()
	}
	return true, nil
}

// WriteAdministrator adds or updates an administrator to the given database.
func WriteAdministrator(cdb couchdb.CouchDB, userID, password string, params ...couchdb.Parameter) error {
	rs := cdb.Put("/_config/admins/"+userID, password, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteAdministrator deletes an administrator from the given database.
func DeleteAdministrator(cdb couchdb.CouchDB, userID string, params ...couchdb.Parameter) error {
	rs := cdb.Delete("/_config/admins/"+userID, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// HasUser checks if a given user account exists.
func HasUser(cdb couchdb.CouchDB, userID string, params ...couchdb.Parameter) (bool, error) {
	id := userDocumentID(userID)
	rs := cdb.Get("/_users/"+id, nil, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusNotFound {
			return false, nil
		}
		return false, rs.Error()
	}
	return true, nil
}

// WriteUser adds or updates a user to the system.
func WriteUser(cdb couchdb.CouchDB, userID, password string, roles []string, params ...couchdb.Parameter) error {
	user := &couchdbUser{
		ID:       userDocumentID(userID),
		UserID:   userID,
		Password: password,
		Type:     "user",
		Roles:    roles,
	}
	rs := cdb.Put("/_users/"+user.ID, user, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusConflict {
			return errors.New(ErrUserExists, errorMessages)
		}
		return rs.Error()
	}
	return nil
}

// DeleteUser deletes a user from the system.
func DeleteUser(cdb couchdb.CouchDB, userID string, params ...couchdb.Parameter) error {
	id := userDocumentID(userID)
	rs := cdb.Delete("/_users/"+id, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// ReadSecurity returns the security for the given database.
func ReadSecurity(cdb couchdb.CouchDB, params ...couchdb.Parameter) (*Security, error) {
	path := "/" + cdb.Database() + "/_security"
	rs := cdb.Get(path, nil, params...)
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
func WriteSecurity(cdb couchdb.CouchDB, security Security, params ...couchdb.Parameter) error {
	path := "/" + cdb.Database() + "/_security"
	rs := cdb.Put(path, security, params...)
	if !rs.IsOK() {
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
