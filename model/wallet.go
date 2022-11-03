package model

import "time"

type Wallet struct {
	Id         string    `json:"id"`
	OwnedBy    string    `json:"owned_by"`
	Status     bool      `json:"status"`
	EnabledAt  time.Time `json:"enabled_at"`
	DisabledAt time.Time `json:"disabled_at"`
	Balance    int       `json:"balance"`
}
