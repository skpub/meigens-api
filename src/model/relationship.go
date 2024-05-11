package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Relationships struct {
	bun.BaseModel `bun:"table:relationships"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Who *Users `bun:"rel:belongs-to"`
	Whom *Users `bun:"rel:belongs-to"`
}