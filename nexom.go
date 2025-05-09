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

func (o *Orm) ASC(value string) string {
	return "ASC " + value
}

func (o *Orm) DESC(value string) string {
	return "DESC " + value
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
		db:        o.qb.db,
		tableName: o.qb.tableName,
	}

	return &updateSecondLevel{ub: ub}
}

func (o *Orm) Prepare(query string) (*sql.Stmt, error) {
	return o.qb.db.Prepare(query)
}

func (o *Orm) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return o.qb.db.PrepareContext(ctx, query)
}

// level 2
type ss struct {
	qb *QueryBuilder
}

// select modifier (JOIN)
type sm struct {
	qb *QueryBuilder
}

func (s *sm) On(conditions string) SelectSecondLevel {
	s.qb.joinStatement += "ON " + conditions
	return &ss{qb: s.qb}
}

// select bound (LIMIT)
type sb struct {
	qb *QueryBuilder
}

func (s *sb) Exec() (*sql.Rows, error) {
	query, args := s.qb.SelectQuery()
	if s.qb.context != nil {
		return s.qb.db.QueryContext(s.qb.context, query, args...)
	}

	return s.qb.db.Query(query, args...)
}

func (s *sb) ExecContext(ctx context.Context) (*sql.Rows, error) {
	s.qb.context = ctx
	return s.Exec()
}

func (s *sb) Log() (string, []any) {
	return s.qb.SelectQuery()
}

func (s *ss) Join(tableName string) SelectModifier {
	s.qb.joinStatement = "JOIN " + tableName + " "
	return &sm{qb: s.qb}
}

func (s *ss) JoinLeft(tableName string) SelectModifier {
	s.qb.joinStatement = "LEFT JOIN " + tableName + " "
	return &sm{qb: s.qb}
}

func (s *ss) JoinRight(tableName string) SelectModifier {
	s.qb.joinStatement = "RIGHT JOIN " + tableName + " "
	return &sm{qb: s.qb}
}

func (s *ss) Where(condition string, args ...any) SelectThirdLevel {
	s.qb.whereClauses = condition
	s.qb.args = args
	return &st{qb: s.qb}
}

func (s *ss) Log() (string, []any) {
	return s.qb.SelectQuery()
}

func (s *ss) Limit(n int) SelectBounder {
	s.qb.limit = n
	return &sb{qb: s.qb}
}

func (s *ss) Exec() (*sql.Rows, error) {
	query, args := s.Log()
	if s.qb.context != nil {
		s.qb.db.QueryContext(s.qb.context, query, args...)
	}
	return s.qb.db.Query(query, args...)
}

func (s *ss) ExecContext(ctx context.Context) (*sql.Rows, error) {
	s.qb.context = ctx
	return s.Exec()
}

type ds struct {
	qb *QueryBuilder
}

func (d *ds) Where(condition string, args ...any) DeleteThirdLevel {
	d.qb.whereClauses = condition
	d.qb.args = args
	return &dt{qb: d.qb}
}

func (d *ds) Log() (string, []any) {
	return d.qb.DeleteQuery()
}

func (d *ds) Exec() (sql.Result, error) {
	query, args := d.Log()
	if d.qb.context != nil {
		return d.qb.db.ExecContext(d.qb.context, query, args...)
	}

	return d.qb.db.Exec(query, args...)
}

func (d *ds) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type dropLevel struct {
	qb *QueryBuilder
}

func (d *dropLevel) Log() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", d.qb.tableName)
}

func (d *dropLevel) Exec() (sql.Result, error) {
	query := d.Log()

	if d.qb.context != nil {
		return d.qb.db.ExecContext(d.qb.context, query)
	}
	return d.qb.db.Exec(query)
}

func (d *dropLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type insertSecondLevel struct {
	ib *InsertBuilder
}

func (is *insertSecondLevel) Values(values V) InsertThirdLevel {
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

func (s *st) Log() (string, []any) {
	return s.qb.SelectQuery()
}

func (s *st) Order(fields ...string) SelectFourthLevel {
	s.qb.orderFields = fields
	return &sf{qb: s.qb}
}

func (s *st) Limit(n int) SelectBounder {
	s.qb.limit = n
	return &sb{qb: s.qb}
}

func (s *st) Exec() (*sql.Rows, error) {
	query, args := s.Log()
	if s.qb.context != nil {
		s.qb.db.QueryContext(s.qb.context, query, args...)
	}
	return s.qb.db.Query(query, args...)
}

func (s *st) ExecContext(ctx context.Context) (*sql.Rows, error) {
	s.qb.context = ctx
	return s.Exec()
}

type dt struct {
	qb *QueryBuilder
}

func (d *dt) Log() (string, []any) {
	return d.qb.DeleteQuery()
}

func (d *dt) Exec() (sql.Result, error) {
	query, args := d.Log()

	if d.qb.context != nil {
		return d.qb.db.ExecContext(d.qb.context, query, args...)
	}

	return d.qb.db.Exec(query, args...)
}

func (d *dt) ExecContext(ctx context.Context) (sql.Result, error) {
	d.qb.context = ctx
	return d.Exec()
}

type insertThirdLevel struct {
	ib *InsertBuilder
}

func (it *insertThirdLevel) Log() (string, []any) {
	return it.ib.InsertQuery()
}

func (it *insertThirdLevel) Exec() (sql.Result, error) {
	query, args := it.Log()

	if it.ib.context != nil {
		return it.ib.db.ExecContext(it.ib.context, query, args...)
	}

	return it.ib.db.Exec(query, args...)
}

func (it *insertThirdLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	it.ib.context = ctx
	return it.Exec()
}

type updateThirdLevel struct {
	ub *UpdateBuilder
}

func (ut *updateThirdLevel) Where(condition string, args ...any) UpdateFourthLevel {
	ut.ub.whereClauses = condition
	ut.ub.args = args
	return &updateFourthLevel{ub: ut.ub}
}

func (ut *updateThirdLevel) Log() (string, []any) {
	return ut.ub.UpdateQuery()
}

func (ut *updateThirdLevel) Exec() (sql.Result, error) {
	query, args := ut.Log()

	if ut.ub.context != nil {
		return ut.ub.db.ExecContext(ut.ub.context, query, args...)
	}

	return ut.ub.db.Exec(query, args...)
}

func (ut *updateThirdLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	ut.ub.context = ctx
	return ut.Exec()
}

// level 4
type sf struct {
	qb *QueryBuilder
}

func (s *sf) Log() (string, []any) {
	return s.qb.SelectQuery()
}

func (s *sf) Limit(n int) SelectBounder {
	s.qb.limit = n
	return &sb{qb: s.qb}
}

func (s *sf) Exec() (*sql.Rows, error) {
	query, args := s.qb.SelectQuery()

	if s.qb.context != nil {
		return s.qb.db.QueryContext(s.qb.context, query, args...)
	}

	return s.qb.db.Query(query, args...)
}

func (s *sf) ExecContext(context context.Context) (*sql.Rows, error) {
	s.qb.context = context
	return s.Exec()
}

type updateFourthLevel struct {
	ub *UpdateBuilder
}

func (uf *updateFourthLevel) Log() (string, []any) {
	return uf.ub.UpdateQuery()
}

func (uf *updateFourthLevel) Exec() (sql.Result, error) {
	query, args := uf.Log()

	if uf.ub.context != nil {
		return uf.ub.db.ExecContext(uf.ub.context, query, args...)
	}

	return uf.ub.db.Exec(query, args...)
}

func (uf *updateFourthLevel) ExecContext(ctx context.Context) (sql.Result, error) {
	uf.ub.context = ctx
	return uf.Exec()
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
