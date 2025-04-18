package main

import (
	"fmt"
	"strings"
)

type LevelOne interface {
	Select(fields ...string) LevelTwo
}

type LevelTwo interface {
	Where(condition string) LevelThree
	Exec()
}

type LevelThree interface {
	And(condition string) LevelThree
	Or(condition string) LevelThree
	Exec()
}

type QueryBuilder struct {
	tableName    string
	selectFields []string
	whereClauses string
	andClauses   []string
	orClauses    []string
}

func (q *QueryBuilder) execute() {
	// TODO; main logic is in here
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

	query := fmt.Sprintf("SELECT %s from %s %s %s %s", fields, q.tableName, q.whereClauses, andConditions, orConditions)
	fmt.Println(query)
}

// level 1
type Orm struct {
	qb *QueryBuilder
}

func NewOrm(tableName string) *Orm {
	return &Orm{qb: &QueryBuilder{tableName: tableName}}
}

func (o *Orm) Select(fields ...string) LevelTwo {
	o.qb.selectFields = fields
	return &l2{qb: o.qb}
}

// level 2
type l2 struct {
	qb *QueryBuilder
}

func (l *l2) Where(condition string) LevelThree {
	l.qb.whereClauses = "WHERE " + condition
	return &l3{qb: l.qb}
}

func (l *l2) Exec() {
	l.qb.execute()
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

func (l *l3) Exec() {
	l.qb.execute()
}

func main() {
	// Users := NewOrm("users")
	// Users.Select("id").Where("id > 1").Exec()

	// Books := NewOrm("books")
	// Books.Select().Where("age > 4").And("city = 'phill'").Or("city = 'every'").Exec()
}
