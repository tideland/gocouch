// Tideland Go CouchDB Client - CouchDB - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

const (
	EmptyCfg       = "{etc}"
	LocalhostCfg   = "{etc {hostname localhost}{port 5984}}"
	TestingDBCfg   = "{etc {hostname localhost}{port 5984}{database tgocouch-testing}}"
	TemporaryDBCfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-temporary}}"
)

//--------------------
// TESTS
//--------------------

// TestNoConfig tests opening the database without a configuration.
func TestNoConfig(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cdb, err := couchdb.Open(nil)
	assert.ErrorMatch(err, ".* cannot open database without configuration")
	assert.Nil(cdb)
}

// TestAllDatabases tests the retrieving of all databases.
func TestAllDatabases(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := etc.ReadString(EmptyCfg)
	assert.Nil(err)

	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	ids, err := cdb.AllDatabases()
	assert.Nil(err)
	assert.True(len(ids) != 0)
}

// TestCreateDeleteDatabase tests the creation and deletion
// of a database.
func TestCreateDeleteDatabase(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := etc.ReadString(TemporaryDBCfg)
	assert.Nil(err)

	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	ids, err := cdb.AllDatabases()
	assert.Nil(err)
	dbNo := len(ids)
	resp := cdb.CreateDatabase()
	assert.True(resp.IsOK())
	ids, err = cdb.AllDatabases()
	assert.Nil(err)
	assert.Equal(len(ids), dbNo+1)
	resp = cdb.DeleteDatabase()
	assert.True(resp.IsOK())
	ids, err = cdb.AllDatabases()
	assert.Nil(err)
	assert.Equal(len(ids), dbNo)
}

// EOF
