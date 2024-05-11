package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type GroupPoets struct {
	bun.BaseModel `bun:"table:group_poets"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Group *Groups `bun:"rel:belongs-to"`
	Poet *Poets `bun:"rel:belongs-to"`
}