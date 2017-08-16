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
	"encoding/json"
	"strings"
)

//--------------------
// CONSTANTS
//--------------------

// CombinationOperator sets how to combine multiple selectors.
type CombinationOperator int

const (
	CombineAnd CombinationOperator = iota + 1
	CombineOr
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

	// GreaterThan checks if the field is greater than the argument.
	GreaterThan(field, argument interface{})

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
	// First operator and argument(s).
	var jArguments []string
	var jArgument string
	for _, argument := range s.arguments {
		b, err := json.Marshal(argument)
		if err != nil {
			return nil, err
		}
		jArguments = append(jArguments, string(b))
	}
	if len(jArguments) == 1 {
		jArgument = jArguments[0]
	} else {
		jArgument = "[" + strings.Join(jArguments, ",") + "]"
	}
	jOperatorArgument := "{\"" + s.operator + "\":" + jArgument + "}"
	if s.field == "" {
		return []byte(jOperatorArgument), nil
	}
	jField := "{\"" + s.field + "\":" + jOperatorArgument + "}"
	return []byte(jField), nil
}

// EOF
