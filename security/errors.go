// Tideland Go CouchDB Client - Security - Errors
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
	"github.com/tideland/golib/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	ErrNoSession = iota + 1
	ErrUserExists
)

var errorMessages = errors.Messages{
	ErrNoSession:  "command needs authenticated session",
	ErrUserExists: "user already exists",
}

// EOF
