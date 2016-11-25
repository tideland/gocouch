// Tideland Go CouchDB Client - Startup - Errors
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package startup

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
	ErrIllegalVersion = iota + 1
	ErrStartupActionFailed
)

var errorMessages = errors.Messages{
	ErrIllegalVersion:      "illegal database version",
	ErrStartupActionFailed: "startup action failed for version '%v'",
}

// EOF
