package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type GroupPoets struct {
	bun.BaseModel `bun:"table:group_poets"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	GroupID uuid.UUID `bun:"group_id,notnull,type:uuid"`
	PoetID uuid.UUID `bun:"poet_id,notnull,type:uuid"`
	Group *Groups `bun:"group,rel:belongs-to,join:group_id=id"`
	Poet *Poets `bun:"poet,rel:belongs-to,join:poet_id=id"`
}