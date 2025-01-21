package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) query(db *sql.DB, sql string, params ...any) (*sql.Rows, error) {
	rows, err := db.Query(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return rows, nil
}

func (s *Postgres) exec(db *sql.DB, sql string, params ...any) (sql.Result, error) {
	result, err := db.Exec(sql, params...)
	if err != nil {
		sql = jdb.SQLParse(sql, params...)
		return nil, console.QueryError(err, sql)
	}

	return result, nil
}

func (s *Postgres) Exec(sql string, params ...any) error {
	_, err := s.exec(s.db, sql, params...)
	if err != nil {
		return err
	}

	return nil

}

func (s *Postgres) SQL(sql string, params ...any) (et.Items, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	var result = et.Items{Result: []et.Json{}}
	if s.typeSelect == jdb.Select {
		result = jdb.RowsToItems(rows)
	} else {
		result = jdb.SourceToItems(rows)
	}

	return result, nil
}

func (s *Postgres) One(sql string, params ...any) (et.Item, error) {
	rows, err := s.query(s.db, sql, params...)
	if err != nil {
		return et.Item{}, err
	}
	defer rows.Close()

	var result = et.Item{Result: et.Json{}}
	if s.typeSelect == jdb.Select {
		result = jdb.RowsToItem(rows)
	} else {
		result = jdb.SourceToItem(rows)
	}
	return result, nil
}

func (s *Postgres) Query(ql *jdb.Ql) (et.Items, error) {
	s.typeSelect = ql.TypeSelect
	ql.Sql = ""
	ql.Sql = strs.Append(ql.Sql, s.sqlSelect(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlFrom(ql.Froms), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlJoin(ql.Joins), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlWhere(ql.Wheres), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlGroupBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlHaving(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlOrderBy(ql), "\n")
	ql.Sql = strs.Append(ql.Sql, s.sqlLimit(ql), "\n")
	ql.Sql = strs.Format(`%s;`, ql.Sql)

	if ql.Show {
		console.Debug(ql.Sql)
	}

	result, err := s.SQL(ql.Sql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

func (s *Postgres) ExecDDL(id, sql string, params ...any) error {
	err := s.Exec(sql, params...)
	if err != nil {
		return err
	}

	go s.upsertDDL(id, sql)

	return nil
}
