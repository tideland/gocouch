// Tideland Go CouchDB Client - Find - Parameters
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package find

//--------------------
// PARAMETERIZABLE
//--------------------

// Parameterizable defines the methods needed to apply the parameters.
type Parameterizable interface {
	// SetParameter sets one of the request parameters.
	SetParameter(key string, parameter interface{})
}

//--------------------
// PARAMETERS
//--------------------

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(pa Parameterizable)

// Fields specifies which fields of each object should be returned.
// if it's omitted, the entire object is returned.
func Fields(fields ...string) Parameter {
	return func(pa Parameterizable) {
		pa.SetParameter("fields", fields)
	}
}

// Limit sets the maximum number of results returned. Default
// by database is 25.
func Limit(limit int) Parameter {
	return func(pa Parameterizable) {
		pa.SetParameter("limit", limit)
	}
}

// UseIndex instructs a query to use a specific index. Name is allowed
// be empty.
func UseIndex(designDocument, name string) Parameter {
	return func(pa Parameterizable) {
		if name == "" {
			pa.SetParameter("use_index", designDocument)
		} else {
			pa.SetParameter("use_index", []string{designDocument, name})
		}
	}
}

// Sort sets how to sort the result by ascending or descending
// sorted fields.
func Sort(fields ...Direction) Parameter {
	return func(pa Parameterizable) {
		sort := []map[string]string{}
		for _, field := range fields {
			sort = append(sort, map[string]string{
				field.Field(): field.Direction(),
			})
		}
		pa.SetParameter("sort", sort)
	}
}

//--------------------
// DIRECTION
//--------------------

// Direction controls the sorting of a find result. It contains a
// name of the field and if it should be ascending or descending.
type Direction interface {
	// Field returns the field taken for sorting.
	Field() string

	// Direction returns the direction as "asc" or "desc".
	Direction() string
}

// direction implements Direction.
type direction struct {
	field     string
	direction string
}

// Field implements Direction.
func (d direction) Field() string {
	return d.field
}

// Direction implements Direction.
func (d direction) Direction() string {
	return d.direction
}

// Ascending returns an ascending sort direction.
func Ascending(field string) Direction {
	return direction{
		field:     field,
		direction: "asc",
	}
}

// Descending returns an descending sort direction.
func Descending(field string) Direction {
	return direction{
		field:     field,
		direction: "desc",
	}
}

// EOF
