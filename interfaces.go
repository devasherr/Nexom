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
}

type SelectSecondLevel interface {
	Join(tableName string) SelectModifier
	JoinLeft(tableName string) SelectModifier
	JoinRight(tableName string) SelectModifier

	Where(conditions ...string) SelectThirdLevel
	Log() (string, []any)
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
}

type SelectModifier interface {
	On(conditions string) SelectSecondLevel
}

type SelectThirdLevel interface {
	Log() (string, []any)
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
}

type DeleteSecondLevel interface {
	Where(conditions ...string) DeleteThirdLevel
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
	Where(fields ...string) UpdateFourthLevel
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type UpdateFourthLevel interface {
	Log() (string, []any)
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}
