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

// TestNewUserManagement tests the starting of the
// user management and the creation of an administrator
// if needed.
func TestNewUserManagement(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("new-user-management", assert)
	defer cleanup()

	um, err := security.NewUserManagement(cdb, "administrator", "administrator")
	assert.Nil(err)
	assert.NotNil(um)
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

// prepareDatabase opens the database, deletes a possible test
// database, and creates it newly.
func prepareDatabase(database string, assert audit.Assertion) (couchdb.CouchDB, func()) {
	cfgstr := strings.Replace(TemplateDBcfg, "<<DATABASE>>", database, 1)
	cfg, err := etc.ReadString(cfgstr)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	resp := cdb.DeleteDatabase()
	resp = cdb.CreateDatabase()
	assert.True(resp.IsOK())
	return cdb, func() { cdb.DeleteDatabase() }
}

// EOF
