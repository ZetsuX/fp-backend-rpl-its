package dto

type AreaCreateRequest struct {
	Name       string `json:"name" binding:"required"`
	SpotCount  int    `json:"spot_count" binding:"required"`
	SpotPerRow int    `json:"spot_per_row" binding:"required"`
}