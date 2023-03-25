package dto

type FilmRegisterRequest struct {
	Title      string    `json:"title" binding:"required"`
	Slug       string    `json:"slug" binding:"required"`
	Synopsis   string    `json:"synopsis" binding:"required"`
	Duration   int       `json:"duration" binding:"required"`
	Genre      string    `json:"genre" binding:"required"`
	Producer   string    `json:"producer" binding:"required"`
	Director   string    `json:"director" binding:"required"`
	Writer     string    `json:"writer" binding:"required"`
	Production string    `json:"production" binding:"required"`
	Cast       string    `json:"cast" binding:"required"`
	Trailer    string    `json:"trailer" binding:"required"`
	Image      string    `json:"image" binding:"required"`
	Status     string    `json:"status" binding:"required"`
}
