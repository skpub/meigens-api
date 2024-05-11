package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type UserGroupRels struct {
	bun.BaseModel `bun:"table:user_group_rels"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Group *Groups `bun:"rel:belongs-to"`
	User *Users `bun:"rel:belongs-to"`
}