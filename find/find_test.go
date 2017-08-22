// Tideland Go CouchDB Client - Find - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find_test

import (
	"strings"
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/identifier"
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocouch/couchdb"
	"github.com/tideland/gocouch/find"
)

//--------------------
// IMPORTS
//--------------------

//--------------------
// CONSTANTS
//--------------------

const (
	Cfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}{debug-logging true}}"
)

//--------------------
// TESTS
//--------------------

// TestSimpleFind tests calling find with a simple selector.
func TestSimpleFind(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("find-documents", assert)
	defer cleanup()

	// Try to find some documents a simple way.
	selector := find.SelectOr(func(os find.Selector) {
		os.Append(find.SelectAnd(func(as find.Selector) {
			as.LowerThan("age", 30)
			as.Equal("active", false)
		}))
		os.Append(find.SelectAnd(func(as find.Selector) {
			as.GreaterThan("age", 60)
			as.Equal("active", true)
		}))
	})
	frs := find.Find(cdb, find.Request{
		Selector: selector,
		Fields:   []string{"_id", "name", "age", "active"},
	})
	assert.NotNil(frs)
	assert.True(frs.IsOK())

	err := frs.Do(func(document couchdb.Unmarshable) error {
		fields := struct {
			ID     string `json:"_id"`
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Active bool   `json:"active"`
		}{}
		if err := document.Unmarshal(&fields); err != nil {
			return err
		}
		assert.True((fields.Age < 30 && !fields.Active) || (fields.Age > 60 && fields.Active))
		assert.Logf("person with name %s has age %d and activity status is %v", fields.Name, fields.Age, fields.Active)
		return nil
	})
	assert.Nil(err)
}

//--------------------
// HELPERS
//--------------------

// Description is used for the tests.
type Description struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Person is used for the tests.
type Person struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Name         string        `json:"name"`
	Age          int           `json:"age"`
	Active       bool          `json:"active"`
	Descriptions []Description `json:"descriptions"`
}

// prepareFilledDatabase opens the database, deletes a possible test
// database, creates it newly and adds some data.
func prepareFilledDatabase(database string, assert audit.Assertion) (couchdb.CouchDB, func()) {
	logger.SetLevel(logger.LevelDebug)
	cfgstr := strings.Replace(Cfg, "<<DATABASE>>", database, 1)
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
		doc := Person{
			DocumentID: identifier.Identifier(last, first, i),
			Name:       first + " " + middle + " " + last,
			Age:        gen.Int(18, 65),
			Active:     gen.FlipCoin(75),
		}
		for j := 0; j < gen.Int(3, 9); j++ {
			description := Description{
				Title: gen.Sentence(),
				Text:  gen.Paragraph(),
			}
			doc.Descriptions = append(doc.Descriptions, description)
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
