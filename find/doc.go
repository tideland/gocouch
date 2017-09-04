// Tideland Go CouchDB Client - Find
//
// Copyright (C) 2016-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package find of the Tideland Go CouchDB Client allows to create Mango
// Query selectors the Go way and search for documents via Find().
//
// A selector is created using a number of functions returning criteria.
//
//     selector := find.Select(find.Or(
//         find.And(
//             find.LowerThan("age", 30),
//             find.Equal("active", false),
//         ),
//         find.And(
//             find.GreaterThan("age", 60),
//             find.Equal("active", "true"),
//         ),
//     ))
//     frs := find.Find(cdb, selector, find.Fields("name", "age", "active"))
//
// Results can be retrieved by iterating over the result set.
//
//     err := frs.Do(func(document couchdb.Unmarshable) error {
//         fields := struct {
//             Name   string `json:"name"`
//             Age    int    `json:"age"`
//             Active bool   `json:"active"`
//         }{}
//         if err := document.Unmarshal(&fields); err != nil {
//             return err
//         }
//         ...
//         return nil
//     })
//
// More parameters allow restrictions to fields, sorting, filtering, and paging.
package find

// EOF
