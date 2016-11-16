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
	// AllDatabases returns a list of all databases
	// of the connected server.
	AllDatabases() ([]string, error)

}

// connection implements Connection.
type connection struct {
	url *url.URL
}

// Open returns a connection to a CouchDB server. If the
// database name is part of the URL that database will
// be used.
func Open(url *url.URL) (Connection, error) {
	conn := &connection{
		url: url,
	}
	return conn, nil
}

// AllDatabases implements connection.
func (conn *connection) AllDatabases() ([]string, error) {
	req := newRequest(conn.url, methGet, "/_all_dbs", nil)
	resp, err := req.do()
	if err != nil {
		return nil, err
	}
	res := []string{}
	err = resp.ResultValue(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateDB() (resp *Response, code int, err error) {
	req, err := c.NewRequest("PUT", c.UrlString(c.DBPath(), nil), nil, nil)
	if err != nil {
		return
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}

	return
}

// EOF
