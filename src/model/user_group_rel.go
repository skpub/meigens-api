package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type UserGroupRels struct {
	bun.BaseModel `bun:"table:user_group_rels"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	GroupID uuid.UUID `bun:"group_id,notnull,type:uuid"`
	UserID uuid.UUID `bun:"user_id,notnull,type:uuid"`

	Group *Groups `bun:"group,rel:belongs-to,join:group_id=id"`
	User *Users `bun:"user,rel:belongs-to,join:user_id=id"`
}