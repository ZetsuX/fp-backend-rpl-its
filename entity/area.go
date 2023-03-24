package entity

import "fp-rpl/common"

type Area struct {
	common.Model
	Name       string    `json:"name" binding:"required"`
	SpotCount  int       `json:"spot_count" binding:"required"`
	SpotPerRow int       `json:"spot_per_row" binding:"required"`
	Sessions   []Session `json:"session,omitempty"`
}
