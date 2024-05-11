package model

import (
	"github.com/uptrace/bun"
	"github.com/google/uuid"
)

type Groups struct {
	bun.BaseModel `bun:"table:groups"`
	Id uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Name string `bun:",type:varchar(255)"`
}