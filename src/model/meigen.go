package model

import (
	"github.com/uptrace/bun"
)

type Meigen struct {
	bun.BaseModel `bun:"table:meigens,alias:m"`
	Id int64 `bun:"id,pk,autoincrement"`
	Meigen string `bun:"name,type:text"`
	date string `pg:"type:date"`
	whom int64 `bun:"id"`
}