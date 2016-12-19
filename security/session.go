// Tideland Go CouchDB Client - Security - Session
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package security

//--------------------
// IMPORTS
//--------------------

import (
	"strings"

	"github.com/tideland/gocouch/couchdb"
)

//--------------------
// SESSION
//--------------------

// Session contains the information of a CouchDB session.
type Session interface {
	// UserID returns the user ID of this session.
	UserID() string

	// Cookie returns the session cookie as parameter
	// to be used in the individual database requests.
	Cookie() couchdb.Parameter

	// Stop ends the session.
	Stop() error
}

// session implements the Session interface.
type session struct {
	cdb         couchdb.CouchDB
	userID      string
	authSession string
}

// NewSession starts a cookie based session for the given user.
func NewSession(cdb couchdb.CouchDB, userID, password string) (Session, error) {
	user := User{
		UserID:   userID,
		Password: password,
	}
	rs := cdb.Post("_session", user)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	roles := couchdbRoles{}
	err := rs.Document(&roles)
	if err != nil {
		return nil, err
	}
	setCookie := rs.Header("Set-Cookie")
	authSession := ""
	for _, part := range strings.Split(setCookie, ";") {
		if strings.HasPrefix(part, "AuthSession=") {
			authSession = part
			break
		}
	}
	s := &session{
		cdb:         cdb,
		userID:      roles.UserID,
		authSession: authSession,
	}
	return s, nil
}

// UserID implements the Session interface.
func (s *session) UserID() string {
	return s.userID
}

// Cookie implements the Session interface.
func (s *session) Cookie() couchdb.Parameter {
	return func(pa couchdb.Parameterizable) {
		pa.SetHeader("Cookie", s.authSession)
		pa.SetHeader("X-CouchDB-WWW-Authenticate", "Cookie")
	}
}

// Stop implements the Session interface.
func (s *session) Stop() error {
	rs := s.cdb.Delete("/_session", nil, s.Cookie())
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// EOF
