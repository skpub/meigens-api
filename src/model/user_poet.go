package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type UserPoets struct {
	bun.BaseModel `bun:"table:user_poets"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID uuid.UUID `bun:"user_id,notnull,type:uuid"`
	PoetID uuid.UUID `bun:"poet_id,notnull,type:uuid"`
	User *Users `bun:"user,rel:belongs-to,join:user_id=id"`
	Poet *Poets `bun:"poet,rel:belongs-to,join:poet_id=id"`
}