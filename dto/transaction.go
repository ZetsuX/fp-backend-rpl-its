package dto

type TransactionMakeRequest struct {
	Code       string
	TotalPrice float64
	Timestamp  string   `json:"timestamp" binding:"required"`
	SpotsName  []string `json:"spots_name" binding:"required"`
	UserID     uint64
	SessionID  uint64
}