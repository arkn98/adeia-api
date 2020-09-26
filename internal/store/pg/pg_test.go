package pg

import (
	"context"
	"errors"
	"testing"

	"adeia/internal/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (p *PostgresDB, mock sqlmock.Sqlmock, closer func()) {
	t.Parallel()
	mockDB, mock, _ := sqlmock.New()
	closer = func() {
		_ = mockDB.Close()
	}
	return &PostgresDB{sqlx.NewDb(mockDB, "sqlmock")}, mock, closer
}

func TestBuildDSN(t *testing.T) {
	testCases := []struct {
		in   *config.DBConfig
		want []string
		msg  string
	}{
		{
			in: &config.DBConfig{
				DBName:   "test",
				User:     "root",
				Password: "password",
				Host:     "localhost",
				Port:     5432,
				// we only have to worry passing the SSLParams when SSLMode is 'disable',
				// all other errors are returned when sqlx.Connect() is called
				SSLMode:     "foobar",
				SSLCert:     "ssl-cert",
				SSLKey:      "ssl-key",
				SSLRootCert: "ssl-root-cert",
			},
			want: []string{
				`dbname='test'`,
				`user='root'`,
				`password='password'`,
				`host='localhost'`,
				`port='5432'`,
				`sslmode='foobar'`,
				`sslcert='ssl-cert'`,
				`sslkey='ssl-key'`,
				`sslrootcert='ssl-root-cert'`,
			},
			msg: "return DSN with SSLParams",
		},
		{
			in: &config.DBConfig{
				DBName:   "test",
				User:     "root",
				Password: "password",
				Host:     "localhost",
				Port:     5432,
				SSLMode:  "disable",
			},
			want: []string{
				`dbname='test'`,
				`user='root'`,
				`password='password'`,
				`host='localhost'`,
				`port='5432'`,
				`sslmode='disable'`,
			},
			msg: "return DSN without SSLParams",
		},
		{
			in: &config.DBConfig{
				DBName:   `test`,
				User:     `ro''ot`,
				Password: `'password`,
				Host:     `localhost'`,
				Port:     5432,
				SSLMode:  `disable`,
			},
			want: []string{
				`dbname='test'`,
				`user='ro\'\'ot'`,
				`password='\'password'`,
				`host='localhost\'`,
				`port='5432'`,
				`sslmode='disable'`,
			},
			msg: "escape single quotes",
		},
		{
			in: &config.DBConfig{
				DBName:   `test`,
				User:     `ro\\ot`,
				Password: `\password`,
				Host:     `localhost\`,
				Port:     5432,
				SSLMode:  `disable`,
			},
			want: []string{
				`dbname='test'`,
				`user='ro\\\\ot'`,
				`password='\\password`,
				`host='localhost\\`,
				`port='5432'`,
				`sslmode='disable'`,
			},
			msg: "escape backslashes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			t.Parallel()
			got := buildDSN(tc.in)
			for _, w := range tc.want {
				assert.Contains(t, got, w)
			}
		})
	}
}

