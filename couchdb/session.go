// Tideland Go CouchDB Client - CouchDB - Session
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
	"strings"
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
	Cookie() Parameter
}

// session implements the Session interface.
type session struct {
	userID      string
	authSession string
}

// newSession creates a new session instance.
func newSession(rs *resultSet) (*session, error) {
	roles := couchdbRoles{}
	err := rs.Document(&roles)
	if err != nil {
		return nil, err
	}
	setCookie := rs.header("Set-Cookie")
	authSession := ""
	for _, part := range strings.Split(setCookie, ";") {
		if strings.HasPrefix(part, "AuthSession=") {
			authSession = part
			break
		}
	}
	s := &session{
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
func (s *session) Cookie() Parameter {
	return func(pa Parameterizable) {
		pa.SetHeader("Cookie", s.authSession)
	}
}

// EOF
