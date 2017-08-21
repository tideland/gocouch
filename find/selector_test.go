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
	assert.Logf("SELECTOR A) %s", string(b))
}

// EOF
