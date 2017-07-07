// Tideland Go CouchDB Client - Views - Unit Tests
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package views_test

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

	"github.com/tideland/gocouch/couchdb"
	"github.com/tideland/gocouch/views"
)

//--------------------
// CONSTANTS
//--------------------

const (
	Cfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}{debug-logging true}}"
)

//--------------------
// TESTS
//--------------------

// TestCallingView tests calling a view.
func TestCallingView(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("view-documents", assert)
	defer cleanup()

	// Create design document.
	design, err := cdb.Design("testing")
	assert.Nil(err)
	design.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc);  } }", "")
	design.SetView("age", "function(doc){ emit(doc.age, doc.name); }", "")
	resp := design.Write()
	assert.True(resp.IsOK())

	// Call the view for the first time.
	vrs := views.View(cdb, "testing", "index-a")
	assert.True(vrs.IsOK())
	trOld := vrs.TotalRows()
	assert.True(trOld > 0)

	// Add a matching document and view again.
	docA := MyDocument{
		DocumentID: "black-jack-4711",
		Name:       "Jack Black",
	}
	resp = cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	vrs = views.View(cdb, "testing", "index-a")
	assert.True(vrs.IsOK())
	trNew := vrs.TotalRows()
	assert.Equal(trNew, trOld+1)
	err = vrs.RowsDo(func(id string, key, value, document couchdb.Unmarshable) error {
		valueA := MyDocument{}
		err := value.Unmarshal(&valueA)
		assert.Nil(err)
		assert.True(strings.Contains(valueA.DocumentID, "a"))
		return err
	})
	assert.Nil(err)

	// Add a non-matching document and view again.
	docB := MyDocument{
		DocumentID: "doe-john-999",
		Name:       "John Doe",
	}
	resp = cdb.CreateDocument(docB)
	assert.True(resp.IsOK())
	vrs = views.View(cdb, "testing", "index-a")
	assert.True(vrs.IsOK())
	trFinal := vrs.TotalRows()
	assert.Equal(trFinal, trNew)

	// Call age view with a key.
	vrs = views.View(cdb, "testing", "age", views.OneKey(51))
	assert.True(vrs.IsOK())
	assert.True(vrs.TotalRows() > vrs.ReturnedRows())
	err = vrs.RowsDo(func(id string, key, value, document couchdb.Unmarshable) error {
		var age int
		var name string
		err := key.Unmarshal(&age)
		assert.Nil(err)
		assert.Equal(age, 51)
		err = value.Unmarshal(&name)
		assert.Nil(err)
		return err
	})
	assert.Nil(err)

	// Call age view with the oldest 5 peaple below 50.
	vrs = views.View(cdb, "testing", "age", views.StartKey(50), views.Descending(), views.Limit(5))
	assert.True(vrs.IsOK())
	assert.True(vrs.ReturnedRows() <= 5)
	err = vrs.RowsDo(func(id string, key, value, document couchdb.Unmarshable) error {
		var age int
		var name string
		err := key.Unmarshal(&age)
		assert.Nil(err)
		assert.True(age <= 50)
		err = value.Unmarshal(&name)
		assert.Nil(err)
		assert.Logf("Tester %s has age %d", name, age)
		return err
	})
	assert.Nil(err)

	// Call age view with multiple keys (even multiple times).
	vrs = views.View(cdb, "testing", "age", views.Keys(50, 51, 52), views.Keys(53, 54))
	assert.True(vrs.IsOK())
	err = vrs.RowsDo(func(id string, key, value, document couchdb.Unmarshable) error {
		var age int
		var name string
		err := key.Unmarshal(&age)
		assert.Nil(err)
		assert.True(age >= 50)
		assert.True(age <= 54)
		err = value.Unmarshal(&name)
		assert.Nil(err)
		assert.Logf("Tester %s has age %d", name, age)
		return err
	})
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
