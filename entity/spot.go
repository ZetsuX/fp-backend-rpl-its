package entity

import "fp-rpl/common"

type Spot struct {
	common.Model
	Row           rune         `gorm:"type:char;" json:"row" binding:"required"`
	Number        int          `json:"number" binding:"required"`
	SessionID     uint64       `gorm:"foreignKey" json:"session_id" binding:"required"`
	Session       *Session     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"session,omitempty"`
	TransactionID uint64       `gorm:"foreignKey" json:"transaction_id" binding:"required"`
	Transaction   *Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transaction,omitempty"`
}
