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

	// StartSession starts a cookie based session for the given user.
	// StartSession(UserID, password string) (*Session, error)

	// HasDatabase checks if the configured database exists.
	HasDatabase() (bool, error)

	// CreateDatabase creates the configured database.
	CreateDatabase() Response

	// DeleteDatabase removes the configured database.
	DeleteDatabase() Response

	// AllDesignDocuments returns the lsit of all design
	// document IDs of the configured database.
	AllDesignDocuments() ([]string, error)

	// CreateDesignDocument creates a new design document.
	CreateDesignDocument(doc *DesignDocument, rps ...Parameter) Response

	// ReadDesignDocument reads an existing design document.
	ReadDesignDocument(id string, rps ...Parameter) (*DesignDocument, error)

	// UpdateDesignDocument update an existing design document.
	UpdateDesignDocument(doc *DesignDocument, rps ...Parameter) Response

	// DeleteDesignDocument deletes an existing design document.
	DeleteDesignDocument(doc *DesignDocument, rps ...Parameter) Response

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

	// BulkWriteDocuments allows to create or update many
	// documents en bloc.
	BulkWriteDocuments(docs ...interface{}) (BulkResults, error)

	// ViewDocuments reads the output of a view.
	ViewDocuments(design, view string, rps ...Parameter) Response
}

// couchdb implements CouchDB.
type couchdb struct {
	host       string
	database   string
	parameters []Parameter
}

// Open returns a configured connection to a CouchDB server.
// Permanent parameters, e.g. for authentication, are possible.
func Open(cfg etc.Etc, rps ...Parameter) (CouchDB, error) {
	return OpenPath(cfg, "", rps...)
}

// OpenPath returns a configured connection to a CouchDB server.
// The configuration is part of a larger configuration and the path
// leads to its location. Permanent parameters, e.g. for authentication,
// are possible.
func OpenPath(cfg etc.Etc, path string, rps ...Parameter) (CouchDB, error) {
	if cfg == nil {
		return nil, errors.New(ErrNoConfiguration, errorMessages)
	}
	if path != "" {
		var err error
		cfg, err = cfg.Split(path)
		if err != nil {
			return nil, errors.New(ErrNoConfiguration, errorMessages)
		}
	}
	host := fmt.Sprintf("%s:%d",
		cfg.ValueAsString("hostname", "localhost"),
		cfg.ValueAsInt("port", 5984),
	)
	cdb := &couchdb{
		host:       host,
		database:   cfg.ValueAsString("database", "default"),
		parameters: rps,
	}
	return cdb, nil
}

