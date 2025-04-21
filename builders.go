package nexom

import (
	"context"
	"database/sql"
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
	values    []string
}

type UpdateBuilder struct {
	db      *sql.DB
	context context.Context

	tableName    string
	whereClauses []string
	values       map[string]interface{}
}
