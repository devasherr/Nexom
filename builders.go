package nexom

import (
	"database/sql"
)

type QueryBuilder struct {
	db        *sql.DB
	queryType string

	tableName    string
	selectFields []string
	whereClauses []string
}

type InsertBuilder struct {
	db        *sql.DB
	tableName string
	columns   []string
	values    []string
}

type UpdateBuilder struct {
	db           *sql.DB
	tableName    string
	whereClauses []string
	values       map[string]interface{}
}
