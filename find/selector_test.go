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

// TestCriteria tests creating and composing query criteria.
func TestCriteria(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	// Nested criteria.
	criterion := find.And(
		find.Equal("foo", 4711),
		find.In("year", 1965, 1989, 2017),
		find.Or(
			find.All("genre", "comedy", "short"),
			find.NotEqual("age", 18),
		),
		find.GreaterThan("count", 4711).Not(),
	)
	b, err := json.Marshal(criterion)
	assert.Nil(err)
	assert.Equal(string(b), `{"$and":[{"foo":{"$eq":4711}},{"year":{"$in":[1965,1989,2017]}},`+
		`{"$or":[{"genre":{"$all":["comedy","short"]}},{"age":{"$ne":18}}]},{"$not":{"count":{"$gt":4711}}}]}`)

	// Only one sub-criterion, but must render.
	criterion = find.Or(find.Equal("foo", "bar"))
	b, err = json.Marshal(criterion)
	assert.Nil(err)
	assert.Equal(string(b), `{"$or":[{"foo":{"$eq":"bar"}}]}`)
}

// TestSelector tests creating selectors.
func TestSelector(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	assert.True(true)
}

// EOF
