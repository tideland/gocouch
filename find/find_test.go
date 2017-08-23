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
	"time"

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
	cdb, cleanup := prepareFilledDatabase("find-simple", 1000, assert)
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
	frs := find.Find(cdb, selector, find.Fields("name", "age", "active"))
	assert.NotNil(frs)
	assert.True(frs.IsOK())

	err := frs.Do(func(document couchdb.Unmarshable) error {
		fields := struct {
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

// TestLimitedFind tests retrieving a larger number but set the limit.
func TestLimitedFind(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("find-limited", 1000, assert)
	defer cleanup()

	// Try to find some documents a simple way.
	selector := find.SelectAnd(func(as find.Selector) {
		as.Equal("active", true)
	})
	frs := find.Find(cdb, selector, find.Fields("name", "active"), find.Limit(5))
	assert.NotNil(frs)
	assert.True(frs.IsOK())
	assert.Length(frs, 5)

	frs = find.Find(cdb, selector, find.Fields("name", "active"), find.Limit(50))
	assert.NotNil(frs)
	assert.True(frs.IsOK())
	assert.Length(frs, 50)
}

// TestFindExists tests calling find with an exists selector.
func TestFindExists(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("find-exists", 1000, assert)
	defer cleanup()

	// Try to find some documents having an existing "last_active".
	selector := find.SelectAnd(func(as find.Selector) {
		as.Exists("last_active")
		as.LowerEqualThan("age", 25)
	})
	frs := find.Find(cdb, selector, find.Fields("name", "age", "active", "last_active"))
	assert.NotNil(frs)
	assert.True(frs.IsOK())

	err := frs.Do(func(document couchdb.Unmarshable) error {
		fields := struct {
			Name       string `json:"name"`
			Age        int    `json:"age"`
			Active     bool   `json:"active"`
			LastActive int64  `json:"last_active"`
		}{}
		if err := document.Unmarshal(&fields); err != nil {
			return err
		}
		assert.True(fields.Age <= 25 && fields.LastActive > 0 && fields.Active)
		lastActive := time.Unix(fields.LastActive, 0).Format(time.RFC1123)
		assert.Logf("person with name %s (age %d) has been last active at %v", fields.Name, fields.Age, lastActive)
		return nil
	})
	assert.Nil(err)

	// Now look for existing "last_active" but "active" is false. So
	// no results.
	selector = find.SelectAnd(func(as find.Selector) {
		as.Exists("last_active")
		as.NotEqual("active", true)
	})
	frs = find.Find(cdb, selector, find.Fields("name", "age", "active", "last_active"))
	assert.NotNil(frs)
	assert.True(frs.IsOK())
	assert.Equal(frs.Len(), 0)
}

// TestSingleMatch tests using only one selector, here a regular expression.
func TestSingleMatch(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	cdb, cleanup := prepareFilledDatabase("find-match", 1000, assert)
	defer cleanup()

	// Try to find some documents having an existing "last_active".
	selector := find.SelectAnd(func(as find.Selector) {
		as.RegExp("name", ".*Adam.*")
	})
	frs := find.Find(cdb, selector, find.Fields("name", "age", "active"))
	assert.NotNil(frs)
	assert.Nil(frs.Error())
	assert.True(frs.IsOK())

	err := frs.Do(func(document couchdb.Unmarshable) error {
		fields := struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Active bool   `json:"active"`
		}{}
		if err := document.Unmarshal(&fields); err != nil {
			return err
		}
		assert.Match(fields.Name, ".*Adam.*")
		return nil
	})
	assert.Nil(err)
}

//--------------------
// HELPERS
//--------------------

// Note is used for the tests.
type Note struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Person is used for the tests.
type Person struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Name       string `json:"name"`
	Age        int    `json:"age"`
	Active     bool   `json:"active"`
	LastActive int64  `json:"last_active,omitempty"`
	Notes      []Note `json:"notes"`
}

// prepareFilledDatabase opens the database, deletes a possible test
// database, creates it newly and adds some data.
func prepareFilledDatabase(database string, count int, assert audit.Assertion) (couchdb.CouchDB, func()) {
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
	runs := count / 1000
	for outer := 0; outer < runs; outer++ {
		assert.Logf("filling database run %d of %d", outer+1, runs)
		docs := []interface{}{}
		for inner := 0; inner < 1000; inner++ {
			first, middle, last := gen.Name()
			person := Person{
				DocumentID: identifier.Identifier(last, first, outer, inner),
				Name:       first + " " + middle + " " + last,
				Age:        gen.Int(18, 65),
				Active:     gen.FlipCoin(75),
			}
			if person.Active {
				person.LastActive = gen.Time(time.UTC, time.Now().Add(-24*time.Hour), 24*time.Hour).Unix()
			}
			for j := 0; j < gen.Int(3, 9); j++ {
				note := Note{
					Title: gen.Sentence(),
					Text:  gen.Paragraph(),
				}
				person.Notes = append(person.Notes, note)
			}
			docs = append(docs, person)
		}
		results, err := cdb.BulkWriteDocuments(docs)
		assert.Nil(err)
		for _, result := range results {
			assert.True(result.OK)
		}
	}

	return cdb, func() { cdb.DeleteDatabase() }
}

// EOF
