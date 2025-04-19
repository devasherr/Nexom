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
	And(conditionKey, conditionValue string) LevelThree
	Or(conditionKey, conditionValue string) LevelThree
	Exec() (*sql.Rows, error)
}

type QueryBuilder struct {
	db *sql.DB

	tableName    string
	selectFields []string
	whereClauses [2]string
	andClauses   [][]string
	orClauses    [][]string
}

func (q *QueryBuilder) execute() (*sql.Rows, error) {
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

func (l *l3) Exec() (*sql.Rows, error) {
	return l.qb.execute()
}

func main() {
	norm := New("mysql", "root:1234@/income_expense")

	income := norm.NewOrm("income")
	rows, err := income.Select().Where("price >", "400").Or("product_id =", "11").Or("income_id =", "10").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)
}
