package nexom

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (d *Driver) CreateTable(tableName string) CreateSecondLevel {
	cb := &CreateBuilder{
		db:        d.DB,
		tableName: tableName,
	}

	return &cs{cb: cb}
}

func (c *cs) Values(values M) CreateThirdLevel {
	c.cb.values = values
	return &ct{cb: c.cb}
}

type ct struct {
	cb *CreateBuilder
}

func (c *ct) Exec() (sql.Result, error) {
	return c.cb.handleCreateTable()
}

func (c *ct) ExecContext(ctx context.Context) (sql.Result, error) {
	c.cb.context = ctx
	return c.cb.handleCreateTable()
}

type cs struct {
	cb *CreateBuilder
}

func (o *Orm) Select(fields ...string) SelectSecondLevel {
	o.qb.selectFields = fields
	return &ss{qb: o.qb}
}

func (o *Orm) Delete() DeleteSecondLevel {
	return &ds{qb: o.qb}
}

func (o *Orm) Drop() DropLevel {
	return &dropLevel{qb: o.qb}
}

func (o *Orm) Insert(columns ...string) InsertSecondLevel {
	ib := &InsertBuilder{
		db:        o.qb.db,
		columns:   columns,
		tableName: o.qb.tableName,
	}

	return &insertSecondLevel{ib: ib}
}

func (o *Orm) Update() UpdateSecondLevel {
	ub := &UpdateBuilder{
		db:           o.qb.db,
		tableName:    o.qb.tableName,
		whereClauses: []string{},
		values:       M{},
	}

	return &updateSecondLevel{ub: ub}
}

// level 2
type ss struct {
	qb *QueryBuilder
}

func (s *ss) Where(conditions ...string) SelectThirdLevel {
	s.qb.whereClauses = conditions
	return &st{qb: s.qb}
}

func (s *ss) Exec() (*sql.Rows, error) {
	return s.qb.handleSelect()
}

func (s *ss) ExecContext(ctx context.Context) (*sql.Rows, error) {
	s.qb.context = ctx
	return s.Exec()
}

type ds struct {
	qb *QueryBuilder
}

func (d *ds) Where(conditions ...string) DeleteThirdLevel {
	d.qb.whereClauses = conditions
	return &dt{qb: d.qb}
}

func (d *ds) Exec() (sql.Result, error) {
	return d.qb.handleDelete()
}

func (d *ds) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type dropLevel struct {
	qb *QueryBuilder
}

func (d *dropLevel) Exec() (sql.Result, error) {
	return d.handleDrop()
}

