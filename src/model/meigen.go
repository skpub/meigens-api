package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Meigen struct {
	bun.BaseModel `bun:"table:meigens,alias:m"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Meigen string `bun:"meigen,notnull,type:text"`
	Date time.Time `bun:"createdAt,notnull,default:current_timestamp"`
	Whom uuid.UUID `bun:"whom,notnull,type:uuid"`
}
