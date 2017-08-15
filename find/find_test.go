// Tideland Go CouchDB Client - Find - Unit Tests
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

const (
	TemplateDBcfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-<<DATABASE>>}{debug-logging true}}"
)

//--------------------
// TESTS
//--------------------

// TestSelector tests creating selectors.
func TestSelector(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	// And selector.
	andSel := find.NewAndSelector().
		Equal("foo", 4711).
		Equal("bar", "42")

	assert.NotNil(andSel)

	b, err := json.Marshal(andSel)
	assert.Nil(err)
	assert.Logf("SELECTOR %s", string(b))

	// Or selector.
	orSel := find.NewOrSelector().
		Equal("yadda", true)

	assert.NotNil(orSel)

	b, err = json.Marshal(orSel)
	assert.Nil(err)
	assert.Logf("SELECTOR %s", string(b))

	// Combine these two.
	combSel := find.NewAndSelector().
		SubSelectors(andSel, orSel)

	b, err = json.Marshal(combSel)
	assert.Nil(err)
	assert.Logf("SELECTOR %s", string(b))

}

// EOF
