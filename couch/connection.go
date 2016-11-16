// Tideland Go CouchDB Client - Couch - Connection
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couch

//--------------------
// IMPORTS
//--------------------

import (
	"net/url"
)

//--------------------
// CONNECTION
//--------------------

type Connection interface {

}

type connection struct {
	location *url.URL
}

func Open(location *url.URL) (Connection, error) {
	conn := &connection{
		location: location,
	}
	return conn, nil
}


// EOF
