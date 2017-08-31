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
	arguments []interface{}
}

// newValuesCriterion creates a createrion with N values as arguments.
func newValuesCriterion(field, operator string, values ...interface{}) *criterion {
	return &criterion{
		field:     field,
		operator:  operator,
		arguments: values,
	}
}

// newCriteriaCriterion creates a createrion with N criteria as arguments.
func newCriteriaCriterion(field, operator string, criteria ...Criterion) *criterion {
	arguments := make([]interface{}, len(criteria))
	for i, c := range criteria {
		arguments[i] = c
	}
	return newValuesCriterion(field, operator, arguments...)
}

// Not implements Criterion.
func (c *criterion) Not() Criterion {
	c.not = true
	return c
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

// And creates a criterion where all sub-criteria have to be true.
func And(criteria ...Criterion) Criterion {
	return newCriteriaCriterion("", "$and", criteria...)
}

// Or creates a criterion where any sub-criteria have to be true.
func Or(criteria ...Criterion) Criterion {
	return newCriteriaCriterion("", "$or", criteria...)
}

// None creates a criterion where none of the sub-criteria may be true.
func None(criteria ...Criterion) Criterion {
	return newCriteriaCriterion("", "$nor", criteria...)
}

// Exists checks if the field exists.
func Exists(field string) Criterion {
	return newValuesCriterion(field, "$exists", true)
}

// Type checks the type of the field.
func Type(field string, fieldType FieldType) Criterion {
	return newValuesCriterion(field, "$type", fieldType)
}

// Equal checks if the field is equal to the value.
func Equal(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$eq", value)
}

// Equal checks if the field is not equal to the value.
func NotEqual(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$ne", value)
}

// Size checks the length of the array addressed with field.
func Size(field string, size int) Criterion {
	return newValuesCriterion(field, "$size", size)
}

// In checks if the field contains one of the values.
func In(field string, values ...interface{}) Criterion {
	return newValuesCriterion(field, "$in", values...)
}

// NotIn checks if the field contains none of the values.
func NotIn(field string, values ...interface{}) Criterion {
	return newValuesCriterion(field, "$nin", values...)
}

// All checks if the field contains all of the values.
func All(field string, values ...interface{}) Criterion {
	return newValuesCriterion(field, "$all", values...)
}

// GreaterThan checks if the field is greater than the value.
func GreaterThan(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$gt", value)
}

// GreaterEqualThan checks if the field is greater or equal than the value.
func GreaterEqualThan(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$gte", value)
}

// LowerThan checks if the field is lower than the value.
func LowerThan(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$lt", value)
}

// LowerEqualThan checks if the field is loweer or equal than the value.
func LowerEqualThan(field string, value interface{}) Criterion {
	return newValuesCriterion(field, "$lte", value)
}

// Modulo checks the remainder of the field devided by divisor.
func Modulo(field string, divisor, remainder int) Criterion {
	return newValuesCriterion(field, "$mod", divisor, remainder)
}

// RegExp checks if the field matches the given pattern.
func RegExp(field, pattern string) Criterion {
	return newValuesCriterion(field, "$regex", pattern)
}

//--------------------
// SELECTOR
//--------------------

// Selector contains one or more criteria to find documents.
type Selector json.Marshaler

// selector implements Selector.
type selector []Criterion

// MarshalJSON implements json.Marshaler.
func (s selector) MarshalJSON() ([]byte, error) {
	// Special case: Only one criterion.
	if len(s) == 1 {
		return s[0].MarshalJSON()
	}
	// Regular case.
	var buf bytes.Buffer
	var slen = len(s)

	fmt.Fprint(&buf, "{")
	for i, c := range s {
		b, err := c.MarshalJSON()
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

// Select creates a selector based on the passed criteria.
func Select(criteria ...Criterion) Selector {
	return selector(criteria)
}

// EOF
