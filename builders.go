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
