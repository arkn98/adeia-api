/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pg

import (
	"context"
	"database/sql"
	"strconv"

	"adeia/internal/config"
	"adeia/pkg/util/ioutil"
	"adeia/pkg/util/stringutil"

	"github.com/jmoiron/sqlx"
)

var escapeDSNValue = stringutil.NewEscaper(`\`, `'`, `\`)

// PostgresDB represents an instance of the Postgres database connection.
type PostgresDB struct {
	*sqlx.DB
}

// New creates a new *PostgresDB.
func New(conf *config.DBConfig) (*PostgresDB, error) {
	d, err := sqlx.Connect(conf.Driver, buildDSN(conf))
	if err != nil {
		return nil, err
	}

	return &PostgresDB{d}, nil
}

// buildDSN is a helper that builds the DSN string.
func buildDSN(conf *config.DBConfig) (dsn string) {
	params := map[string]string{
		"dbname":   conf.DBName,
		"user":     conf.User,
		"password": conf.Password,
		"host":     conf.Host,
		"port":     strconv.Itoa(conf.Port),
		"sslmode":  conf.SSLMode,
	}
	for k, v := range params {
		dsn += k + "='" + escapeDSNValue(v) + "'"
	}

	// other sslParams are dont-cares when sslmode is "disable"
	// https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
	if params["sslmode"] == "disable" {
		return
	}

	sslParams := map[string]string{
		"sslcert":     conf.SSLCert,
		"sslkey":      conf.SSLKey,
		"sslrootcert": conf.SSLRootCert,
	}
	for k, v := range sslParams {
		dsn += k + "='" + escapeDSNValue(v) + "'"
	}
	return
}

// Insert inserts a row into the database. It returns the lastInsertID.
func (p *PostgresDB) Insert(ctx context.Context, query string, args ...interface{}) (lastInsertID int, err error) {
	err = p.QueryRowContext(ctx, query, args...).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// InsertNamed inserts a row into the database using a named query. It returns the
// lastInsertID. This method only accepts named queries.
func (p *PostgresDB) InsertNamed(ctx context.Context, namedQuery string, arg interface{}) (lastInsertID int, err error) {
	query, args, err := sqlx.Named(namedQuery, arg)
	if err != nil {
		return 0, err
	}

	query = p.Rebind(query)
	err = p.GetContext(ctx, &lastInsertID, query, args...)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// GetMany is a generic database SELECT that returns multiple records.
func (p *PostgresDB) GetMany(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	rows, err := p.QueryxContext(ctx, query, args...)
	if err != nil {
		return
	}
	defer ioutil.CheckCloseErr(rows, &err)

	if err = sqlx.StructScan(rows, dest); err != nil {
		return
	}
	return
}

// GetOne is a generic database SELECT that returns a single record.
func (p *PostgresDB) GetOne(ctx context.Context, dest interface{}, query string, args ...interface{}) (ok bool, err error) {
	if err := p.GetContext(ctx, dest, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Update is a generic database UPDATE that updates records.
func (p *PostgresDB) Update(ctx context.Context, query string, args ...interface{}) (rowsAffected int64, err error) {
	return p.exec(ctx, query, args...)
}

// UpdateNamed is a generic database UPDATE that updates records using a named query.
func (p *PostgresDB) UpdateNamed(ctx context.Context, query string, arg interface{}) (rowsAffected int64, err error) {
	return p.execNamed(ctx, query, arg)
}

// Delete is a generic database DELETE.
func (p *PostgresDB) Delete(ctx context.Context, query string, args ...interface{}) (rowsAffected int64, err error) {
	return p.exec(ctx, query, args...)
}

func (p *PostgresDB) exec(ctx context.Context, query string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := p.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (p *PostgresDB) execNamed(ctx context.Context, namedQuery string, arg interface{}) (rowsAffected int64, err error) {
	result, err := p.NamedExecContext(ctx, namedQuery, arg)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
