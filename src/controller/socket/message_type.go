package socket

type MsgTLState struct {
	State int8 `json:"state"`
}

type MsgMeigen struct {
	Meigen string `json:"meigen"`
	Poet   string `json:"poet"`
}

type MsgMeigenGroup struct {
	Meigen  string `json:"meigen"`
	Poet    string `json:"poet"`
	GroupID string `json:"group_id"`
}
