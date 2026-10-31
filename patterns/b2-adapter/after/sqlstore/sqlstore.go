// Package sqlstore will be the adapter between database/sql and orders. Your
// work: see the GOAL in ../../after.go.
package sqlstore

import "database/sql"

// New wraps db. For now it hands db straight back, so orders still holds the
// raw database.
//
// TODO: return your adapter type instead.
func New(db *sql.DB) *sql.DB { return db }
