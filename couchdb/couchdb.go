// Tideland Go CouchDB Client - CouchDB
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tideland/golib/errors"
	"github.com/tideland/golib/etc"
	"github.com/tideland/golib/identifier"
)

//--------------------
// COUCHDB
//--------------------

// CouchDB provides the access to a database.
type CouchDB interface {
	// AllDatabases returns a list of all database IDs
	// of the connected server.
	AllDatabases() ([]string, error)

	// CreateDatabase creates the configured database.
	CreateDatabase() Response

	// DeleteDatabase removes the configured database.
	DeleteDatabase() Response

	// AllDesignDocuments returns the lsit of all design
	// document IDs of the configured database.
	AllDesignDocuments() ([]string, error)

	// CreateDesignDocument creates a new design document.
	CreateDesignDocument(doc *DesignDocument) Response

	// AllDocuments returns a list of all document IDs
	// of the configured database.
	AllDocuments() ([]string, error)

	// CreateDocument creates a new document.
	CreateDocument(doc interface{}, rps ...Parameter) Response

	// ReadDocument reads an existing document.
	ReadDocument(id string, rps ...Parameter) Response

	// UpdateDocument update an existing document.
	UpdateDocument(doc interface{}, rps ...Parameter) Response

	// DeleteDocument deletes an existing document.
	DeleteDocument(doc interface{}, rps ...Parameter) Response

	// ViewDocuments reads the output of a view.
	ViewDocuments(design, view string, rps ...Parameter) Response
}

// couchdb implements CouchDB.
type couchdb struct {
	host     string
	database string
}

// Open returns a configured connection to a CouchDB server.
func Open(cfg etc.Etc) (CouchDB, error) {
	if cfg == nil {
		return nil, errors.New(ErrNoConfiguration, errorMessages)
	}
	host := fmt.Sprintf("%s:%d",
		cfg.ValueAsString("hostname", "localhost"),
		cfg.ValueAsInt("port", 5984),
	)
	db := &couchdb{
		host:     host,
		database: cfg.ValueAsString("database", "default"),
	}
	return db, nil
}

// AllDatabases implements connection.
func (db *couchdb) AllDatabases() ([]string, error) {
	req := newRequest(db, "/_all_dbs", nil)
	resp := req.get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	ids := []string{}
	err := resp.ResultValue(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// CreateDatabase implements the CouchDB interface.
func (db *couchdb) CreateDatabase() Response {
	req := newRequest(db, db.databasePath(), nil)
	return req.put()
}

// DeleteDatabase implements the CouchDB interface.
func (db *couchdb) DeleteDatabase() Response {
	req := newRequest(db, db.databasePath(), nil)
	return req.delete()
}

// AllDesignDocuments implements the CouchDB interface.
func (db *couchdb) AllDesignDocuments() ([]string, error) {
	req := newRequest(db, db.databasePath("_all_docs"), nil)
	resp := req.setParameters(StartEndKey("_design/", "_design0")).get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	vr := ViewResult{}
	err := resp.ResultValue(&vr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range vr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// CreateDesignDocument implements the CouchDB interface.
func (db *couchdb) CreateDesignDocument(doc *DesignDocument) Response {
	req := newRequest(db, db.databasePath(doc.ID), doc)
	return req.put()
}

// AllDocuments implements the CouchDB interface.
func (db *couchdb) AllDocuments() ([]string, error) {
	req := newRequest(db, db.databasePath("_all_docs"), nil)
	resp := req.get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	vr := ViewResult{}
	err := resp.ResultValue(&vr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range vr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// CreateDocument implements the CouchDB interface.
func (db *couchdb) CreateDocument(doc interface{}, rps ...Parameter) Response {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	req := newRequest(db, db.databasePath(id), doc)
	return req.setParameters(rps...).put()
}

// ReadDocument implements the CouchDB interface.
func (db *couchdb) ReadDocument(id string, rps ...Parameter) Response {
	req := newRequest(db, db.databasePath(id), nil)
	return req.setParameters(rps...).get()
}

// UpdateDocument implements the CouchDB interface.
func (db *couchdb) UpdateDocument(doc interface{}, rps ...Parameter) Response {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	if id == "" {
		return newResponse(nil, errors.New(ErrNoIdentifier, errorMessages))
	}
	req := newRequest(db, db.databasePath(id), doc)
	return req.setParameters(rps...).put()
}

// DeleteDocument implements the CouchDB interface.
func (db *couchdb) DeleteDocument(doc interface{}, rps ...Parameter) Response {
	id, rev, err := db.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	rps = append(rps, Revision(rev))
	req := newRequest(db, db.databasePath(id), nil)
	return req.setParameters(rps...).delete()
}

// ViewDocuments implements the CouchDB interface.
func (db *couchdb) ViewDocuments(design, view string, rps ...Parameter) Response {
	req := newRequest(db, db.databasePath("_design", design, "_view", view), nil)
	req = req.setParameters(rps...)
	if len(req.keys) > 0 {
		return req.post()
	}
	return req.get()
}

// databasePath creates a path containing the passed
// elements based on the path of the database.
func (db *couchdb) databasePath(parts ...string) string {
	fullParts := append([]string{db.database}, parts...)
	return "/" + strings.Join(fullParts, "/")
}

// idAndRevision retrieves the ID and the revision of the
// passed document.
func (db *couchdb) idAndRevision(doc interface{}) (string, string, error) {
	marshalled, err := json.Marshal(doc)
	if err != nil {
		return "", "", errors.Annotate(err, ErrMarshallingDoc, errorMessages)
	}
	metadata := &struct {
		DocumentID       string `json:"_id,omitempt"`
		DocumentRevision string `json:"_rev,omitempty"`
	}{}
	if err = json.Unmarshal(marshalled, metadata); err != nil {
		return "", "", errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return metadata.DocumentID, metadata.DocumentRevision, nil
}

// EOF
