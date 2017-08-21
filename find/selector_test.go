// Tideland Go CouchDB Client - Find - Unit Tests - Selector
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find_test

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"testing"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocouch/find"
)

//--------------------
// CONSTANTS
//--------------------

//--------------------
// TESTS
//--------------------

// TestSelector tests creating selectors.
func TestSelector(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	selectorA := find.SelectAnd(func(as find.Selector) {
		as.Equal("foo", 4711)
		as.In("years", 1965, 1989, 2017)
		as.Append(find.SelectOr(func(os find.Selector) {
			os.All("genre", "comedy", "short")
			os.NotEqual("age", 18)
		}))
		as.GreaterThan("count", 4711).Not()
	})

	b, err := json.Marshal(selectorA)
	assert.Nil(err)
	assert.Equal(string(b), `{"$and":[{"foo":{"$eq":4711}},{"years":{"$in":[1965,1989,2017]}},`+
		`{"$or":[{"genre":{"$all":["comedy","short"]}},{"age":{"$ne":18}}]},{"$not":{"count":{"$gt":4711}}}]}`)

	selectorB := find.SelectOr(nil)
	selectorB.Equal("left", true)
	selectorB.GreaterThan("height", 100).Not()

	b, err = json.Marshal(selectorB)
	assert.Nil(err)
	assert.Equal(string(b), `{"$or":[{"left":{"$eq":true}},{"$not":{"height":{"$gt":100}}}]}`)
}

// EOF
