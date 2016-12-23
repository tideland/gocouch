// Tideland Go CouchDB Client - Security - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package security_test

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"

	"github.com/tideland/gocouch/couchdb"
	"github.com/tideland/gocouch/security"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TemplateDBcfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}}"
)

//--------------------
// TESTS
//--------------------

// TestAdministraotor tests the administrator related functions.
func TestAdministrator(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase("administrator", assert)

	// Check first admin before it exists.
	ok, err := security.HasAdministrator(cdb, "admin1")
	assert.Nil(err)
	assert.False(ok)

	err = security.WriteAdministrator(cdb, "admin1", "admin1")
	assert.Nil(err)
	defer func() {
		// Let the administator remove himself.
		session, err := security.NewSession(cdb, "admin1", "admin1")
		assert.Nil(err)
		err = security.DeleteAdministrator(cdb, "admin1", session.Cookie())
		assert.Nil(err)
	}()

	// Check first admin after creation without session.
	ok, err = security.HasAdministrator(cdb, "admin1")
	assert.ErrorMatch(err, ".*status code 401.*")
	assert.False(ok)

	// Check first admin after creation with session.
	session, err := security.NewSession(cdb, "admin1", "admin1")
	assert.Nil(err)
	ok, err = security.HasAdministrator(cdb, "admin1", session.Cookie())
	assert.Nil(err)
	assert.True(ok)

	// Now care for second administrator, first withour session,
	// then with.
	err = security.WriteAdministrator(cdb, "admin2", "admin2")
	assert.ErrorMatch(err, ".*status code 401.*")

	err = security.WriteAdministrator(cdb, "admin2", "admin2", session.Cookie())
	assert.Nil(err)

	ok, err = security.HasAdministrator(cdb, "admin2", session.Cookie())
	assert.Nil(err)
	assert.True(ok)

	auth := security.BasicAuthentication("admin1", "admin1")
	err = security.DeleteAdministrator(cdb, "admin2", auth)
	assert.Nil(err)

	ok, err = security.HasAdministrator(cdb, "admin2", auth)
	assert.Nil(err)
	assert.False(ok)
}

// TestUser tests the user management related functions.
func TestUser(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase("security", assert)

	userA, err := security.ReadUser(cdb, "user1")
	assert.Nil(userA)
	assert.ErrorMatch(err, ".*status code 404.*")

	userB := &security.User{
		UserID:   "user1",
		Password: "user1",
		Roles:    []string{"developer"},
	}
	err = security.CreateUser(cdb, userB)
	assert.Nil(err)

	userA, err = security.ReadUser(cdb, "user1")
	assert.Nil(err)
	assert.Equal(userA.UserID, "user1")

	err = security.DeleteUser(cdb, userA)
	assert.Nil(err)
}

// TestSecurity tests the security related functions.
func TestSecurity(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase("security", assert)

	// Without database and authentication.
	in := security.Security{
		Admins: security.UserIDsRoles{
			UserIDs: []string{"admin"},
		},
	}
	err := security.WriteSecurity(cdb, in)
	assert.ErrorMatch(err, ".*status code 404.*")

	// Without database but with authentication.
	err = security.WriteAdministrator(cdb, "admin", "admin")
	assert.Nil(err)
	defer func() {
		// Let the administator remove himself.
		session, err := security.NewSession(cdb, "admin", "admin")
		assert.Nil(err)
		err = security.DeleteAdministrator(cdb, "admin", session.Cookie())
		assert.Nil(err)
	}()
	session, err := security.NewSession(cdb, "admin", "admin")
	assert.Nil(err)
	err = security.WriteSecurity(cdb, in, session.Cookie())
	assert.ErrorMatch(err, ".*status code 404.*")

	// With database and without authentication.
	rs := cdb.CreateDatabase()
	assert.ErrorMatch(rs.Error(), ".*status code 401.*")
	rs = cdb.CreateDatabase(session.Cookie())
	assert.True(rs.IsOK())
	defer func() {
		rs := cdb.DeleteDatabase(session.Cookie())
		assert.True(rs.IsOK())
	}()
	err = security.WriteSecurity(cdb, in)
	assert.ErrorMatch(err, ".*status code 401.*")

	// With database and authentication.
	err = security.WriteSecurity(cdb, in, session.Cookie())
	assert.Nil(err)

	// Now read the security information.
	out, err := security.ReadSecurity(cdb, security.BasicAuthentication("admin", "admin"))
	assert.Nil(err)
	assert.Equal(out.Admins, in.Admins)
}

// TestScenario tests a scenario with administrator, user,
// amd unauthorized access to the database.
func TestScenario(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase("scenario", assert)

	// Create administrator.
	err := security.WriteAdministrator(cdb, "admin", "admin")
	assert.Nil(err)
	session, err := security.NewSession(cdb, "admin", "admin")
	assert.Nil(err)
	defer func() {
		// Let the administator remove himself.
		err = security.DeleteAdministrator(cdb, "admin", session.Cookie())
		assert.Nil(err)
	}()

	// Create users.
	user := &security.User{
		UserID:   "user",
		Password: "user",
		Roles:    []string{"developer"},
	}
	err = security.CreateUser(cdb, user, session.Cookie())
	assert.Nil(err)
	user = &security.User{
		UserID:   "somebody",
		Password: "somebody",
		Roles:    []string{"visitors"},
	}
	err = security.CreateUser(cdb, user, session.Cookie())
	assert.Nil(err)
	defer func() {
		user, err := security.ReadUser(cdb, "user", session.Cookie())
		assert.Nil(err)
		err = security.DeleteUser(cdb, user, session.Cookie())
		assert.Nil(err)
		user, err = security.ReadUser(cdb, "somebody", session.Cookie())
		assert.Nil(err)
		err = security.DeleteUser(cdb, user, session.Cookie())
		assert.Nil(err)
	}()

	// Create database.
	rs := cdb.CreateDatabase(session.Cookie())
	assert.True(rs.IsOK())
	defer func() {
		rs := cdb.DeleteDatabase(session.Cookie())
		assert.True(rs.IsOK())
	}()

	// Add security.
	sec := security.Security{
		Admins: security.UserIDsRoles{
			UserIDs: []string{"user"},
		},
	}
	err = security.WriteSecurity(cdb, sec, session.Cookie())
	assert.Nil(err)

	// Add document.
	doc := MyDocument{
		DocumentID: "foo",
		Data:       "foo",
	}
	rs = cdb.CreateDocument(doc)
	assert.False(rs.IsOK())
	rs = cdb.CreateDocument(doc, security.BasicAuthentication("somebody", "somebody"))
	assert.False(rs.IsOK())
	rs = cdb.CreateDocument(doc, security.BasicAuthentication("user", "user"))
	assert.True(rs.IsOK())
}

//--------------------
// HELPERS
//--------------------

// MyDocument is used for the tests.
type MyDocument struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Data string `json:"data"`
}

// prepareDatabase opens the database and deletes a
// possible test database.
func prepareDatabase(database string, assert audit.Assertion) couchdb.CouchDB {
	cfgstr := strings.Replace(TemplateDBcfg, "<<DATABASE>>", database, 1)
	cfg, err := etc.ReadString(cfgstr)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	cdb.DeleteDatabase()
	return cdb
}

// EOF
