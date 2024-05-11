package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Reactions struct {
	bun.BaseModel `bun:"table:reactions"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Meigen *Meigens `bun:"rel:belongs-to"`
	Whom *Users `bun:"rel:belongs-to"`
}