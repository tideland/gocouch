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
	// StartSession starts a cookie based session for the given user.
	StartSession(userID, password string) (Session, error)

	// AllDatabases returns a list of all database IDs
	// of the connected server.
	AllDatabases() ([]string, error)

	// HasDatabase checks if the configured database exists.
	HasDatabase() (bool, error)

	// CreateDatabase creates the configured database.
	CreateDatabase() ResultSet

	// DeleteDatabase removes the configured database.
	DeleteDatabase() ResultSet

	// AllDesigns returns the list of all design
	// document IDs of the configured database.
	AllDesigns() ([]string, error)

	// Design returns the design document instance for
	// the given ID.
	Design(id string) (Design, error)

	// AllDocuments returns a list of all document IDs
	// of the configured database.
	AllDocuments() ([]string, error)

	// HasDocument checks if the document with the ID exists.
	HasDocument(id string) (bool, error)

	// CreateDocument creates a new document.
	CreateDocument(doc interface{}, rps ...Parameter) ResultSet

	// ReadDocument reads an existing document.
	ReadDocument(id string, rps ...Parameter) ResultSet

	// UpdateDocument update an existing document.
	UpdateDocument(doc interface{}, rps ...Parameter) ResultSet

	// DeleteDocument deletes an existing document.
	DeleteDocument(doc interface{}, rps ...Parameter) ResultSet

	// BulkWriteDocuments allows to create or update many
	// documents en bloc.
	BulkWriteDocuments(docs ...interface{}) (Statuses, error)

	// View performs a view request.
	View(design, view string, rps ...Parameter) ViewResultSet
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

// StartSession implements the CouchDB interface.
func (cdb *couchdb) StartSession(userID, password string) (Session, error) {
	authentication := couchdbAuthentication{
		UserID:   userID,
		Password: password,
	}
	req := newRequest(cdb, "/_session", authentication)
	rs := req.post()
	return newSession(rs)
}

// AllDatabases implements the CouchDB interface.
func (cdb *couchdb) AllDatabases() ([]string, error) {
	req := newRequest(cdb, "/_all_dbs", nil)
	rs := req.get()
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	ids := []string{}
	err := rs.Document(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// HasDatabase implements the CouchDB interface.
func (cdb *couchdb) HasDatabase() (bool, error) {
	req := newRequest(cdb, cdb.databasePath(), nil)
	rs := req.head()
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDatabase implements the CouchDB interface.
func (cdb *couchdb) CreateDatabase() ResultSet {
	req := newRequest(cdb, cdb.databasePath(), nil)
	return req.put()
}

// DeleteDatabase implements the CouchDB interface.
func (cdb *couchdb) DeleteDatabase() ResultSet {
	req := newRequest(cdb, cdb.databasePath(), nil)
	return req.delete()
}

// AllDesigns implements the CouchDB interface.
func (cdb *couchdb) AllDesigns() ([]string, error) {
	req := newRequest(cdb, cdb.databasePath("_all_docs"), nil)
	rs := req.apply(StartEndKey("_design/", "_design0")).get()
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	vr := viewResult{}
	err := rs.Document(&vr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range vr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// Design implements the CouchDB interface.
func (cdb *couchdb) Design(id string) (Design, error) {
	return newDesign(cdb, id)
}

// AllDocuments implements the CouchDB interface.
func (cdb *couchdb) AllDocuments() ([]string, error) {
	req := newRequest(cdb, cdb.databasePath("_all_docs"), nil)
	rs := req.get()
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	vr := viewResult{}
	err := rs.Document(&vr)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range vr.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// HasDocument implements the CouchDB interface.
func (cdb *couchdb) HasDocument(id string) (bool, error) {
	req := newRequest(cdb, cdb.databasePath(id), nil)
	rs := req.head()
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDocument implements the CouchDB interface.
func (cdb *couchdb) CreateDocument(doc interface{}, rps ...Parameter) ResultSet {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	req := newRequest(cdb, cdb.databasePath(id), doc)
	return req.apply(rps...).put()
}

// ReadDocument implements the CouchDB interface.
func (cdb *couchdb) ReadDocument(id string, rps ...Parameter) ResultSet {
	req := newRequest(cdb, cdb.databasePath(id), nil)
	return req.apply(rps...).get()
}

// UpdateDocument implements the CouchDB interface.
func (cdb *couchdb) UpdateDocument(doc interface{}, rps ...Parameter) ResultSet {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		return newResultSet(nil, errors.New(ErrNoIdentifier, errorMessages))
	}
	req := newRequest(cdb, cdb.databasePath(id), doc)
	return req.apply(rps...).put()
}

// DeleteDocument implements the CouchDB interface.
func (cdb *couchdb) DeleteDocument(doc interface{}, rps ...Parameter) ResultSet {
	id, rev, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	rps = append(rps, Revision(rev))
	req := newRequest(cdb, cdb.databasePath(id), nil)
	return req.apply(rps...).delete()
}

// BulkWriteDocuments implements the CouchDB interface.
func (cdb *couchdb) BulkWriteDocuments(docs ...interface{}) (Statuses, error) {
	bulk := &couchdbBulkDocuments{
		Docs: docs,
	}
	req := newRequest(cdb, cdb.databasePath("_bulk_docs"), bulk)
	rs := req.post()
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	statuses := Statuses{}
	err := rs.Document(&statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// View implements the CouchDB interface.
func (cdb *couchdb) View(design, view string, rps ...Parameter) ViewResultSet {
	req := newRequest(cdb, cdb.databasePath("_design", design, "_view", view), nil)
	req = req.apply(rps...)
	var rs ResultSet
	if len(req.keys) > 0 {
		rs = req.post()
	}
	rs = req.get()
	return newView(rs)
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
