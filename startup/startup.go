// Tideland Go CouchDB Client - Startup
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package startup

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/version"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// DOCUMENTS
//--------------------

// DatabaseVersionID is used for the database version document.
const DatabaseVersionID = "database-version"

// DatabaseVersion stores the current database version with
// the document ID "database-version".
type DatabaseVersion struct {
	ID       string `json:"_id"`
	Revision string `json:"_rev,omitempty"`
	Version  string `json:"version"`
}

//--------------------
// STEP
//--------------------

// StepAction is the concrete action of a step.
type StepAction func(cdb couchdb.CouchDB) error

// Step returns the version after a startup step and the action
// that shall be performed on the database. The returned action
// will only be performed, if the current if the new version is
// than the current version.
type Step func() (version.Version, StepAction)

// run performs one step.
func (step Step) run(cdb couchdb.CouchDB) error {
	// Retrieve current database version.
	resp := cdb.ReadDocument(DatabaseVersionID)
	if !resp.IsOK() {
		return resp.Error()
	}
	dv := DatabaseVersion{}
	err := resp.ResultValue(&dv)
	if err != nil {
		return err
	}
	cv, err := version.Parse(dv.Version)
	if err != nil {
		return errors.Annotate(err, ErrIllegalVersion, errorMessages)
	}
	// Get new version of the step and action.
	nv, action := step()
	// Check the new version.
	precedence, _ := nv.Compare(cv)
	if precedence != version.Newer {
		return nil
	}
	// Now perform the step action and update the
	// version document.
	err = action(cdb)
	if err != nil {
		return errors.Annotate(err, ErrStartupActionFailed, errorMessages, nv)
	}
	dv.Version = nv.String()
	resp = cdb.UpdateDocument(&dv)
	if !resp.IsOK() {
		return resp.Error()
	}
	return nil
}

// Steps is just an ordered number of steps.
type Steps []Step

// run performs the steps.
func (steps Steps) run(cdb couchdb.CouchDB) error {
	for _, step := range steps {
		if err := step.run(cdb); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// RUN
//--------------------

// Run checks and creates the database if needed and performs
// the individual steps.
func Run(cdb couchdb.CouchDB, steps ...Step) error {
	// Check database.
	ok, err := cdb.HasDatabase()
	if err != nil {
		return err
	}
	// Create and initialize it.
	if !ok {
		resp := cdb.CreateDatabase()
		if !resp.IsOK() {
			return resp.Error()
		}
		dv := DatabaseVersion{
			ID:      DatabaseVersionID,
			Version: version.New(0, 0, 0).String(),
		}
		resp = cdb.CreateDocument(&dv)
		if !resp.IsOK() {
			return resp.Error()
		}
	}
	// Run the steps.
	return Steps(steps).run(cdb)
}

// EOF
