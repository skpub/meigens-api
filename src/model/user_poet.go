package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type UserPoets struct {
	bun.BaseModel `bun:"table:user_poets"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	User *Users `bun:"rel:belongs-to"`
	Poet *Poets `bun:"rel:belongs-to"`
}