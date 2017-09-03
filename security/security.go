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
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// SECURITY FUNCTIONS
//--------------------

// HasAdministrator checks if a given administrator account exists.
func HasAdministrator(cdb couchdb.CouchDB, nodename, name string, params ...couchdb.Parameter) (bool, error) {
	path := cdb.Path("_node", nodename, "_config", "admins", name)
	rs := cdb.Get(path, nil, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == couchdb.StatusNotFound {
			return false, nil
		}
		return false, rs.Error()
	}
	return true, nil
}

// WriteAdministrator adds or updates an administrator to the given database.
func WriteAdministrator(cdb couchdb.CouchDB, nodename, name, password string, params ...couchdb.Parameter) error {
	path := cdb.Path("_node", nodename, "_config", "admins", name)
	rs := cdb.Put(path, password, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteAdministrator deletes an administrator from the given database.
func DeleteAdministrator(cdb couchdb.CouchDB, nodename, name string, params ...couchdb.Parameter) error {
	path := cdb.Path("_node", nodename, "_config", "admins", name)
	rs := cdb.Delete(path, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// CreateUser adds a new user to the system.
func CreateUser(cdb couchdb.CouchDB, user *User, params ...couchdb.Parameter) error {
	if err := ensureUsersDatabase(cdb); err != nil {
		return err
	}
	user.DocumentID = userDocumentID(user.Name)
	user.Type = "user"
	path := cdb.Path("_users", user.DocumentID)
	rs := cdb.Put(path, user, params...)
	return rs.Error()
}

// ReadUser reads an existing user from the system.
func ReadUser(cdb couchdb.CouchDB, name string, params ...couchdb.Parameter) (*User, error) {
	path := cdb.Path("_users", userDocumentID(name))
	rs := cdb.Get(path, nil, params...)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	var user User
	err := rs.Document(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user in the system.
func UpdateUser(cdb couchdb.CouchDB, user *User, params ...couchdb.Parameter) error {
	if err := ensureUsersDatabase(cdb); err != nil {
		return err
	}
	path := cdb.Path("_users", user.DocumentID)
	rs := cdb.Put(path, user, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteUser deletes a user from the system.
func DeleteUser(cdb couchdb.CouchDB, user *User, params ...couchdb.Parameter) error {
	params = append(params, couchdb.Revision(user.DocumentRevision))
	path := cdb.Path("_users", user.DocumentID)
	rs := cdb.Delete(path, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// ReadSecurity returns the security for the given database.
func ReadSecurity(cdb couchdb.CouchDB, params ...couchdb.Parameter) (*Security, error) {
	path := cdb.DatabasePath("_security")
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
	path := cdb.DatabasePath("_security")
	rs := cdb.Put(path, security, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

//--------------------
// HELPERS
//--------------------

// ensureUsersDatabase checks if the _users database exists and
// creates it if needed.
func ensureUsersDatabase(cdb couchdb.CouchDB) error {
	rs := cdb.Get("_users", nil)
	if rs.IsOK() {
		return nil
	}
	return cdb.Put("_users", nil).Error()
}

// userDocumentID builds the document ID based
// on the name.
func userDocumentID(name string) string {
	return "org.couchdb.user:" + name
}

// EOF
