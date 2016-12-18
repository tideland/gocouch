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
	cdb := prepareCouchDB("administrator", assert)

	// Check first admin before it exists.
	ok, err := security.HasAdministrator(cdb, nil, "admin1")
	assert.Nil(err)
	assert.False(ok)

	err = security.WriteAdministrator(cdb, nil, "admin1", "admin1")
	assert.Nil(err)
	defer func() {
		// Let the administator remove himself.
		session, err := security.NewSession(cdb, "admin1", "admin1")
		assert.Nil(err)
		err = security.DeleteAdministrator(cdb, session, "admin1")
		assert.Nil(err)
	}()

	// Check first admin after creation without session.
	ok, err = security.HasAdministrator(cdb, nil, "admin1")
	assert.ErrorMatch(err, ".*status code 401.*")
	assert.False(ok)

	// Check first admin after creation with session.
	session, err := security.NewSession(cdb, "admin1", "admin1")
	assert.Nil(err)
	ok, err = security.HasAdministrator(cdb, session, "admin1")
	assert.Nil(err)
	assert.True(ok)

	// Now care for second administrator, first withour session,
	// then with.
	err = security.WriteAdministrator(cdb, nil, "admin2", "admin2")
	assert.ErrorMatch(err, ".*status code 401.*")

	err = security.WriteAdministrator(cdb, session, "admin2", "admin2")
	assert.Nil(err)

	ok, err = security.HasAdministrator(cdb, session, "admin2")
	assert.Nil(err)
	assert.True(ok)

	err = security.DeleteAdministrator(cdb, session, "admin2")
	assert.Nil(err)

	ok, err = security.HasAdministrator(cdb, session, "admin2")
	assert.Nil(err)
	assert.False(ok)
}

// TestSecurity tests the security related functions.
func TestSecurity(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareCouchDB("security", assert)

	// Without database and admin.
	in := security.Security{
		Admins: security.UserIDsRoles{
			UserIDs: []string{"admin"},
		},
	}
	err := security.WriteSecurity(cdb, nil, in)
	assert.ErrorMatch(err, ".*status code 404.*")

	// With database and without admin.
	rs := cdb.CreateDatabase()
	assert.True(rs.IsOK())
	defer func() {
		cdb.DeleteDatabase()
	}()
	err = security.WriteSecurity(cdb, nil, in)
	assert.Nil(err)
}

//--------------------
// HELPERS
//--------------------

// MyDocument is used for the tests.
type MyDocument struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Name        string `json:"name"`
	Age         int    `json:"age"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
}

// prepareCouchDB opens the DBMS for one database
// w/o creating it. It deletes the named database to
// avoid conflicts.
func prepareCouchDB(database string, assert audit.Assertion) couchdb.CouchDB {
	cfgstr := strings.Replace(TemplateDBcfg, "<<DATABASE>>", database, 1)
	cfg, err := etc.ReadString(cfgstr)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	cdb.DeleteDatabase()
	return cdb
}

// EOF
