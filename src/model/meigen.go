package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Meigens struct {
	bun.BaseModel `bun:"table:meigens"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Meigen string `bun:"meigen,notnull,type:text"`
	CreatedAt time.Time `bun:"createdAt,notnull,default:current_timestamp"`
	WhomID uuid.UUID `bun:"whom_id,notnull,type:uuid"`
	GroupID uuid.UUID `bun:"group_id,type:uuid"`
	PoetID uuid.UUID `bun:"poet_id,notnull,type:uuid"`

	Whom *Users `bun:"rel:belongs-to,join:whom_id=id"`
	Group *Groups `bun:"rel:belongs-to,join:group_id=id"`
	Poet *Poets `bun:"rel:belongs-to,join:poet_id=id"`
}
