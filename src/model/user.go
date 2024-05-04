package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Users struct {
	bun.BaseModel `bun:"table:users"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Name string `bun:"name,notnull,type:varchar(127)"`
	Bio string `bun:"bio,notnull,type:text"`
	Since time.Time `bun:"since,notnull,default:current_timestamp"`
}
