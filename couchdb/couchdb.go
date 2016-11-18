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
// DOCUMENT INTERFACES
//--------------------

// Identifiable are document type which can provide their
// document identifier and revision.
type Identifiable interface {
	// DocumentID returns the identifier of the document.
	DocumentID() string

	// Document revision returns the revision of the document.
	DocumentRevision() string
}

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

	// AllDocuments returns a list of all document IDs
	// of the configured database.
	AllDocuments() ([]string, error)

	// CreateDocument creates a new document.
	CreateDocument(doc interface{}) Response
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
	query := NewQuery().StartEndKey("_design/", "_design0")
	req := newRequest(db, db.databasePath("_all_docs"), nil).setQuery(query)
	resp := req.get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	cvr := couchdbViewResult{}
	err := resp.ResultValue(&cvr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range cvr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// AllDocuments implements the CouchDB interface.
func (db *couchdb) AllDocuments() ([]string, error) {
	req := newRequest(db, db.databasePath("_all_docs"), nil)
	resp := req.get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	cvr := couchdbViewResult{}
	err := resp.ResultValue(&cvr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range cvr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// CreateDocument implements the CouchDB interface.
func (db *couchdb) CreateDocument(doc interface{}) Response {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	req := newRequest(db, db.databasePath(id), doc)
	return req.put()
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
	// Can the type provide it by itself?
	if identifiable, ok := doc.(Identifiable); ok {
		return identifiable.DocumentID(), identifiable.DocumentRevision(), nil
	}
	// OK, use marshalling.
	marshalled, err := json.Marshal(doc)
	if err != nil {
		return "", "", errors.Annotate(err, ErrMarshallingDoc, errorMessages)
	}
	iar := &struct {
		ID       string `json:"_id"`
		Revision string `json:"_rev"`
	}{}
	if err = json.Unmarshal(marshalled, iar); err != nil {
		return "", "", errors.Annotate(err, ErrUnmarshallingDoc, errorMessages)
	}
	return iar.ID, iar.Revision, nil
}

// EOF
