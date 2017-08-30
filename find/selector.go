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

// FieldType describes the valid types of document fields.
type FieldType string

// List of valid field types.
const (
	FieldTypeNull    FieldType = "null"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeNumber  FieldType = "number"
	FieldTypeString  FieldType = "string"
	FieldTypeArray   FieldType = "array"
	FieldTypeObject  FieldType = "object"
)

// Operators expecting an array.
var arrayOperators = map[string]bool{
	"$and": true,
	"$or":  true,
	"$nor": true,
	"$all": true,
	"$in":  true,
	"$nin": true,
}

//--------------------
// NEGATABLE
//--------------------

// Negatable allows to negate a selector.
type Negatable interface {
	// Not negates a selector.
	Not()
}

//--------------------
// CRITERION
//--------------------

type Criterion interface {
	// Negatable allows to negate this selector.
	Negatable

	// Marshaler allows to write this selector in its JSON encoding.
	json.Marshaler
}

type criterion struct {
	not       bool
	field     string
	operator  string
	arguments []interface{}
}

func newValuesCriterion(field, operator string, values ...interface{}) *criterion {
	return &criterion{
		field:     field,
		operator:  operator,
		arguments: values,
	}
}

func newCriteriaCriterion(field, operator string, criteria ...Criterion) *criterion {
	arguments := make([]interface{}, len(criteria))
	for i, c := range criteria {
		arguments[i] = c
	}
	return newValuesCriterion(field, operator, arguments...)
}

// Not implements Negatable.
func (c *criterion) Not() {
	c.not = true
}

// MarshalJSON implements json.Marshaler.
func (c *criterion) MarshalJSON() ([]byte, error) {
	var sbuf bytes.Buffer
	var jargs [][]byte
	var jargslen int

	// Praparations first.
	for _, argument := range c.arguments {
		jarg, err := json.Marshal(argument)
		if err != nil {
			return nil, err
		}
		jargs = append(jargs, jarg)
	}
	jargslen = len(jargs)

	// Is negated?
	if c.not {
		fmt.Fprint(&sbuf, "{\"$not\":")
	}
	// Prepend with field if needed.
	if c.field != "" {
		fmt.Fprintf(&sbuf, "{%q:", c.field)
	}
	// Now operator and argument(s).
	fmt.Fprintf(&sbuf, "{%q:", c.operator)
	if arrayOperators[c.operator] {
		fmt.Fprint(&sbuf, "[")
	}
	for i, jarg := range jargs {
		fmt.Fprintf(&sbuf, "%s", jarg)
		if i < jargslen-1 {
			fmt.Fprint(&sbuf, ",")
		}
	}
	if arrayOperators[c.operator] {
		fmt.Fprint(&sbuf, "]")
	}
	fmt.Fprint(&sbuf, "}")
	// Append closing brace if field has been prepended.
	if c.field != "" {
		fmt.Fprint(&sbuf, "}")
	}
	// Append closing brace if selector is negated.
	if c.not {
		fmt.Fprint(&sbuf, "}")
	}

	return sbuf.Bytes(), nil
}

func And(criteria ...Criterion) Criterion {
	return newCriteriaCriterion("", "$and", criteria...)
}

//--------------------
// SELECTOR
//--------------------

// Selector contains one or more conditions to find documents.
type Selector interface {
	// Equal checks if the field is equal to the argument.
	Equal(field string, argument interface{}) Negatable

	// NotEqual checks if the field is not equal to the argument.
	NotEqual(field string, argument interface{}) Negatable

	// In checks if the field is in the arguments.
	In(field string, arguments ...interface{}) Negatable

	// NotIn checks if the field is not in the arguments.
	NotIn(field string, arguments ...interface{}) Negatable

	// Size checks the length of the array addressed with field.
	Size(field string, size int) Negatable

	// All checks if the field is an array and contains all the arguments.
	All(field string, arguments ...interface{}) Negatable

	// GreaterThan checks if the field is greater than the argument.
	GreaterThan(field string, argument interface{}) Negatable

	// GreaterEqualThan checks if the field is greater or equal than the argument.
	GreaterEqualThan(field string, argument interface{}) Negatable

	// LowerThan checks if the field is greater than the argument.
	LowerThan(field string, argument interface{}) Negatable

	// LowerEqualThan checks if the field is greater or equal than the argument.
	LowerEqualThan(field string, argument interface{}) Negatable

	// Exists checks if the field exists.
	Exists(field string) Negatable

	// Type checks the type of the field.
	Type(field string, argument FieldType) Negatable

	// Modulo checks the remainder of the field devided by divisor.
	Modulo(field string, divisor, remainder int) Negatable

	// RegExp checks if the field matches the given pattern.
	RegExp(field, pattern string) Negatable

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
func SelectOr(conditioner func(s Selector)) Selector {
	return newSelector("$or", conditioner)
}

// SelectNone creates a combination selector where none of the
// selectors has to be true.
func SelectNone(conditioner func(s Selector)) Selector {
	return newSelector("$nor", conditioner)
}

// Equal implements Selector.
func (s *selector) Equal(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$eq", argument)
}

// NotEqual implements Selector.
func (s *selector) NotEqual(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$ne", argument)
}

// In implements Selector.
func (s *selector) In(field string, arguments ...interface{}) Negatable {
	return s.appendSubselector(field, "$in", arguments...)
}

// NotIn implements Selector.
func (s *selector) NotIn(field string, arguments ...interface{}) Negatable {
	return s.appendSubselector(field, "$nin", arguments...)
}

// Size implements Selector.
func (s *selector) Size(field string, size int) Negatable {
	return s.appendSubselector(field, "$size", size)
}

// All implements Selector.
func (s *selector) All(field string, arguments ...interface{}) Negatable {
	return s.appendSubselector(field, "$all", arguments...)
}

// GreaterThan implements Selector.
func (s *selector) GreaterThan(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$gt", argument)
}

// GreaterEqualThan implements Selector.
func (s *selector) GreaterEqualThan(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$gte", argument)
}

// LowerThan implements Selector.
func (s *selector) LowerThan(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$lt", argument)
}

// LowerEqualThan implements Selector.
func (s *selector) LowerEqualThan(field string, argument interface{}) Negatable {
	return s.appendSubselector(field, "$lte", argument)
}

// Exists implements Selector.
func (s *selector) Exists(field string) Negatable {
	return s.appendSubselector(field, "$exists", true)
}

// Type implements Selector.
func (s *selector) Type(field string, argument FieldType) Negatable {
	return s.appendSubselector(field, "$type", argument)
}

// Modulo implements Selector.
func (s *selector) Modulo(field string, divisor, remainder int) Negatable {
	return s.appendSubselector(field, "$mod", divisor, remainder)
}

// RegExp implements Selector.
func (s *selector) RegExp(field, pattern string) Negatable {
	return s.appendSubselector(field, "$regex", pattern)
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

	// Praparations first.
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
	fmt.Fprintf(&sbuf, "{%q:", s.operator)
	if arrayOperators[s.operator] {
		fmt.Fprint(&sbuf, "[")
	}
	for i, jarg := range jargs {
		fmt.Fprintf(&sbuf, "%s", jarg)
		if i < jargslen-1 {
			fmt.Fprint(&sbuf, ",")
		}
	}
	if arrayOperators[s.operator] {
		fmt.Fprint(&sbuf, "]")
	}
	fmt.Fprint(&sbuf, "}")
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
func (s *selector) appendSubselector(field, operator string, arguments ...interface{}) Negatable {
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
	if conditioner != nil {
		conditioner(s)
	}
	return s
}

// EOF
