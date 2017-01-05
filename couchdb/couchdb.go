// Tideland Go CouchDB Client - CouchDB
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"reflect"
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
	// Path creates a document path starting at root.
	Path(parts ...string) string

	// DatabasePath creates a document path for the database.
	DatabasePath(parts ...string) string

	// Head performs a GET request against the configured database.
	Head(path string, doc interface{}, params ...Parameter) ResultSet

	// Get performs a GET request against the configured database.
	Get(path string, doc interface{}, params ...Parameter) ResultSet

	// Put performs a GET request against the configured database.
	Put(path string, doc interface{}, params ...Parameter) ResultSet

	// Post performs a GET request against the configured database.
	Post(path string, doc interface{}, params ...Parameter) ResultSet

	// Delete performs a GET request against the configured database.
	Delete(path string, doc interface{}, params ...Parameter) ResultSet

	// AllDatabases returns a list of all database IDs
	// of the connected server.
	AllDatabases() ([]string, error)

	// HasDatabase checks if the configured database exists.
	HasDatabase() (bool, error)

	// CreateDatabase creates the configured database.
	CreateDatabase(params ...Parameter) ResultSet

	// DeleteDatabase removes the configured database.
	DeleteDatabase(params ...Parameter) ResultSet

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
	CreateDocument(doc interface{}, params ...Parameter) ResultSet

	// ReadDocument reads an existing document.
	ReadDocument(id string, params ...Parameter) ResultSet

	// UpdateDocument update an existing document.
	UpdateDocument(doc interface{}, params ...Parameter) ResultSet

	// DeleteDocument deletes an existing document.
	DeleteDocument(doc interface{}, params ...Parameter) ResultSet

	// DeleteDocumentByID deletes an existing document simply by
	// its identifier and revision.
	DeleteDocumentByID(id, revision string, params ...Parameter) ResultSet

	// BulkWriteDocuments allows to create or update many
	// documents en bloc.
	BulkWriteDocuments(docs []interface{}, params ...Parameter) (Statuses, error)

	// View performs a view request.
	View(design, view string, params ...Parameter) ViewResultSet
}

// couchdb implements CouchDB.
type couchdb struct {
	host       string
	database   string
	parameters []Parameter
}

// Open returns a configured connection to a CouchDB server.
// Permanent parameters, e.g. for authentication, are possible.
func Open(cfg etc.Etc, params ...Parameter) (CouchDB, error) {
	return OpenPath(cfg, "", params...)
}

// OpenPath returns a configured connection to a CouchDB server.
// The configuration is part of a larger configuration and the path
// leads to its location. Permanent parameters, e.g. for authentication,
// are possible.
func OpenPath(cfg etc.Etc, path string, params ...Parameter) (CouchDB, error) {
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
		parameters: params,
	}
	return cdb, nil
}

// Path implements the CouchDB interface.
func (cdb *couchdb) Path(parts ...string) string {
	return strings.Join(append([]string{""}, parts...), "/")
}

// DatabasePath implements the CouchDB interface.
func (cdb *couchdb) DatabasePath(parts ...string) string {
	return cdb.Path(append([]string{cdb.database}, parts...)...)
}

// Head implements the CouchDB interface.
func (cdb *couchdb) Head(path string, doc interface{}, params ...Parameter) ResultSet {
	req := newRequest(cdb, path, doc)
	return req.apply(params...).head()
}

// Get implements the CouchDB interface.
func (cdb *couchdb) Get(path string, doc interface{}, params ...Parameter) ResultSet {
	req := newRequest(cdb, path, doc)
	return req.apply(params...).get()
}

// Put implements the CouchDB interface.
func (cdb *couchdb) Put(path string, doc interface{}, params ...Parameter) ResultSet {
	req := newRequest(cdb, path, doc)
	return req.apply(params...).put()
}

// Post implements the CouchDB interface.
func (cdb *couchdb) Post(path string, doc interface{}, params ...Parameter) ResultSet {
	req := newRequest(cdb, path, doc)
	return req.apply(params...).post()
}

// Delete implements the CouchDB interface.
func (cdb *couchdb) Delete(path string, doc interface{}, params ...Parameter) ResultSet {
	req := newRequest(cdb, path, doc)
	return req.apply(params...).delete()
}

