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
	Whom uuid.UUID `bun:"whom,notnull,type:uuid"`
	Group *Groups `bun:"rel:belongs-to"`
	Poet *Poets `bun:"rel:belongs-to"`
}
