package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Relationships struct {
	bun.BaseModel `bun:"table:relationships"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	WhoID uuid.UUID `bun:"who_id,notnull,type:uuid"`
	WhomID uuid.UUID `bun:"whom_id,notnull,type:uuid"`

	Who *Users `bun:"who,rel:belongs-to,join:who_id=id"`
	Whom *Users `bun:"whom,rel:belongs-to,join:whom_id=id"`
}