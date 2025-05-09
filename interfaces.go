package nexom

import (
	"context"
	"database/sql"
)

type M map[string]interface{}
type V [][]string

type CreateSecondLevel interface {
	Values(values M) CreateThirdLevel
}

type CreateThirdLevel interface {
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type LevelOne interface {
	Select(fields ...string) SelectSecondLevel
	Delete() DeleteSecondLevel
	Drop() DropLevel
	Insert(columns ...string) InsertSecondLevel
	Update() UpdateSecondLevel

	Prepare() (*sql.Stmt, error)
}

type SelectSecondLevel interface {
	Join(tableName string) SelectModifier
	JoinLeft(tableName string) SelectModifier
	JoinRight(tableName string) SelectModifier

	Where(condition string, args ...any) SelectThirdLevel
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
	Log() (string, []any)
	Limit(n int) SelectBounder
}

type SelectModifier interface {
	On(conditions string) SelectSecondLevel
}

type SelectBounder interface {
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
	Log() (string, []any)
}

type SelectThirdLevel interface {
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
	Order(fields ...string) SelectFourthLevel
	Log() (string, []any)
	Limit(n int) SelectBounder
}

type SelectFourthLevel interface {
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
	Log() (string, []any)
	Limit(n int) SelectBounder
}

type DeleteSecondLevel interface {
	Where(condition string, args ...any) DeleteThirdLevel
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type DeleteThirdLevel interface {
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type DropLevel interface {
	Log() string
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type InsertSecondLevel interface {
	Values(values V) InsertThirdLevel
}

type InsertThirdLevel interface {
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type UpdateSecondLevel interface {
	Set(M) UpdateThirdLevel
}

type UpdateThirdLevel interface {
	Where(condition string, args ...any) UpdateFourthLevel
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type UpdateFourthLevel interface {
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}
