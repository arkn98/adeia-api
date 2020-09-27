/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package store

import (
	"context"
	"io"
)

// Getter is the interface for all GET-related methods of the database.
type Getter interface {
	GetMany(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetOne(ctx context.Context, dest interface{}, query string, args ...interface{}) (ok bool, err error)
}

// Deleter is the interface for all DELETE-related methods of the database.
type Deleter interface {
	Delete(ctx context.Context, query string, args ...interface{}) (rowsAffected int64, err error)
}

// Inserter is the interface for all INSERT-related methods of the database.
type Inserter interface {
	Insert(ctx context.Context, query string, args ...interface{}) (lastInsertID int, err error)
	InsertNamed(ctx context.Context, namedQuery string, arg interface{}) (lastInsertID int, err error)
}

// Updater is the interface for all UPDATE-related methods of the database.
type Updater interface {
	Update(ctx context.Context, query string, args ...interface{}) (rowsAffected int64, err error)
	UpdateNamed(ctx context.Context, query string, arg interface{}) (rowsAffected int64, err error)
}

// DB is the interface for all the methods of the database.
type DB interface {
	io.Closer
	Deleter
	Getter
	Inserter
	Updater
}
