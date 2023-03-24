package entity

import "fp-rpl/common"

type Session struct {
	common.Model
	Time         string        `gorm:"type:time" json:"time" binding:"required"`
	Price        float64       `json:"price" binding:"required"`
	Transactions []Transaction `json:"transaction,omitempty"`
	Spots        []Spot        `json:"spot,omitempty"`
	FilmID       uint64        `gorm:"foreignKey" json:"film_id" binding:"required"`
	Film         *Film         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"film,omitempty"`
	AreaID       uint64        `gorm:"foreignKey" json:"area_id" binding:"required"`
	Area         *Area         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"area,omitempty"`
}
