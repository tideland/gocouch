// Tideland Go CouchDB Client - Startup - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package startup_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/version"

	"github.com/tideland/gocouch/couchdb"
	"github.com/tideland/gocouch/startup"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TemporaryDBCfg = "{etc {hostname localhost}{port 5984}{database tgocouch-testing-temporary-startup}}"
)

//--------------------
// TESTS
//--------------------

// TestNoSteps tests creating the database with no steps.
func TestNoSteps(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := etc.ReadString(TemporaryDBCfg)
	assert.Nil(err)

	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	defer func() { cdb.DeleteDatabase() }()

	err = startup.Run(cdb)
	assert.Nil(err)

	ok, err := cdb.HasDatabase()
	assert.Nil(err)
	assert.True(ok)

	resp := cdb.ReadDocument(startup.DatabaseVersionID)
	assert.True(resp.IsOK())

	dv := startup.DatabaseVersion{}
	err = resp.Document(&dv)
	assert.Nil(err)
	assert.Equal(dv.Version, "0.0.0")
}

// TestSomeSteps tests creating the database with some steps.
func TestSomeSteps(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := etc.ReadString(TemporaryDBCfg)
	assert.Nil(err)

	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	defer func() { cdb.DeleteDatabase() }()

	err = startup.Run(cdb, StepA, StepB)
	assert.Nil(err)

	resp := cdb.ReadDocument(startup.DatabaseVersionID)
	assert.True(resp.IsOK())

	dv := startup.DatabaseVersion{}
	err = resp.Document(&dv)
	assert.Nil(err)
	assert.Equal(dv.Version, "0.2.0")

	ids, err := cdb.AllDocuments()
	assert.Nil(err)
	assert.Length(ids, 3)
}

// TestMultipleStartups tests calling startup multiple times.
func TestMultipleStartups(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)

	cfg, err := etc.ReadString(TemporaryDBCfg)
	assert.Nil(err)

	cdb, err := couchdb.Open(cfg)
	assert.Nil(err)
	defer func() { cdb.DeleteDatabase() }()

	err = startup.Run(cdb, StepA)
	assert.Nil(err)

	ids, err := cdb.AllDocuments()
	assert.Nil(err)
	assert.Length(ids, 2)

	resp := cdb.ReadDocument(startup.DatabaseVersionID)
	assert.True(resp.IsOK())

	dv := startup.DatabaseVersion{}
	err = resp.Document(&dv)
	assert.Nil(err)
	assert.Equal(dv.Version, "0.1.0")

	err = startup.Run(cdb, StepA, StepB, StepC)
	assert.Nil(err)

	resp = cdb.ReadDocument(startup.DatabaseVersionID)
	assert.True(resp.IsOK())

	dv = startup.DatabaseVersion{}
	err = resp.Document(&dv)
	assert.Nil(err)

	assert.Equal(dv.Version, "0.3.0")
	ids, err = cdb.AllDocuments()
	assert.Nil(err)
	assert.Length(ids, 4)
}

//--------------------
// HELPERS
//--------------------

// MyDocument is used for the tests.
type MyDocument struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`
	Name             string `json:"name"`
	Age              int    `json:"age"`
}

func StepA() (version.Version, startup.StepAction) {
	v := version.New(0, 1, 0)
	return v, func(cdb couchdb.CouchDB) error {
		md := MyDocument{
			DocumentID: "my-document-a",
			Name:       "Joe Black",
			Age:        25,
		}
		resp := cdb.CreateDocument(&md)
		if !resp.IsOK() {
			return resp.Error()
		}
		return nil
	}
}

func StepB() (version.Version, startup.StepAction) {
	v := version.New(0, 2, 0)
	return v, func(cdb couchdb.CouchDB) error {
		md := MyDocument{
			DocumentID: "my-document-b",
			Name:       "John Doe",
			Age:        51,
		}
		resp := cdb.CreateDocument(&md)
		if !resp.IsOK() {
			return resp.Error()
		}
		return nil
	}
}

func StepC() (version.Version, startup.StepAction) {
	v := version.New(0, 3, 0)
	return v, func(cdb couchdb.CouchDB) error {
		md := MyDocument{
			DocumentID: "my-document-c",
			Name:       "Donald Duck",
			Age:        85,
		}
		resp := cdb.CreateDocument(&md)
		if !resp.IsOK() {
			return resp.Error()
		}
		return nil
	}
}

// EOF
