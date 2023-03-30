package dto

type TransactionMakeRequest struct {
	Code       string
	TotalPrice float64
	SpotsName  []string `json:"spots_name" binding:"required"`
	UserID     uint64
	SessionID  uint64
}