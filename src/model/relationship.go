package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Relationships struct {
	bun.BaseModel `bun:"table:relationships"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Who uuid.UUID `bun:"who,notnull,type:uuid"`
	Whom uuid.UUID `bun:"whom,notnull,type:uuid"`
}