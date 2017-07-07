// Tideland Go CouchDB Client - Changes - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package changes_test

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/identifier"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocouch/changes"
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TemplateDBcfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}{debug-logging true}}"
)

//--------------------
// TESTS
//--------------------

// TestChanges tests retrieving changes.
func TestChanges(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	count := 1000
	cdb, gen, cleanup := prepareFilledDatabase(assert, "changes", count)
	defer cleanup()

	// Simple changes access.
	crs := changes.Changes(cdb)
	assert.True(crs.IsOK())
	assert.Equal(crs.ResultsLen(), count)

	crs.ResultsDo(func(id, sequence string, deleted bool, revisions []string, document couchdb.Unmarshable) error {
		assert.Length(revisions, 1)
		return nil
	})

	lseq := crs.LastSequence()
	crs = changes.Changes(cdb, changes.Since(lseq))
	assert.True(crs.IsOK())
	assert.Equal(crs.ResultsLen(), 0)

	crs = changes.Changes(cdb, changes.Since(changes.SinceNow))
	assert.True(crs.IsOK())
	assert.Equal(crs.ResultsLen(), 0)

	// Add some more documents and check changes.
	docs := generateDocuments(gen, count)
	results, err := cdb.BulkWriteDocuments(docs)
	assert.Nil(err)
	for _, result := range results {
		assert.True(result.OK)
	}

	crs = changes.Changes(cdb)
	assert.True(crs.IsOK())
	assert.Equal(crs.ResultsLen(), 2*count)

	crs = changes.Changes(cdb, changes.Since(lseq))
	assert.True(crs.IsOK())
	assert.Equal(crs.ResultsLen(), count)
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

// prepareFilledDatabase opens the database, deletes a possible test
// database, creates it newly and adds some data.
func prepareFilledDatabase(assert audit.Assertion, database string, count int) (couchdb.CouchDB, *audit.Generator, func()) {
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
	docs := generateDocuments(gen, count)
	results, err := cdb.BulkWriteDocuments(docs)
	assert.Nil(err)
	for _, result := range results {
		assert.True(result.OK)
	}
	return cdb, gen, func() { cdb.DeleteDatabase() }
}

// generateDocuments creates a number of documents.
func generateDocuments(gen *audit.Generator, count int) []interface{} {
	docs := []interface{}{}
	for i := 0; i < count; i++ {
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
	return docs
}

// EOF