func TestPostgresDB_Delete(t *testing.T) {
	t.Run("return error when exec fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "DELETE FROM table WHERE col1=$1 AND col2=$2"
		mock.
			ExpectExec("DELETE FROM table WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2").
			WillReturnError(errors.New("exec failed"))
		_, err := p.Delete(context.Background(), query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("return error when result fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "DELETE FROM table WHERE col1=$1 AND col2=$2"
		mock.
			ExpectExec("DELETE FROM table WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2").
			WillReturnResult(sqlmock.NewErrorResult(errors.New("result failed")))
		_, err := p.Delete(context.Background(), query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("delete on no error", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "DELETE FROM table WHERE col1=$1 AND col2=$2"
		mock.
			ExpectExec("DELETE FROM table WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2").
			WillReturnResult(sqlmock.NewResult(0, 2))
		_, err := p.Delete(context.Background(), query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
	})
}

func TestPostgresDB_Update(t *testing.T) {
	t.Run("return error when exec fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "UPDATE table SET col=$1 WHERE col1=$2 AND col2=$3"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2", "arg3").
			WillReturnError(errors.New("exec failed"))
		_, err := p.Update(context.Background(), query, "arg1", "arg2", "arg3")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("return error when result fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "UPDATE table SET col=$1 WHERE col1=$2 AND col2=$3"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2", "arg3").
			WillReturnResult(sqlmock.NewErrorResult(errors.New("result failed")))
		_, err := p.Update(context.Background(), query, "arg1", "arg2", "arg3")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("update on no error", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "UPDATE table SET col=$1 WHERE col1=$2 AND col2=$3"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs("arg1", "arg2", "arg3").
			WillReturnResult(sqlmock.NewResult(0, 2))
		_, err := p.Update(context.Background(), query, "arg1", "arg2", "arg3")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
	})
}

func TestPostgresDB_UpdateNamed(t *testing.T) {
	t.Run("return error when exec fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		data := &struct {
			Col  string `db:"col"`
			Col1 string `db:"col1"`
			Col2 string `db:"col2"`
		}{
			"foo",
			"bar",
			"foobar",
		}
		query := "UPDATE table SET col=:col WHERE col1=:col1 AND col2=:col2"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs(data.Col, data.Col1, data.Col2).
			WillReturnError(errors.New("exec failed"))
		_, err := p.UpdateNamed(context.Background(), query, data)

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("return error when result fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		data := &struct {
			Col  string `db:"col"`
			Col1 string `db:"col1"`
			Col2 string `db:"col2"`
		}{
			"foo",
			"bar",
			"foobar",
		}
		query := "UPDATE table SET col=:col WHERE col1=:col1 AND col2=:col2"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs(data.Col, data.Col1, data.Col2).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("result failed")))
		_, err := p.UpdateNamed(context.Background(), query, data)

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("update on no error", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		data := &struct {
			Col  string `db:"col"`
			Col1 string `db:"col1"`
			Col2 string `db:"col2"`
		}{
			"foo",
			"bar",
			"foobar",
		}
		query := "UPDATE table SET col=:col WHERE col1=:col1 AND col2=:col2"
		mock.
			ExpectExec("UPDATE table SET col=(.+) WHERE col1=(.+) AND col2=(.+)").
			WithArgs(data.Col, data.Col1, data.Col2).
			WillReturnResult(sqlmock.NewResult(0, 1))
		_, err := p.UpdateNamed(context.Background(), query, data)

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
	})
}

