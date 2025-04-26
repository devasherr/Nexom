package nexom

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	db      *sql.DB
	context context.Context

	tableName    string
	selectFields []string
	whereClauses []string
}

type InsertBuilder struct {
	db      *sql.DB
	context context.Context

	tableName string
	columns   []string
	values    [][]string
}

type UpdateBuilder struct {
	db      *sql.DB
	context context.Context

	tableName    string
	whereClauses []string
	values       M
}

type CreateBuilder struct {
	db      *sql.DB
	context context.Context

	tableName string
	values    M
}

func (qb *QueryBuilder) SelectQuery() (string, []any) {
	fields := "*"
	if len(qb.selectFields) > 0 {
		fields = strings.Join(qb.selectFields, ", ")
	}

	whereConditions := ""
	args := []any{}

	if len(qb.whereClauses) > 0 {
		whereConditions = "WHERE " + qb.whereClauses[0]
		for i := 1; i < len(qb.whereClauses); i++ {
			args = append(args, qb.whereClauses[i])
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s", fields, qb.tableName, whereConditions)
	return query, args
}

func (qb *QueryBuilder) DeleteQuery() (string, []any) {
	whereConditions := ""
	args := []any{}

	if len(qb.whereClauses) > 0 {
		whereConditions = "WHERE " + qb.whereClauses[0]
		for i := 1; i < len(qb.whereClauses); i++ {
			args = append(args, qb.whereClauses[i])
		}
	}

	query := fmt.Sprintf("DELETE FROM %s %s", qb.tableName, whereConditions)
	return query, args
}

func (ib *InsertBuilder) InsertQuery() (string, []any) {
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

	return query, args
}

func (ub *UpdateBuilder) UpdateQuery() (string, []any) {
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

	return query, args
}
