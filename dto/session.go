package dto

type SessionCreateRequest struct {
	Time   string  `json:"time" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
	FilmID uint64  `json:"film_id" binding:"required"`
	AreaID uint64  `json:"area_id" binding:"required"`
}