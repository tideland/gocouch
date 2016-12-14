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
	ErrUserExists = iota + 1
)

var errorMessages = errors.Messages{
	ErrUserExists: "user already exists",
}

// EOF
