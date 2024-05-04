package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Reactions struct {
	bun.BaseModel `bun:"table:reactions"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Meigen uuid.UUID `bun:"meigen,notnull,type:uuid"`
	Whom uuid.UUID `bun:"whom,notnull,type:uuid"`
}