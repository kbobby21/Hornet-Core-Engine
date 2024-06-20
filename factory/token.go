package factory

import "time"

type Tokens struct {
	Token      string    `json:"token"`
	IsAdmin    bool      `json:"is_admin"`
	Privileges string    `json:"privileges"`
	CreateTime time.Time `json:"creation_time"`
	ValidTill  time.Time `json:"valid_till"`
}