// AllDatabases implements connection.
func (cdb *couchdb) AllDatabases() ([]string, error) {
	req := newRequest(cdb, "/_all_dbs", nil)
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

// HasDatabase implements the CouchDB interface.
func (cdb *couchdb) HasDatabase() (bool, error) {
	req := newRequest(cdb, cdb.databasePath(), nil)
	resp := req.head()
	if resp.IsOK() {
		return true, nil
	}
	if resp.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, resp.Error()
}

// CreateDatabase implements the CouchDB interface.
func (cdb *couchdb) CreateDatabase() Response {
	req := newRequest(cdb, cdb.databasePath(), nil)
	return req.put()
}

// DeleteDatabase implements the CouchDB interface.
func (cdb *couchdb) DeleteDatabase() Response {
	req := newRequest(cdb, cdb.databasePath(), nil)
	return req.delete()
}

// AllDesignDocuments implements the CouchDB interface.
func (cdb *couchdb) AllDesignDocuments() ([]string, error) {
	req := newRequest(cdb, cdb.databasePath("_all_docs"), nil)
	resp := req.apply(StartEndKey("_design/", "_design0")).get()
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
func (cdb *couchdb) CreateDesignDocument(doc *DesignDocument, rps ...Parameter) Response {
	if !strings.HasPrefix(doc.ID, "_design/") {
		doc.ID = "_design/" + doc.ID
	}
	req := newRequest(cdb, cdb.databasePath(doc.ID), doc)
	return req.apply(rps...).put()
}

// ReadDesignDocument implements the CouchDB interface.
func (cdb *couchdb) ReadDesignDocument(id string, rps ...Parameter) (*DesignDocument, error) {
	if !strings.HasPrefix(id, "_design/") {
		id = "_design/" + id
	}
	req := newRequest(cdb, cdb.databasePath(id), nil)
	resp := req.apply(rps...).get()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	dd := DesignDocument{}
	err := resp.ResultValue(&dd)
	if err != nil {
		return nil, err
	}
	return &dd, nil
}

// UpdateDesignDocument implements the CouchDB interface.
func (cdb *couchdb) UpdateDesignDocument(doc *DesignDocument, rps ...Parameter) Response {
	req := newRequest(cdb, cdb.databasePath(doc.ID), doc)
	return req.apply(rps...).put()
}

// DeleteDesignDocument implements the CouchDB interface.
func (cdb *couchdb) DeleteDesignDocument(doc *DesignDocument, rps ...Parameter) Response {
	rps = append(rps, Revision(doc.Revision))
	req := newRequest(cdb, cdb.databasePath(doc.ID), nil)
	return req.apply(rps...).delete()
}

// AllDocuments implements the CouchDB interface.
func (cdb *couchdb) AllDocuments() ([]string, error) {
	req := newRequest(cdb, cdb.databasePath("_all_docs"), nil)
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
func (cdb *couchdb) CreateDocument(doc interface{}, rps ...Parameter) Response {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	req := newRequest(cdb, cdb.databasePath(id), doc)
	return req.apply(rps...).put()
}

// ReadDocument implements the CouchDB interface.
func (cdb *couchdb) ReadDocument(id string, rps ...Parameter) Response {
	req := newRequest(cdb, cdb.databasePath(id), nil)
	return req.apply(rps...).get()
}

// UpdateDocument implements the CouchDB interface.
func (cdb *couchdb) UpdateDocument(doc interface{}, rps ...Parameter) Response {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	if id == "" {
		return newResponse(nil, errors.New(ErrNoIdentifier, errorMessages))
	}
	req := newRequest(cdb, cdb.databasePath(id), doc)
	return req.apply(rps...).put()
}

// DeleteDocument implements the CouchDB interface.
func (cdb *couchdb) DeleteDocument(doc interface{}, rps ...Parameter) Response {
	id, rev, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResponse(nil, err)
	}
	rps = append(rps, Revision(rev))
	req := newRequest(cdb, cdb.databasePath(id), nil)
	return req.apply(rps...).delete()
}

// BulkWriteDocuments implements the CouchDB interface.
func (cdb *couchdb) BulkWriteDocuments(docs ...interface{}) (BulkResults, error) {
	bulk := &couchdbBulkDocuments{
		Docs: docs,
	}
	req := newRequest(cdb, cdb.databasePath("_bulk_docs"), bulk)
	resp := req.post()
	if !resp.IsOK() {
		return nil, resp.Error()
	}
	brs := BulkResults{}
	err := resp.ResultValue(&brs)
	if err != nil {
		return nil, err
	}
	return brs, nil
}

// ViewDocuments implements the CouchDB interface.
func (cdb *couchdb) ViewDocuments(design, view string, rps ...Parameter) Response {
	req := newRequest(cdb, cdb.databasePath("_design", design, "_view", view), nil)
	req = req.apply(rps...)
	if len(req.keys) > 0 {
		return req.post()
	}
	return req.get()
}

// databasePath creates a path containing the passed
// elements based on the path of the database.
func (cdb *couchdb) databasePath(parts ...string) string {
	fullParts := append([]string{cdb.database}, parts...)
	return "/" + strings.Join(fullParts, "/")
}

// idAndRevision retrieves the ID and the revision of the
// passed document.
func (cdb *couchdb) idAndRevision(doc interface{}) (string, string, error) {
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
