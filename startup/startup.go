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

// ShallStep returns true if the new version is newer than
// the current version.
func ShallStep(cv, nv version.Version) bool {
	precedence, _ := nv.Compare(cv)
	return precedence == version.Newer
}

// Step defines one step in starting up a CouchDB. It receives
// the current database version and has to return a new one if
// it changes the database. So it initially should compare the
// versions:
//
//    nv := version.New(1, 2, 3)
//    if !startup.ShallStep(v, nv) {
//        return nil, nil
//    }
//    ...
//    return nv, nil
type Step func(cdb couchdb.CouchDB, v version.Version) (version.Version, error)

// run performs one step.
func (s Step) run(cdb couchdb.CouchDB) error {
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
		return err
	}
	// Now perform the step.
	nv, err := s(cdb, cv)
	if err != nil {
		return err
	}
	// Update version document only if needed.
	if nv == nil {
		return nil
	}
	dv.Version = nv.String()
	resp = cdb.UpdateDocument(&dv)
	if !resp.IsOK() {
		return resp.Error()
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
	for _, step := range steps {
		if err := step.run(cdb); err != nil {
			return err
		}
	}
	return nil
}

// EOF
