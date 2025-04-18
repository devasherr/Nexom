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
}

type LevelTwo interface {
	Where(conditionKey, conditionValue string) LevelThree
	Exec() (*sql.Rows, error)
}

type LevelThree interface {
	And(condition string) LevelThree
	Or(condition string) LevelThree
	Exec() (*sql.Rows, error)
}

type QueryBuilder struct {
	db *sql.DB

	tableName    string
	selectFields []string
	whereClauses [2]string
	andClauses   []string
	orClauses    []string
}

func (q *QueryBuilder) execute() (*sql.Rows, error) {
	fields := "*"
	if len(q.selectFields) > 0 {
		fields = strings.Join(q.selectFields, ", ")
	}

	andConditions := ""
	if len(q.andClauses) > 0 {
		andConditions = strings.Join(q.andClauses, ", ")
	}

	orConditions := ""
	if len(q.orClauses) > 0 {
		orConditions = strings.Join(q.orClauses, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s", fields, q.tableName, q.whereClauses[0], andConditions, orConditions)

	args := []any{}
	if q.whereClauses[1] != "" {
		args = append(args, q.whereClauses[1])
	}

	rows, err := q.db.Query(query, args...)
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, nil
}

// level 1
type Orm struct {
	qb *QueryBuilder
}

func (o *Orm) Select(fields ...string) LevelTwo {
	o.qb.selectFields = fields
	return &l2{qb: o.qb}
}

// level 2
type l2 struct {
	qb *QueryBuilder
}

func (l *l2) Where(conditionKey, conditionValue string) LevelThree {
	l.qb.whereClauses = [2]string{"WHERE " + conditionKey + " ?", conditionValue}
	return &l3{qb: l.qb}
}

func (l *l2) Exec() (*sql.Rows, error) {
	return l.qb.execute()
}

// level 3
type l3 struct {
	qb *QueryBuilder
}

func (l *l3) And(condition string) LevelThree {
	l.qb.andClauses = append(l.qb.andClauses, "AND "+condition)
	return l
}

func (l *l3) Or(condition string) LevelThree {
	l.qb.orClauses = append(l.qb.orClauses, "OR "+condition)
	return l
}

func (l *l3) Exec() (*sql.Rows, error) {
	return l.qb.execute()
}

func main() {
	norm := New("mysql", "root:1234@/income_expense")

	Users := norm.NewOrm("users")
	rows, err := Users.Select().Where("username =", "Abebe").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)
}
