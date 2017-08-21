// Tideland Go CouchDB Client - Find - Selector
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
	"bytes"
	"encoding/json"
	"fmt"
)

//--------------------
// CONSTANTS
//--------------------

//--------------------
// SELECTOR
//--------------------

// Negatable allows to negate a selector.
type Negatable interface {
	// Not negates a selector.
	Not()
}

// Selector contains one or more conditions to find documents.
type Selector interface {
	// Equal checks if the field is equal to the argument.
	Equal(field string, argument interface{}) Negatable

	// NotEqual checks if the field is not equal to the argument.
	NotEqual(field string, argument interface{}) Negatable

	// In checks if the field is in the arguments.
	In(field string, arguments ...interface{}) Negatable

	// All checks if the field is an array and contains all the arguments.
	All(field string, arguments ...interface{}) Negatable

	// GreaterThan checks if the field is greater than the argument.
	GreaterThan(field string, argument interface{}) Negatable

	// Append adds an other selector to this one.
	Append(subselector Selector)

	// Negatable allows to negate this selector.
	Negatable

	// Marshaler allows to write this selector in its JSON encoding.
	json.Marshaler
}

// selector implements Selector.
type selector struct {
	not       bool
	field     string
	operator  string
	arguments []interface{}
}

// SelectAnd creates a combination selector where all selectors
// have to be true.
func SelectAnd(conditioner func(s Selector)) Selector {
	return newSelector("$and", conditioner)
}

// SelectOr creates a combination selector where any selectors
// have to be true.
func SelectOr(conditioner func(as Selector)) Selector {
	return newSelector("$or", conditioner)
}

// Equal implements Selector.
func (s *selector) Equal(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$eq", []interface{}{argument})
}

// NotEqual implements Selector.
func (s *selector) NotEqual(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$ne", []interface{}{argument})
}

// In implements Selector.
func (s *selector) In(field string, arguments ...interface{}) Negatable {
	return s.appendSubselector(field, "$in", arguments)
}

// All implements Selector.
func (s *selector) All(field string, arguments ...interface{}) Negatable {
	return s.appendSubselector(field, "$all", arguments)
}

// GreaterThan implements Selector.
func (s *selector) GreaterThan(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$gt", []interface{}{argument})
}

// Append implements Selector.
func (s *selector) Append(subselector Selector) {
	s.arguments = append(s.arguments, subselector)
}

// Not implements Negatable.
func (s *selector) Not() {
	s.not = true
}

// MarshalJSON implements json.Marshaler.
func (s *selector) MarshalJSON() ([]byte, error) {
	var sbuf bytes.Buffer
	var jargs [][]byte
	var jargslen int

	for _, argument := range s.arguments {
		jarg, err := json.Marshal(argument)
		if err != nil {
			return nil, err
		}
		jargs = append(jargs, jarg)
	}
	jargslen = len(jargs)

	// Is negated?
	if s.not {
		fmt.Fprint(&sbuf, "{\"$not\":")
	}
	// Prepend with field if needed.
	if s.field != "" {
		fmt.Fprintf(&sbuf, "{%q:", s.field)
	}
	// Now operator and argument(s).
	if s.operator != "" {
		fmt.Fprintf(&sbuf, "{%q:", s.operator)
	}
	if jargslen > 1 {
		fmt.Fprintf(&sbuf, "[")
	}
	for i, jarg := range jargs {
		fmt.Fprintf(&sbuf, "%s", jarg)
		if i < jargslen-1 {
			fmt.Fprint(&sbuf, ",")
		}
	}
	if jargslen > 1 {
		fmt.Fprint(&sbuf, "]")
	}
	if s.operator != "" {
		fmt.Fprint(&sbuf, "}")
	}
	// Append closing brace if field has been prepended.
	if s.field != "" {
		fmt.Fprint(&sbuf, "}")
	}
	// Append closing brace if selector is negated.
	if s.not {
		fmt.Fprint(&sbuf, "}")
	}

	return sbuf.Bytes(), nil
}

// appendSubselector creates and appends a subselector based on field,
// operator, and arguments.
func (s *selector) appendSubselector(field, operator string, arguments []interface{}) Negatable {
	subselector := &selector{
		field:     field,
		operator:  operator,
		arguments: arguments,
	}
	s.Append(subselector)
	return subselector
}

//--------------------
// HELPERS
//--------------------

// newSelector helps to create new selectors.
func newSelector(operator string, conditioner func(Selector)) Selector {
	s := &selector{
		operator: operator,
	}
	conditioner(s)
	return s
}

// EOF
