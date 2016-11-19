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
	"fmt"
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

// TestCreateDocument tests creating new documents.
func TestCreateDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	docA := DocWithoutID{
		FieldA: "foo",
		FieldB: 4711,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Match(id, "[0-9a-f]{32}")

	docB := DocWithID{
		Identificator: "bar-12345",
		FieldA:        "bar",
		FieldB:        12345,
	}
	resp = cdb.CreateDocument(docB)
	assert.True(resp.IsOK())
	id = resp.ID()
	assert.Equal(id, "bar-12345")

	docC := &IdentifiableDoc{
		FieldA: "yadda",
		FieldB: 54321,
	}
	resp = cdb.CreateDocument(docC)
	assert.True(resp.IsOK())
	id = resp.ID()
	assert.Equal(id, "yadda-54321")
}

// TestReadDocument tests reading new documents.
func TestReadDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	docA := DocWithID{
		Identificator: "foo-12345",
		FieldA:        "foo",
		FieldB:        12345,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := DocWithID{}
	err := resp.ResultValue(&docB)
	assert.Nil(err)
	assert.Equal(docB.Identificator, docA.Identificator)
	assert.Equal(docB.FieldA, docA.FieldA)
	assert.Equal(docB.FieldB, docA.FieldB)
}

//--------------------
// HELPERS
//--------------------

// DocWithoutID is for document tests without an ID.
type DocWithoutID struct {
	FieldA string
	FieldB int
}

// DocWithUD is for document tests with an ID.
type DocWithID struct {
	Identificator string `json:"_id"`
	FieldA        string
	FieldB        int
}

// IdentifiableDoc is for document tests with the Identifiable
// intterface.
type IdentifiableDoc struct {
	revision string
	FieldA   string
	FieldB   int
}

func (d *IdentifiableDoc) DocumentID() string {
	return fmt.Sprintf("%s-%d", d.FieldA, d.FieldB)
}

func (d *IdentifiableDoc) DocumentRevision() string {
	return d.revision
}

// prepareDatabase opens the database deletes a potention test
// database and creates it newly.
func prepareDatabase(assert audit.Assertion) couchdb.CouchDB {
	cfg, err := etc.ReadString(TestingDBCfg)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	resp := cdb.DeleteDatabase()
	resp = cdb.CreateDatabase()
	assert.True(resp.IsOK())
	return cdb
}

// EOF
