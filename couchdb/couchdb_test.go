// Tideland Go CouchDB Client - CouchDB - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb_test

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/identifier"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

const (
	EmptyCfg      = "{etc}"
	TemplateDBcfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}{debug-logging true}}"
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

	cfg, err := couchdb.Configure("localhost", 5984, "tgocouch-testing-temporary")
	assert.Nil(err)

	// This time also use OpenPath() to check its behavior.
	cdb, err := couchdb.OpenPath(cfg, "couchdb")
	assert.Nil(err)
	_, err = cdb.AllDatabases()
	assert.Nil(err)

	cfg, err = etc.ReadString(EmptyCfg)
	assert.Nil(err)

	cdb, err = couchdb.OpenPath(cfg, "")
	assert.Nil(err)
	_, err = cdb.AllDatabases()
	assert.Nil(err)
}

// TestCreateDeleteDatabase tests the creation and deletion
// of a database.
func TestCreateDeleteDatabase(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := couchdb.Configure("localhost", 5984, "tgocouch-testing-temporary")
	assert.Nil(err)

	// Open and check existence.
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	has, err := cdb.HasDatabase()
	assert.Nil(err)
	assert.False(has)

	// Create and check existence,
	resp := cdb.CreateDatabase()
	assert.True(resp.IsOK())
	has, err = cdb.HasDatabase()
	assert.Nil(err)
	assert.True(has)

	// Delete and check existence.
	resp = cdb.DeleteDatabase()
	assert.True(resp.IsOK())
	has, err = cdb.HasDatabase()
	assert.Nil(err)
	assert.False(has)
}

// TestCreateDesignDocument tests creating new design documents.
func TestCreateDesignDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("create-design", assert)
	defer cleanup()

	// Create design document and check if it has been created.
	allDesignA, err := cdb.AllDesigns()
	assert.Nil(err)

	design, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(design.ID(), "testing-a")
	design.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := design.Write()
	assert.True(resp.IsOK())

	design, err = cdb.Design("testing-b")
	assert.Nil(err)
	assert.Equal(design.ID(), "testing-b")
	design.SetView("index-b", "function(doc){ if (doc._id.indexOf('b') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp = design.Write()
	assert.True(resp.IsOK())

	allDesignB, err := cdb.AllDesigns()
	assert.Nil(err)
	assert.Equal(len(allDesignB), len(allDesignA)+2)
}

// TestReadDesignDocument tests reading design documents.
func TestReadDesignDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("read-design", assert)
	defer cleanup()

	// Create design document and read it again.
	designA, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(designA.ID(), "testing-a")
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	designB, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(designB.ID(), "testing-a")
}

// TestUpdateDesignDocument tests updating design documents.
func TestUpdateDesignDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("update-design", assert)
	defer cleanup()

	// Create design document and read it again.
	designA, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(designA.ID(), "testing-a")
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	designB, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(designB.ID(), "testing-a")

	// Now update it and read it again.
	designB.SetView("index-b", "function(doc){ if (doc._id.indexOf('b') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp = designB.Write()
	assert.True(resp.IsOK())

	designC, err := cdb.Design("testing-a")
	assert.Nil(err)
	assert.Equal(designC.ID(), "testing-a")
	_, _, ok := designC.View("index-a")
	assert.True(ok)
	_, _, ok = designC.View("index-b")
	assert.True(ok)
}

// TestDeleteDesignDocument tests deleting design documents.
func TestDeleteDesignDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("delete-design", assert)
	defer cleanup()

	// Create design document and check if it has been created.
	allDesignA, err := cdb.AllDesigns()
	assert.Nil(err)

	designA, err := cdb.Design("testing")
	assert.Nil(err)
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	allDesignB, err := cdb.AllDesigns()
	assert.Nil(err)
	assert.Equal(len(allDesignB), len(allDesignA)+1)

	// Read it and delete it.
	designB, err := cdb.Design("testing")
	assert.Nil(err)
	resp = designB.Delete()
	assert.True(resp.IsOK())

	allDesignC, err := cdb.AllDesigns()
	assert.Nil(err)
	assert.Equal(len(allDesignC), len(allDesignA))
}

// TestCreateDocument tests creating new documents.
func TestCreateDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("create-document", assert)
	defer cleanup()

	// Create document without ID.
	docA := MyDocument{
		Name: "foo",
		Age:  50,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Match(id, "[0-9a-f]{32}")

	// Create document with ID.
	docB := MyDocument{
		DocumentID: "bar-12345",
		Name:       "bar",
		Age:        25,
		Active:     true,
	}
	resp = cdb.CreateDocument(docB)
	assert.True(resp.IsOK())
	id = resp.ID()
	assert.Equal(id, "bar-12345")
}

