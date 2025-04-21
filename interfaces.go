package nexom

import "database/sql"

type LevelOne interface {
	Select(fields ...string) LevelTwo
	Delete() LevelTwo
	Drop() DropLevel
	Insert(columns ...string) InsertSecondLevel
	Update() UpdateSecondLevel
}

type LevelTwo interface {
	Where(conditions ...string) LevelThree
	Exec() (*QueryResult, error)
}

type LevelThree interface {
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

type UpdateSecondLevel interface {
	Set(map[string]interface{}) UpdateThirdLevel
}

type UpdateThirdLevel interface {
	Where(fields ...string) UpdateFourthLevel
	Exec() (sql.Result, error)
}

type UpdateFourthLevel interface {
	Exec() (sql.Result, error)
}
