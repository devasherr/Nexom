package nexom

import (
	"context"
	"database/sql"
)

type M map[string]interface{}

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
	Where(conditions ...string) SelectThirdLevel
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
}

type SelectThirdLevel interface {
	Exec() (*sql.Rows, error)
	ExecContext(ctx context.Context) (*sql.Rows, error)
}

type DeleteSecondLevel interface {
	Where(conditions ...string) DeleteThirdLevel
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type DeleteThirdLevel interface {
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type DropLevel interface {
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type InsertSecondLevel interface {
	Values(fields ...string) InsertThirdLevel
}

type InsertThirdLevel interface {
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type UpdateSecondLevel interface {
	Set(M) UpdateThirdLevel
}

type UpdateThirdLevel interface {
	Where(fields ...string) UpdateFourthLevel
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}

type UpdateFourthLevel interface {
	Exec() (sql.Result, error)
	ExecContext(ctx context.Context) (sql.Result, error)
}
