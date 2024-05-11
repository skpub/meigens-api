package model

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Users struct {
	bun.BaseModel `bun:"table:users"`
	Id uuid.UUID `bun:",pk,unique,type:uuid,default:uuid_generate_v4()"`
	Name string `bun:"name,unique,notnull,type:varchar(127)"`
	Bio string `bun:"bio,type:text"`
	Since time.Time `bun:"since,unique,notnull,default:current_timestamp"`
	Email string `bun:"email,unique,notnull,type:varchar(127)"`
	Password string `bun:"password,unique,notnull,type:char(64)"`
}
