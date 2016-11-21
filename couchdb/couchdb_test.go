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

// TestCreateDocument tests creating new documents.
func TestCreateDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	// Create document without ID.
	docA := MyDocument{
		FieldA: "foo",
		FieldB: 4711,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Match(id, "[0-9a-f]{32}")

	// create document with ID.
	docB := MyDocument{
		DocumentID: "bar-12345",
		FieldA:     "bar",
		FieldB:     12345,
	}
	resp = cdb.CreateDocument(docB)
	assert.True(resp.IsOK())
	id = resp.ID()
	assert.Equal(id, "bar-12345")
}

// TestReadDocument tests reading a document.
func TestReadDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		FieldA:     "foo",
		FieldB:     12345,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.ResultValue(&docB)
	assert.Nil(err)
	assert.Equal(docB.DocumentID, docA.DocumentID)
	assert.Equal(docB.FieldA, docA.FieldA)
	assert.Equal(docB.FieldB, docA.FieldB)

	// Try to read non-existant document.
	resp = cdb.ReadDocument("i-do-not-exist")
	assert.False(resp.IsOK())
	assert.ErrorMatch(resp.Error(), ".* 404,.*")
}

// TestUpdateDocument tests updating documents.
func TestUpdateDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	// Create first revision.
	docA := MyDocument{
		DocumentID: "foo-12345",
		FieldA:     "foo",
		FieldB:     12345,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	revision := resp.Revision()
	assert.Equal(id, "foo-12345")

	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.ResultValue(&docB)
	assert.Nil(err)

	// Update the document.
	docB.FieldB = 54321

	resp = cdb.UpdateDocument(docB)
	assert.True(resp.IsOK())

	// read the updated revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docC := MyDocument{}
	err = resp.ResultValue(&docC)
	assert.Nil(err)
	assert.Equal(docC.DocumentID, docB.DocumentID)
	assert.Substring("2-", docC.DocumentRevision)
	assert.Equal(docC.FieldA, docB.FieldA)
	assert.Equal(docC.FieldB, docB.FieldB)

	// Read the first revision.
	resp = cdb.ReadDocument(id, couchdb.Revision(revision))
	assert.True(resp.IsOK())
	docD := MyDocument{}
	err = resp.ResultValue(&docD)
	assert.Nil(err)
	assert.Equal(docD.DocumentRevision, revision)
	assert.Equal(docD.FieldB, docA.FieldB)
}

// TestDeleteDocument tests deleting a document.
func TestDeleteDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb := prepareDatabase(assert)

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		FieldA:     "foo",
		FieldB:     12345,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document, we need it including the revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.ResultValue(&docB)
	assert.Nil(err)

	// Delete the test document.
	resp = cdb.DeleteDocument(docB)
	assert.True(resp.IsOK())

	// Try to read deleted document.
	resp = cdb.ReadDocument(id)
	assert.False(resp.IsOK())
	assert.ErrorMatch(resp.Error(), ".* 404,.*")
}

//--------------------
// HELPERS
//--------------------

// MyDocument is used for the tests.
type MyDocument struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	FieldA string
	FieldB int
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
