// Tideland GoCouch - CouchDB - Design
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// DESIGN
//--------------------

// Design provides convenient access to a design document.
type Design interface {
	// ID returns the ID of the design.
	ID() string

	// Language returns the language for views and shows.
	Language() string

	// Language sets the language for views and shows.
	SetLanguage(language string)

	// View returns the map and the reduce functions of the
	// view with the ID, otherwise false.
	View(id string) (string, string, bool)

	// SetView sets the map and the reduce functions of the
	// view with the ID.
	SetView(id, mapf, reducef string)

	// Show returns the show function with the ID, otherwise false.
	Show(id string) (string, bool)

	// SetShow sets the show function with the ID.
	SetShow(id, showf string)

	// Write creates a new design document or updates an
	// existing one.
	Write(rps ...Parameter) ResultSet

	// Delete a design document.
	Delete(rps ...Parameter) ResultSet
}

// design implements the Design interface.
type design struct {
	cdb      *couchdb
	id       string
	document *designDocument
}

// newDesign creates a design instance.
func newDesign(cdb *couchdb, id string) (*design, error) {
	designID := "_design/" + id
	ok, err := cdb.HasDocument(designID)
	if err != nil {
		return nil, err
	}
	document := designDocument{}
	if ok {
		// Read the design document.
		resp := cdb.ReadDocument(designID)
		if !resp.IsOK() {
			return nil, resp.Error()
		}
		err = resp.Document(&document)
		if err != nil {
			return nil, err
		}
	} else {
		// Create the design document.
		document = designDocument{
			ID:       designID,
			Language: "javascript",
		}
	}
	d := &design{
		cdb:      cdb,
		id:       id,
		document: &document,
	}
	return d, nil
}

// ID implements the Design interface.
func (d *design) ID() string {
	return d.id
}

// Language implements the Design interface.
func (d *design) Language() string {
	return d.document.Language
}

// SetLanguage implements the Design interface.
func (d *design) SetLanguage(language string) {
	d.document.Language = language
}

// View implements the Design interface.
func (d *design) View(id string) (string, string, bool) {
	if d.document.Views == nil {
		d.document.Views = designViews{}
	}
	view, ok := d.document.Views[id]
	if !ok {
		return "", "", false
	}
	return view.Map, view.Reduce, true
}

// SetView implements the Design interface.
func (d *design) SetView(id, mapf, reducef string) {
	if d.document.Views == nil {
		d.document.Views = designViews{}
	}
	d.document.Views[id] = designView{
		Map:    mapf,
		Reduce: reducef,
	}
}

// Show implements the Design interface.
func (d *design) Show(id string) (string, bool) {
	if d.document.Shows == nil {
		d.document.Shows = map[string]string{}
	}
	show, ok := d.document.Shows[id]
	if !ok {
		return "", false
	}
	return show, true
}

// SetShow implements the Design interface.
func (d *design) SetShow(id, showf string) {
	if d.document.Shows == nil {
		d.document.Shows = map[string]string{}
	}
	d.document.Shows[id] = showf
}

// Write implements the Design interface.
func (d *design) Write(rps ...Parameter) ResultSet {
	if d.document.Revision == "" {
		return d.cdb.CreateDocument(d.document, rps...)
	}
	return d.cdb.UpdateDocument(d.document, rps...)
}

// Delete implements the Design interface.
func (d *design) Delete(rps ...Parameter) ResultSet {
	return d.cdb.DeleteDocument(d.document, rps...)
}

//--------------------
// DESIGN DOCUMENT
//--------------------

// designView defines a view inside a design document.
type designView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type designViews map[string]designView

// designAttachment defines an attachment inside a design document.
type designAttachment struct {
	Stub        bool   `json:"stub,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Length      int    `json:"length,omitempty"`
}

type designAttachments map[string]designAttachment

// designDocument contains the data of view design documents.
type designDocument struct {
	ID                     string            `json:"_id"`
	Revision               string            `json:"_rev,omitempty"`
	Language               string            `json:"language,omitempty"`
	ValidateDocumentUpdate string            `json:"validate_doc_update,omitempty"`
	Views                  designViews       `json:"views,omitempty"`
	Shows                  map[string]string `json:"shows,omitempty"`
	Attachments            designAttachments `json:"_attachments,omitempty"`
	Signatures             map[string]string `json:"signatures,omitempty"`
	Libraries              interface{}       `json:"libs,omitempty"`
}

// EOF
