// Tideland Go CouchDB Client - Find - Index
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// INDEX
//--------------------

// Index defines the needed information for creation of an index.
type Index interface {
	// Parameters returns the parameters of the index to create.
	Parameters() []Parameter
}

// index implements index.
type index struct {
	parameters []Parameter
}

// NewIndex returns a new index containing the referenced fields.
func NewIndex(fields ...string) Index {
	idx := &index{}
	idx.parameters = append(idx.parameters, Fields(fields...))
	return idx
}

// Parameters implements Index.
func (idx *index) Parameters() []Parameter {
	return idx.parameters
}

// CreateIndex creates a new index for finds. The parameter
// Fields is needed while others are optional.
func CreateIndex(cdb couchdb.CouchDB, index Index) error {
	// Create request object.
	idxReq := request{}
	idxReq.apply(index.Parameters()...)
	req := request{}
	req.SetParameter("index", idxReq)
	// Perform index command.
	rs := cdb.Post(cdb.DatabasePath("_index"), req)
	return rs.Error()
}

// EOF