func (d *dropLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type insertSecondLevel struct {
	ib *InsertBuilder
}

func (is *insertSecondLevel) Values(values [][]string) InsertThirdLevel {
	ib := &InsertBuilder{
		db:        is.ib.db,
		tableName: is.ib.tableName,
		columns:   is.ib.columns,
		values:    values,
	}

	return &insertThirdLevel{ib: ib}
}

type updateSecondLevel struct {
	ub *UpdateBuilder
}

func (us *updateSecondLevel) Set(values M) UpdateThirdLevel {
	us.ub.values = values
	return &updateThirdLevel{ub: us.ub}
}

// level 3
type st struct {
	qb *QueryBuilder
}

func (s *st) Exec() (*sql.Rows, error) {
	return s.qb.handleSelect()
}

func (s *st) ExecContext(ctx context.Context) (*sql.Rows, error) {
	s.qb.context = ctx
	return s.Exec()
}

type dt struct {
	qb *QueryBuilder
}

func (d *dt) Exec() (sql.Result, error) {
	return d.qb.handleDelete()
}

func (d *dt) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type insertThirdLevel struct {
	ib *InsertBuilder
}

func (it *insertThirdLevel) Exec() (sql.Result, error) {
	return it.ib.handleInsert()
}

func (it *insertThirdLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	it.ib.context = ctx
	return it.Exec()
}

type updateThirdLevel struct {
	ub *UpdateBuilder
}

func (ut *updateThirdLevel) Where(fields ...string) UpdateFourthLevel {
	ut.ub.whereClauses = fields
	return &updateFourthLevel{ub: ut.ub}
}

func (ut *updateThirdLevel) Exec() (sql.Result, error) {
	return ut.ub.handleUpdate()
}

func (ut *updateThirdLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	ut.ub.context = ctx
	return ut.Exec()
}

// level 4
type updateFourthLevel struct {
	ub *UpdateBuilder
}

func (uf *updateFourthLevel) Exec() (sql.Result, error) {
	return uf.ub.handleUpdate()
}

func (uf *updateFourthLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	uf.ub.context = ctx
	return uf.Exec()
}

func (q *QueryBuilder) handleSelect() (*sql.Rows, error) {
	fields := "*"
	if len(q.selectFields) > 0 {
		fields = strings.Join(q.selectFields, ", ")
	}

	whereConditions := ""
	args := []any{}

	if len(q.whereClauses) > 0 {
		whereConditions = "WHERE " + q.whereClauses[0]
		for i := 1; i < len(q.whereClauses); i++ {
			args = append(args, q.whereClauses[i])
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s", fields, q.tableName, whereConditions)

	if q.context != nil {
		return q.db.QueryContext(q.context, query, args...)
	}

	return q.db.Query(query, args...)
}

func (q *QueryBuilder) handleDelete() (sql.Result, error) {
	whereConditions := ""
	args := []any{}

	if len(q.whereClauses) > 0 {
		whereConditions = "WHERE " + q.whereClauses[0]
		for i := 1; i < len(q.whereClauses); i++ {
			args = append(args, q.whereClauses[i])
		}
	}

	query := fmt.Sprintf("DELETE FROM %s %s", q.tableName, whereConditions)

	if q.context != nil {
		return q.db.ExecContext(q.context, query, args...)
	}

	return q.db.Exec(query, args...)
}

func (d *dropLevel) handleDrop() (sql.Result, error) {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", d.qb.tableName)
	if d.qb.context != nil {
		return d.qb.db.ExecContext(d.qb.context, query)
	}

	return d.qb.db.Exec(query)
}

func (ib *InsertBuilder) handleInsert() (sql.Result, error) {
	insertColumns := ""
	if len(ib.columns) > 0 {
		insertColumns = "(" + strings.Join(ib.columns, ", ") + ")"
	}

	var insertValues strings.Builder
	for i := range len(ib.values) {
		curVal := []string{}
		for range len(ib.values[i]) {
			curVal = append(curVal, "?")
		}

		insertValues.WriteString("(" + strings.Join(curVal, ", ") + "), ")
	}

	// users should make sure of this, but
	// helps index out of bound error when insertValues is empty
	if insertValues.Len() == 0 {
		insertValues.WriteString("  ")
	}

	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", ib.tableName, insertColumns, insertValues.String()[:insertValues.Len()-2])
	args := []any{}
	if len(ib.values) > 0 {
		for i := range ib.values {
			for j := range ib.values[i] {
				args = append(args, ib.values[i][j])
			}
		}
	}

	if ib.context != nil {
		return ib.db.ExecContext(ib.context, query, args...)
	}

	return ib.db.Exec(query, args...)
}

func (ub *UpdateBuilder) handleUpdate() (sql.Result, error) {
	var values strings.Builder
	args := []any{}
	for key, val := range ub.values {
		values.WriteString(key + " = ?, ")
		args = append(args, val)
	}

	whereClauses := ""
	for i := range ub.whereClauses {
		if i == 0 {
			whereClauses = ub.whereClauses[i]
			continue
		}

		args = append(args, ub.whereClauses[i])
	}

	setValues := values.String()
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", ub.tableName, setValues[:len(setValues)-2], whereClauses)

	if ub.context != nil {
		return ub.db.ExecContext(ub.context, query, args...)
	}

	return ub.db.Exec(query, args...)
}

func (cb *CreateBuilder) handleCreateTable() (sql.Result, error) {
	var sb strings.Builder
	for key := range cb.values {
		sb.WriteString(fmt.Sprintf("%s %s, ", key, cb.values[key]))
	}

	tableDefinition := sb.String()
	if len(tableDefinition) > 0 {
		tableDefinition = "(" + tableDefinition[:len(tableDefinition)-2] + ")"
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", cb.tableName, tableDefinition)

	if cb.context != nil {
		return cb.db.ExecContext(cb.context, query)
	}

	return cb.db.Exec(query)
}
