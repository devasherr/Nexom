package nexom

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type QueryResult struct {
	Rows   *sql.Rows
	Result sql.Result
}

func (o *Orm) Select(fields ...string) LevelTwo {
	o.qb.queryType = "select"
	o.qb.selectFields = fields
	return &l2{qb: o.qb}
}

func (o *Orm) Delete() LevelTwo {
	o.qb.queryType = "delete"
	return &l2{qb: o.qb}
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
		values:       map[string]interface{}{},
	}

	return &updateSecondLevel{ub: ub}
}

// level 2
type l2 struct {
	qb *QueryBuilder
}

func (l *l2) Where(conditions ...string) LevelThree {
	l.qb.whereClauses = conditions
	return &l3{qb: l.qb}
}

func (l *l2) Exec() (*QueryResult, error) {
	if l.qb.queryType == "select" {
		return l.qb.handleSelect()
	} else if l.qb.queryType == "delete" {
		return l.qb.handleDelete()
	} else {
		return &QueryResult{}, nil
	}
}

func (l *l2) ExecContext(ctx context.Context) (*QueryResult, error) {
	if l.qb.queryType == "select" {
		return l.qb.handleSelectWithContext(ctx)
	} else if l.qb.queryType == "delete" {
		return l.qb.handleDeleteWithContext(ctx)
	} else {
		return &QueryResult{}, nil
	}
}

type dropLevel struct {
	qb *QueryBuilder
}

func (d *dropLevel) Exec() (sql.Result, error) {
	return d.handleDrop()
}

type insertSecondLevel struct {
	ib *InsertBuilder
}

func (is *insertSecondLevel) Values(values ...string) InsertThirdLevel {
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

func (us *updateSecondLevel) Set(values map[string]interface{}) UpdateThirdLevel {
	us.ub.values = values
	return &updateThirdLevel{ub: us.ub}
}

// level 3
type l3 struct {
	qb *QueryBuilder
}

func (l *l3) Exec() (*QueryResult, error) {
	if l.qb.queryType == "select" {
		return l.qb.handleSelect()
	} else if l.qb.queryType == "delete" {
		return l.qb.handleDelete()
	} else {
		return &QueryResult{}, nil
	}
}

func (l *l3) ExecContext(ctx context.Context) (*QueryResult, error) {
	if l.qb.queryType == "select" {
		return l.qb.handleSelectWithContext(ctx)
	} else if l.qb.queryType == "delete" {
		return l.qb.handleDeleteWithContext(ctx)
	} else {
		return &QueryResult{}, nil
	}
}

type insertThirdLevel struct {
	ib *InsertBuilder
}

func (it *insertThirdLevel) Exec() (sql.Result, error) {
	return it.ib.handleInsert()
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

// level 4
type updateFourthLevel struct {
	ub *UpdateBuilder
}

func (uf *updateFourthLevel) Exec() (sql.Result, error) {
	return uf.ub.handleUpdate()
}

func (q *QueryBuilder) handleSelect() (*QueryResult, error) {
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

	rows, err := q.db.Query(query, args...)
	return &QueryResult{Rows: rows}, err
}

func (q *QueryBuilder) handleSelectWithContext(ctx context.Context) (*QueryResult, error) {
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

	rows, err := q.db.QueryContext(ctx, query, args...)
	return &QueryResult{Rows: rows}, err
}

func (q *QueryBuilder) handleDelete() (*QueryResult, error) {
	whereConditions := ""
	args := []any{}

	if len(q.whereClauses) > 0 {
		whereConditions = "WHERE " + q.whereClauses[0]
		for i := 1; i < len(q.whereClauses); i++ {
			args = append(args, q.whereClauses[i])
		}
	}

	query := fmt.Sprintf("DELETE FROM %s %s", q.tableName, whereConditions)

	result, err := q.db.Exec(query, args...)
	return &QueryResult{Result: result}, err
}

func (q *QueryBuilder) handleDeleteWithContext(ctx context.Context) (*QueryResult, error) {
	whereConditions := ""
	args := []any{}

	if len(q.whereClauses) > 0 {
		whereConditions = "WHERE " + q.whereClauses[0]
		for i := 1; i < len(q.whereClauses); i++ {
			args = append(args, q.whereClauses[i])
		}
	}

	query := fmt.Sprintf("DELETE FROM %s %s", q.tableName, whereConditions)

	result, err := q.db.ExecContext(ctx, query, args...)
	return &QueryResult{Result: result}, err
}
func (d *dropLevel) handleDrop() (sql.Result, error) {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", d.qb.tableName)
	return d.qb.db.Exec(query)
}

func (ib *InsertBuilder) handleInsert() (sql.Result, error) {
	insertColumns := ""
	if len(ib.columns) > 0 {
		insertColumns = "(" + strings.Join(ib.columns, ", ") + ")"
	}

	insertValues := ""
	protectedValues := make([]string, len(ib.values))
	for i := range ib.values {
		protectedValues[i] = "?"
	}

	if len(ib.values) > 0 {
		insertValues = "(" + strings.Join(protectedValues, ", ") + ")"
	}

	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", ib.tableName, insertColumns, insertValues)
	args := []any{}
	if len(ib.values) > 0 {
		for i := range ib.values {
			args = append(args, ib.values[i])
		}
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

	return ub.db.Exec(query, args...)
}