// TestReadDocument tests reading a document.
func TestReadDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("read-document", assert)
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        18,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)
	assert.Equal(docB.DocumentID, docA.DocumentID)
	assert.Equal(docB.Name, docA.Name)
	assert.Equal(docB.Age, docA.Age)

	// Try to read non-existent document.
	resp = cdb.ReadDocument("i-do-not-exist")
	assert.False(resp.IsOK())
	assert.ErrorMatch(resp.Error(), ".* 404,.*")
}

// TestUpdateDocument tests updating documents.
func TestUpdateDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("update-document", assert)
	defer cleanup()

	// Create first revision.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        22,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	revision := resp.Revision()
	assert.Equal(id, "foo-12345")

	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)

	// Update the document.
	docB.Age = 23

	resp = cdb.UpdateDocument(docB)
	assert.True(resp.IsOK())

	// Read the updated revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docC := MyDocument{}
	err = resp.Document(&docC)
	assert.Nil(err)
	assert.Equal(docC.DocumentID, docB.DocumentID)
	assert.Substring("2-", docC.DocumentRevision)
	assert.Equal(docC.Name, docB.Name)
	assert.Equal(docC.Age, docB.Age)

	// Read the first revision.
	resp = cdb.ReadDocument(id, couchdb.Revision(revision))
	assert.True(resp.IsOK())
	assert.Equal(resp.Revision(), revision)

	// Try to update a non-existent document.
	docD := MyDocument{
		DocumentID: "i-do-not-exist",
		Name:       "none",
		Age:        999,
	}
	resp = cdb.UpdateDocument(docD)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
}

// TestDeleteDocument tests deleting a document.
func TestDeleteDocument(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("delete-document", assert)
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        33,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document, we need it including the revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)

	// Delete the test document.
	resp = cdb.DeleteDocument(docB)
	assert.True(resp.IsOK())

	// Try to read deleted document.
	resp = cdb.ReadDocument(id)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)

	// Try to delete it a second time.
	resp = cdb.DeleteDocument(docB)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
}

// TestDeleteDocumentByID tests deleting a document by identifier.
func TestDeleteDocumentByID(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareDatabase("delete-document-by-id", assert)
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        33,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	revision := resp.Revision()
	assert.Equal(id, "foo-12345")

	// Delete the test document by ID.
	resp = cdb.DeleteDocumentByID(id, revision)
	assert.True(resp.IsOK())

	// Try to read deleted document.
	resp = cdb.ReadDocument(id)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)

	// Try to delete it a second time.
	resp = cdb.DeleteDocumentByID(id, revision)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
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
	logger.SetLevel(logger.LevelDebug)
	cfgstr := strings.Replace(TemplateDBcfg, "<<DATABASE>>", database, 1)
	cfg, err := etc.ReadString(cfgstr)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	rs := cdb.DeleteDatabase()
	rs = cdb.CreateDatabase()
	assert.True(rs.IsOK())
	return cdb, func() { cdb.DeleteDatabase() }
}

// prepareFilledDatabase opens the database, deletes a possible test
// database, creates it newly and adds some data.
func prepareFilledDatabase(database string, assert audit.Assertion) (couchdb.CouchDB, func()) {
	logger.SetLevel(logger.LevelDebug)
	cfgstr := strings.Replace(TemplateDBcfg, "<<DATABASE>>", database, 1)
	cfg, err := etc.ReadString(cfgstr)
	assert.Nil(err)
	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	rs := cdb.DeleteDatabase()
	rs = cdb.CreateDatabase()
	assert.True(rs.IsOK())

	gen := audit.NewGenerator(audit.FixedRand())
	docs := []interface{}{}
	for i := 0; i < 1000; i++ {
		first, middle, last := gen.Name()
		doc := MyDocument{
			DocumentID:  identifier.Identifier(last, first, i),
			Name:        first + " " + middle + " " + last,
			Age:         gen.Int(18, 65),
			Active:      gen.FlipCoin(75),
			Description: gen.Sentence(),
		}
		docs = append(docs, doc)
	}
	results, err := cdb.BulkWriteDocuments(docs)
	assert.Nil(err)
	for _, result := range results {
		assert.True(result.OK)
	}

	return cdb, func() { cdb.DeleteDatabase() }
}

// EOF
