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
	selector := find.NewAndSelector().
		Equal("foo", 4711).
		Equal("bar", "42")

	assert.NotNil(selector)
}

// EOF
