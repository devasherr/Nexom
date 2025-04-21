package nexom

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Driver struct {
	DB *sql.DB
}

func New(driver string, details string) *Driver {
	db, err := sql.Open(driver, details)
	if err != nil {
		panic(err)
	}
	return &Driver{DB: db}
}

func (d *Driver) NewOrm(tableName string) *Orm {
	return &Orm{qb: &QueryBuilder{tableName: tableName, db: d.DB}}
}

type Orm struct {
	qb *QueryBuilder
}