func TestPostgresDB_GetOne(t *testing.T) {
	t.Run("return ok=false when no rows", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnRows(sqlmock.NewRows([]string{"foo", "bar"})).
			RowsWillBeClosed()

		type data struct {
			Foo string `db:"foo"`
			Bar int    `db:"bar"`
		}
		var got data
		ok, err := p.GetOne(context.Background(), &got, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.False(t, ok)
		assert.Empty(t, got)
	})

	t.Run("return err when query fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnError(errors.New("query failed"))
		ok, err := p.GetOne(context.Background(), &struct{}{}, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("on success return ok and fill struct", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnRows(sqlmock.NewRows([]string{"foo", "bar"}).AddRow("val1", 1)).
			RowsWillBeClosed()

		type data struct {
			Foo string `db:"foo"`
			Bar int    `db:"bar"`
		}
		var got data
		want := data{"val1", 1}
		ok, err := p.GetOne(context.Background(), &got, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}

func TestPostgresDB_GetMany(t *testing.T) {
	t.Run("return err when query fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnError(errors.New("query failed"))
		err := p.GetMany(context.Background(), []interface{}{}, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("return err when StructScan fails", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		rows := sqlmock.NewRows([]string{"foo", "bar"}).
			AddRow("foo1", 1).
			AddRow("foo2", 2).
			AddRow("foo3", 3)
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnRows(rows).
			RowsWillBeClosed()

		type data struct {
			Foo string `db:"foo"`
			Bar string `db:"bar"`
		}
		var got []data
		err := p.GetMany(context.Background(), got, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("empty slice when no rows", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnRows(sqlmock.NewRows([]string{"foo", "bar"})).
			RowsWillBeClosed()

		type data struct {
			Foo string `db:"foo"`
			Bar int    `db:"bar"`
		}
		var got []data
		err := p.GetMany(context.Background(), &got, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Nil(t, got)
	})

	t.Run("on success fill slice", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "SELECT * FROM users WHERE col1=$1 AND col2=$2"
		rows := sqlmock.NewRows([]string{"foo", "bar"}).
			AddRow("foo1", 1).
			AddRow("foo2", 2).
			AddRow("foo3", 3)
		mock.
			ExpectQuery(`SELECT \* FROM users WHERE col1=(.+) AND col2=(.+)`).
			WithArgs("arg1", "arg2").
			WillReturnRows(rows).
			RowsWillBeClosed()

		type data struct {
			Foo string `db:"foo"`
			Bar int    `db:"bar"`
		}
		var got []data
		want := []data{
			{"foo1", 1},
			{"foo2", 2},
			{"foo3", 3},
		}
		err := p.GetMany(context.Background(), &got, query, "arg1", "arg2")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

func TestPostgresDB_Insert(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "INSERT INTO users (col1, col2, col3) VALUES($1, $2, $3) RETURNING id"
		mock.
			ExpectQuery(`INSERT INTO users \(col1, col2, col3\) VALUES\((.+), (.+), (.+)\) RETURNING id`).
			WithArgs("arg1", "arg2", "arg3").
			WillReturnError(errors.New("insert failed"))
		_, err := p.Insert(context.Background(), query, "arg1", "arg2", "arg3")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("successful insert", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "INSERT INTO users (col1, col2, col3) VALUES($1, $2, $3) RETURNING id"
		want := 1
		mock.
			ExpectQuery(`INSERT INTO users \(col1, col2, col3\) VALUES\((.+), (.+), (.+)\) RETURNING id`).
			WithArgs("arg1", "arg2", "arg3").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(want)).
			RowsWillBeClosed()
		got, err := p.Insert(context.Background(), query, "arg1", "arg2", "arg3")

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}

func TestPostgresDB_InsertNamed(t *testing.T) {
	type data struct {
		Col1 string `db:"col1"`
		Col2 int    `db:"col2"`
		Col3 bool   `db:"col3"`
	}

	t.Run("named query error", func(t *testing.T) {
		p, _, c := setup(t)
		defer c()

		query := "INSERT INTO users (col1, col2, col3) VALUES(:col1, :co:l2, :col3) RETURNING id"
		d := &data{"foo", 1, true}
		_, err := p.InsertNamed(context.Background(), query, d)

		assert.Error(t, err)
	})

	t.Run("error when getting lastInsertID", func(t *testing.T) {
		p, _, c := setup(t)
		defer c()

		query := "INSERT INTO users (col1, col2, col3) VALUES(:col1, :col2, :col3)"
		d := &data{"foo", 1, true}
		_, err := p.InsertNamed(context.Background(), query, d)

		assert.Error(t, err)
	})

	t.Run("successful insert", func(t *testing.T) {
		p, mock, c := setup(t)
		defer c()

		query := "INSERT INTO users (col1, col2, col3) VALUES(:col1, :col2, :col3) RETURNING id"
		want := 1
		d := &data{"foo", 1, true}
		mock.
			ExpectQuery(`INSERT INTO users \(col1, col2, col3\) VALUES\((.+), (.+), (.+)\) RETURNING id`).
			WithArgs(d.Col1, d.Col2, d.Col3).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(want)).
			RowsWillBeClosed()
		got, err := p.InsertNamed(context.Background(), query, d)

		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})
}
