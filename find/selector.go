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

// CombinationOperator sets how to combine multiple selectors.
type CombinationOperator int

const (
	CombineAnd CombinationOperator = iota + 1
	CombineOr
	CombineNot
	CombineNone
)

//--------------------
// FIND SELECTOR
//--------------------

// Selector contains one or more conditions to find documents.
type Selector interface {
	// Equal checks if the field is equal to the argument.
	Equal(field string, argument interface{}) Selector

	// In checks if the field is in the arguments.
	In(field string, arguments ...interface{}) Selector

	// All checks if the field is an array and contains all the arguments.
	All(field string, arguments ...interface{}) Selector

	// GreaterThan checks if the field is greater than the argument.
	GreaterThan(field string, argument interface{}) Selector

	// Marshaler allows to write a selector in its JSON encoding.
	json.Marshaler
}

// selector implements Selector.
type selector struct {
	field     string
	operator  string
	arguments []interface{}
}

// NewSelector creates a selector based on the given
// combination operator.
func NewSelector(co CombinationOperator, selectors ...Selector) Selector {
	// Get combination operator.
	ops := map[CombinationOperator]string{
		CombineAnd:  "$and",
		CombineOr:   "$or",
		CombineNot:  "$not",
		CombineNone: "$nor",
	}
	op, ok := ops[co]
	if !ok {
		op = "$and"
	}
	// Create selector.
	s := &selector{
		operator: op,
	}
	for _, subselector := range selectors {
		s.arguments = append(s.arguments, subselector)
	}
	return s
}

// Equal implements Selector.
func (s *selector) Equal(field string, argument interface{}) Selector {
	s.arguments = append(s.arguments, &selector{
		field:     field,
		operator:  "$eq",
		arguments: []interface{}{argument},
	})
	return s
}

// In implements Selector.
func (s *selector) In(field string, arguments ...interface{}) Selector {
	s.arguments = append(s.arguments, &selector{
		field:     field,
		operator:  "$in",
		arguments: arguments,
	})
	return s
}

// All implements Selector.
func (s *selector) All(field string, arguments ...interface{}) Selector {
	s.arguments = append(s.arguments, &selector{
		field:     field,
		operator:  "$all",
		arguments: arguments,
	})
	return s
}

// GreaterThan implements Selector.
func (s *selector) GreaterThan(field string, argument interface{}) Selector {
	s.arguments = append(s.arguments, &selector{
		field:     field,
		operator:  "$gt",
		arguments: []interface{}{argument},
	})
	return s
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

	// Prepend with field if needed.
	if s.field != "" {
		fmt.Fprintf(&sbuf, "{%q:", s.field)
	}
	// Now operator and argument8s).
	fmt.Fprintf(&sbuf, "{%q:", s.operator)
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
	fmt.Fprint(&sbuf, "}")
	// Append closing brace if field has been prepended.
	if s.field != "" {
		fmt.Fprint(&sbuf, "}")
	}

	return sbuf.Bytes(), nil
}

// EOF
