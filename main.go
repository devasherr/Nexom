package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Driver struct {
	db *sql.DB
}

func New(driver string, details string) *Driver {
	db, err := sql.Open(driver, details)
	if err != nil {
		panic(err)
	}
	return &Driver{db: db}
}

func (d *Driver) NewOrm(tableName string) *Orm {
	return &Orm{qb: &QueryBuilder{tableName: tableName, db: d.db}}
}

type LevelOne interface {
	Select(fields ...string) LevelTwo
	Delete() LevelTwo
	Drop() DropLevel
	Insert(columns ...string) InsertSecondLevel
}

type LevelTwo interface {
	Where(conditionKey, conditionValue string) LevelThree
	Exec() (*QueryResult, error)
}

type LevelThree interface {
	And(conditionKey, conditionValue string) LevelThree
	Or(conditionKey, conditionValue string) LevelThree
	Exec() (*QueryResult, error)
}

type DropLevel interface {
	Exec() (sql.Result, error)
}

type InsertSecondLevel interface {
	Values(fields ...string) InsertThirdLevel
}

type InsertThirdLevel interface {
	Exec() (sql.Result, error)
}

type QueryResult struct {
	Rows   *sql.Rows
	Result sql.Result
}

type QueryBuilder struct {
	db        *sql.DB
	queryType string

	tableName    string
	selectFields []string
	whereClauses [2]string
	andClauses   [][]string
	orClauses    [][]string
}

type InsertBuilder struct {
	db        *sql.DB
	tableName string
	columns   []string
	values    []string
}

func (q *QueryBuilder) handleSelect() (*QueryResult, error) {
	fields := "*"
	if len(q.selectFields) > 0 {
		fields = strings.Join(q.selectFields, ", ")
	}

	var andConditions strings.Builder
	if len(q.andClauses) > 0 {
		for i := range q.andClauses {
			qq := q.andClauses[i]
			andConditions.WriteString(strings.Join(qq[:1], ", "))
			andConditions.WriteString(" ")
		}
	}

	var orConditions strings.Builder
	if len(q.orClauses) > 0 {
		for i := range q.orClauses {
			qq := q.andClauses[i]
			orConditions.WriteString(strings.Join(qq[:1], ", "))
			orConditions.WriteString(" ")
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s", fields, q.tableName, q.whereClauses[0], andConditions.String(), orConditions.String())

	args := []any{}
	if q.whereClauses[1] != "" {
		args = append(args, q.whereClauses[1])
	}

	if len(q.andClauses) > 0 {
		for i := range q.andClauses {
			args = append(args, q.andClauses[i][1])
		}
	}

	if len(q.orClauses) > 0 {
		for i := range q.orClauses {
			args = append(args, q.orClauses[i][1])
		}
	}

	rows, err := q.db.Query(query, args...)
	return &QueryResult{Rows: rows}, err
}

func (q *QueryBuilder) handleDelete() (*QueryResult, error) {
	var andConditions strings.Builder
	if len(q.andClauses) > 0 {
		for i := range q.andClauses {
			qq := q.andClauses[i]
			andConditions.WriteString(strings.Join(qq[:1], ", "))
			andConditions.WriteString(" ")
		}
	}

	var orConditions strings.Builder
	if len(q.orClauses) > 0 {
		for i := range q.orClauses {
			qq := q.andClauses[i]
			orConditions.WriteString(strings.Join(qq[:1], ", "))
			orConditions.WriteString(" ")
		}
	}

	query := fmt.Sprintf("DELETE FROM %s %s %s %s", q.tableName, q.whereClauses[0], andConditions.String(), orConditions.String())

	args := []any{}
	if q.whereClauses[1] != "" {
		args = append(args, q.whereClauses[1])
	}

	if len(q.andClauses) > 0 {
		for i := range q.andClauses {
			args = append(args, q.andClauses[i][1])
		}
	}

	if len(q.orClauses) > 0 {
		for i := range q.orClauses {
			args = append(args, q.orClauses[i][1])
		}
	}

	result, err := q.db.Exec(query, args...)
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

// level 1
type Orm struct {
	qb *QueryBuilder
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

// level 2
type l2 struct {
	qb *QueryBuilder
}

func (l *l2) Where(conditionKey, conditionValue string) LevelThree {
	l.qb.whereClauses = [2]string{"WHERE " + conditionKey + " ?", conditionValue}
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

// level 3
type l3 struct {
	qb *QueryBuilder
}

func (l *l3) And(conditionKey, conditionValue string) LevelThree {
	and_c := []string{"AND " + conditionKey + " ?", conditionValue}
	l.qb.andClauses = append(l.qb.andClauses, and_c)
	return l
}

func (l *l3) Or(conditionKey, conditionValue string) LevelThree {
	or_c := []string{"OR " + conditionKey + " ?", conditionValue}
	l.qb.andClauses = append(l.qb.andClauses, or_c)
	return l
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

type insertThirdLevel struct {
	ib *InsertBuilder
}

func (it *insertThirdLevel) Exec() (sql.Result, error) {
	return it.ib.handleInsert()
}

func main() {
	norm := New("mysql", "root:1234@/income_expense")
	defer norm.db.Close()

	persons := norm.NewOrm("persons")
	// _, err := persons.Insert().Values("1", "bobby", "john", "nowhere", "mars").Exec()
	_, err := persons.Insert("id", "lastName", "firstName", "address", "city").Values("2", "bob", "john", "nowhere", "mars").Exec()
	if err != nil {
		panic(err)
	}
}
