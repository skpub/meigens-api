package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Reactions struct {
	bun.BaseModel `bun:"table:reactions"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	MeigenID uuid.UUID `bun:"meigen_id,notnull,type:uuid"`
	WhomID uuid.UUID `bun:"whom_id,notnull,type:uuid"`
	Meigen *Meigens `bun:"meigen,rel:belongs-to,join:meigen_id=id"`
	Whom *Users `bun:"user,rel:belongs-to,join:whom_id=id"`
}