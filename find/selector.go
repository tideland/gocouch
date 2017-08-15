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
// FIND SELECTOR
//--------------------

// Selector contains one or more conditions to find documents.
type Selector interface {
	// SubSelectors combines extra created selectors to a new one.
	SubSelectors(sels ...Selector) Selector

	// Equal checks if the field is equal to the argument.
	Equal(field string, argument interface{}) Selector

	// Marshaler allows to write a selector in its JSON encoding.
	json.Marshaler
}

// selector implements Selector.
type selector struct {
	field     string
	operator  string
	arguments []interface{}
}

// NewAndSelector returns a new selector combining the called
// operators with "and".
func NewAndSelector() Selector {
	return &selector{
		operator: "$and",
	}
}

// NewOrSelector returns a new selector combining the called
// operators with "or".
func NewOrSelector() Selector {
	return &selector{
		operator: "$or",
	}
}

// SubSelectors implements Selector.
func (sel *selector) SubSelectors(sels ...Selector) Selector {
	for _, subsel := range sels {
		sel.arguments = append(sel.arguments, subsel)
	}
	return sel
}

func (sel *selector) Equal(field string, argument interface{}) Selector {
	sel.arguments = append(sel.arguments, &selector{
		field:     field,
		arguments: []interface{}{argument},
	})
	return sel
}

// MarshalJSON implements json.Marshaler.
func (sel *selector) MarshalJSON() ([]byte, error) {
	if sel.field != "" {
		// Single field.
		var jArguments []string
		var jArgument string
		for _, argument := range sel.arguments {
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
		jOperator := "{\"" + sel.operator + "\":" + jArgument + "}"
		jSelector := "{\"" + sel.field + "\":" + jOperator + "}"
		return []byte(jSelector), nil
	}
	// Operator.
	return nil, nil
}

// EOF