// AllDatabases implements the CouchDB interface.
func (cdb *couchdb) AllDatabases() ([]string, error) {
	rs := cdb.Get("/_all_dbs", nil)
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
	rs := cdb.Head(cdb.DatabasePath(), nil)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDatabase implements the CouchDB interface.
func (cdb *couchdb) CreateDatabase(params ...Parameter) ResultSet {
	return cdb.Put(cdb.DatabasePath(), nil, params...)
}

// DeleteDatabase implements the CouchDB interface.
func (cdb *couchdb) DeleteDatabase(params ...Parameter) ResultSet {
	return cdb.Delete(cdb.DatabasePath(), nil, params...)
}

// AllDesigns implements the CouchDB interface.
func (cdb *couchdb) AllDesigns() ([]string, error) {
	rs := cdb.Get(cdb.DatabasePath("_all_docs"), nil, StartEndKey("_design/", "_design0/"))
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	vr := couchdbViewResult{}
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
	rs := cdb.Get(cdb.DatabasePath("_all_docs"), nil)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	vr := couchdbViewResult{}
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
	rs := cdb.Head(cdb.DatabasePath(id), nil)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDocument implements the CouchDB interface.
func (cdb *couchdb) CreateDocument(doc interface{}, params ...Parameter) ResultSet {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	return cdb.Put(cdb.DatabasePath(id), doc, params...)
}

// ReadDocument implements the CouchDB interface.
func (cdb *couchdb) ReadDocument(id string, params ...Parameter) ResultSet {
	return cdb.Get(cdb.DatabasePath(id), nil, params...)
}

// UpdateDocument implements the CouchDB interface.
func (cdb *couchdb) UpdateDocument(doc interface{}, params ...Parameter) ResultSet {
	id, _, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		return newResultSet(nil, errors.New(ErrNoIdentifier, errorMessages))
	}
	hasDoc, err := cdb.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, errorMessages, id))
	}
	return cdb.Put(cdb.DatabasePath(id), doc, params...)
}

// DeleteDocument implements the CouchDB interface.
func (cdb *couchdb) DeleteDocument(doc interface{}, params ...Parameter) ResultSet {
	id, revision, err := cdb.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	hasDoc, err := cdb.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, errorMessages, id))
	}
	params = append(params, Revision(revision))
	return cdb.Delete(cdb.DatabasePath(id), nil, params...)
}

// DeleteDocumentByID implements the CouchDB interface.
func (cdb *couchdb) DeleteDocumentByID(id, revision string, params ...Parameter) ResultSet {
	hasDoc, err := cdb.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, errorMessages, id))
	}
	params = append(params, Revision(revision))
	return cdb.Delete(cdb.DatabasePath(id), nil, params...)
}

// BulkWriteDocuments implements the CouchDB interface.
func (cdb *couchdb) BulkWriteDocuments(docs []interface{}, params ...Parameter) (Statuses, error) {
	bulk := &couchdbBulkDocuments{
		Docs: docs,
	}
	rs := cdb.Post(cdb.DatabasePath("_bulk_docs"), bulk, params...)
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
func (cdb *couchdb) View(design, view string, params ...Parameter) ViewResultSet {
	var rs ResultSet
	req := newRequest(cdb, cdb.DatabasePath("_design", design, "_view", view), nil).apply(params...)
	if len(req.keys) > 0 {
		rs = req.post()
	} else {
		rs = req.get()
	}
	return newView(rs)
}

// idAndRevision retrieves the ID and the revision of the
// passed document.
func (cdb *couchdb) idAndRevision(doc interface{}) (string, string, error) {
	v := reflect.Indirect(reflect.ValueOf(doc))
	t := v.Type()
	k := t.Kind()
	if k != reflect.Struct {
		return "", "", errors.New(ErrInvalidDocument, errorMessages)
	}
	var id string
	var revision string
	var found int
	for i := 0; i < t.NumField(); i++ {
		vf := v.Field(i)
		tf := t.Field(i)
		if json, ok := tf.Tag.Lookup("json"); ok {
			switch json {
			case "_id", "_id,omitempty":
				id = vf.String()
				found++
			case "_rev", "_rev,omitempty":
				revision = vf.String()
				found++
			}
		}
	}
	if found != 2 {
		return "", "", errors.New(ErrInvalidDocument, errorMessages)
	}
	return id, revision, nil
}

//--------------------
// CONFIGURATION
//--------------------

// Configure creates a configuration out of
// the passed arguments.
func Configure(hostname string, port int, database string) (etc.Etc, error) {
	source := fmt.Sprintf("{etc {hostname %s}{port %d}{database %s}}", hostname, port, database)
	return etc.ReadString(source)
}

// EOF
