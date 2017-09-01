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

// Operators expecting direct fields.
var fieldOperators = map[string]bool{
	"$elemMatch": true,
	"$allMatch":  true,
}

//--------------------
// ENVELOPE
//--------------------

// envelope makes any type implementing json.Marshaler.
type envelope struct {
	value interface{}
}

// MarshalJSON implements json.Marshaler.
func (e envelope) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

//--------------------
// CRITERION
//--------------------

// Criterion defines one selector criterion. They are created by several extra
// functions. Some of those expect a field name for their operation. In case of
// MatchElement() and MatchAll() the fields of the sub-criteria have
type Criterion interface {
	// Not allows to negate this criterion.
	Not() Criterion

	// Marshaler allows to write this selector in its JSON encoding.
	json.Marshaler
}

// criterion implements Criterion.
type criterion struct {
	not       bool
	field     string
	operator  string
	arguments []json.Marshaler
}

// newCriterion creates a createrion with a number of arguments.
func newCriterion(field, operator string, arguments ...interface{}) *criterion {
	c := &criterion{
		field:     field,
		operator:  operator,
		arguments: make([]json.Marshaler, len(arguments)),
	}
	for i, argument := range arguments {
		jm, ok := argument.(json.Marshaler)
		if ok {
			c.arguments[i] = jm
			continue
		}
		c.arguments[i] = envelope{argument}
	}
	return c
}

// Not implements Criterion.
func (c *criterion) Not() Criterion {
	c.not = true
	return c
}

// MarshalJSON implements json.Marshaler.
func (c *criterion) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	var alen = len(c.arguments)

	// Is negated?
	if c.not {
		fmt.Fprint(&buf, "{\"$not\":")
	}
	// Prepend with field if needed.
	if c.field != "" {
		fmt.Fprintf(&buf, "{%q:", c.field)
	}
	// Now operator and arguments(s).
	fmt.Fprintf(&buf, "{%q:", c.operator)
	switch {
	case arrayOperators[c.operator]:
		fmt.Fprint(&buf, "[")
	case fieldOperators[c.operator]:
		fmt.Fprint(&buf, "{")
	}
	for i, argument := range c.arguments {
		b, err := argument.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if fieldOperators[c.operator] {
			buf.Write(b[1 : len(b)-1])
		} else {
			buf.Write(b)
		}
		if i < alen-1 {
			fmt.Fprint(&buf, ",")
		}
	}
	switch {
	case arrayOperators[c.operator]:
		fmt.Fprint(&buf, "]")
	case fieldOperators[c.operator]:
		fmt.Fprint(&buf, "}")
	}
	fmt.Fprint(&buf, "}")
	// Append closing brace if field has been prepended.
	if c.field != "" {
		fmt.Fprint(&buf, "}")
	}
	// Append closing brace if selector is negated.
	if c.not {
		fmt.Fprint(&buf, "}")
	}

	return buf.Bytes(), nil
}

// And creates a criterion where all sub-criteria have to be true.
func And(criteria ...Criterion) Criterion {
	return newCriterion("", "$and", criteriaToArguments(criteria...)...)
}

// Or creates a criterion where any sub-criteria have to be true.
func Or(criteria ...Criterion) Criterion {
	return newCriterion("", "$or", criteriaToArguments(criteria...)...)
}

// None creates a criterion where none of the sub-criteria may be true.
func None(criteria ...Criterion) Criterion {
	return newCriterion("", "$nor", criteriaToArguments(criteria...)...)
}

// MatchElement creates a criterion matching all documents with at least
// one array field element matching the supplied query criteria.
func MatchElement(field string, criteria ...Criterion) {
	return newCriterion(field, "$elemMatch", criteriaToArguments(criteria...)...)
}

// MatchAll creates a criterion matching all documents with all array field
// elements matching the supplied query criteria.
func MatchAll(field string, criteria ...Criterion) {
	return newCriterion(field, "$allMatch", criteriaToArguments(criteria...)...)
}

// Exists checks if the field exists.
func Exists(field string) Criterion {
	return newCriterion(field, "$exists", true)
}

// Type checks the type of the field.
func Type(field string, fieldType FieldType) Criterion {
	return newCriterion(field, "$type", fieldType)
}

// Equal checks if the field is equal to the value.
func Equal(field string, value interface{}) Criterion {
	return newCriterion(field, "$eq", value)
}

// Equal checks if the field is not equal to the value.
func NotEqual(field string, value interface{}) Criterion {
	return newCriterion(field, "$ne", value)
}

// Size checks the length of the array addressed with field.
func Size(field string, size int) Criterion {
	return newCriterion(field, "$size", size)
}

// In checks if the field contains one of the values.
func In(field string, values ...interface{}) Criterion {
	return newCriterion(field, "$in", values...)
}

// NotIn checks if the field contains none of the values.
func NotIn(field string, values ...interface{}) Criterion {
	return newCriterion(field, "$nin", values...)
}

// All checks if the field contains all of the values.
func All(field string, values ...interface{}) Criterion {
	return newCriterion(field, "$all", values...)
}

// GreaterThan checks if the field is greater than the value.
func GreaterThan(field string, value interface{}) Criterion {
	return newCriterion(field, "$gt", value)
}

// GreaterEqualThan checks if the field is greater or equal than the value.
func GreaterEqualThan(field string, value interface{}) Criterion {
	return newCriterion(field, "$gte", value)
}

// LowerThan checks if the field is lower than the value.
func LowerThan(field string, value interface{}) Criterion {
	return newCriterion(field, "$lt", value)
}

// LowerEqualThan checks if the field is loweer or equal than the value.
func LowerEqualThan(field string, value interface{}) Criterion {
	return newCriterion(field, "$lte", value)
}

// Modulo checks the remainder of the field devided by divisor.
func Modulo(field string, divisor, remainder int) Criterion {
	return newCriterion(field, "$mod", divisor, remainder)
}

// RegExp checks if the field matches the given pattern.
func RegExp(field, pattern string) Criterion {
	return newCriterion(field, "$regex", pattern)
}

//--------------------
// SELECTOR
//--------------------

// Selector contains one or more criteria to find documents.
type Selector json.Marshaler

// selector implements Selector.
type selector struct {
	criteria []Criterion
}

// Select creates a selector based on the passed criteria.
func Select(criteria ...Criterion) Selector {
	return &selector{criteria}
}

// MarshalJSON implements json.Marshaler.
func (s *selector) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	var slen = len(s.criteria)

	// Special case, only one criterion.
	if slen == 1 {
		return s.criteria[0].MarshalJSON()
	}
	// Regular case, multiple criteria.
	fmt.Fprint(&buf, "{")
	for i, criterion := range s.criteria {
		b, err := criterion.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(b[1 : len(b)-1])
		if i < slen-1 {
			fmt.Fprint(&buf, ",")
		}
	}
	fmt.Fprint(&buf, "}")

	return buf.Bytes(), nil
}

//--------------------
// HELPER
//--------------------

// criteriaToArguments converts a slice of Criterion to
// a slice of empty interfaces.
func criteriaToArguments(criteria ...Criterion) []interface{} {
	arguments := make([]interface{}, len(criteria))
	for i, criterion := range criteria {
		arguments[i] = criterion
	}
	return arguments
}

// EOF
