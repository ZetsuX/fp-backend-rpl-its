package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type filmController struct {
	filmService service.FilmService
}

type FilmController interface {
	CreateFilm(ctx *gin.Context)
	GetAllFilms(ctx *gin.Context)
	GetFilmDetailBySlug(ctx *gin.Context)
	UpdateFilm(ctx *gin.Context)
	DeleteFilm(ctx *gin.Context)
	GetAllFilmsNowPlaying(ctx *gin.Context)
	GetAllFilmsComingSoon(ctx *gin.Context)
	GetAllFilmsNotPlaying(ctx *gin.Context)
}

func NewFilmController(filmS service.FilmService) FilmController {
	return &filmController{filmService: filmS}
}
func (fc *filmController) CreateFilm(ctx *gin.Context) {
	var filmDTO dto.FilmRegisterRequest
	err := ctx.ShouldBind(&filmDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process create film request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	filmCheck, er := fc.filmService.GetFilmBySlug(ctx, filmDTO.Slug)
	if er != nil {
		resp := common.CreateFailResponse("failed to process create film request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	if !(reflect.DeepEqual(filmCheck, entity.Film{})) {
		resp := common.CreateFailResponse("slug is not unique", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	filmDTO.Status = filmDTO.StatusCode.String()
	if filmDTO.Status == "Unknown" {
		resp := common.CreateFailResponse("status is invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	_, err = fc.filmService.CreateNewFilm(ctx, filmDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process create film request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateEmptySuccessResponse("film succesfully created", http.StatusCreated)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetAllFilms(ctx *gin.Context) {
	films, err := fc.filmService.GetAllFilm(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to get all film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("get film success", http.StatusCreated, films)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetAllFilmsNowPlaying(ctx *gin.Context) {
	films, err := fc.filmService.GetAllFilmByStatus(ctx,"Now Playing")
	if err != nil {
		resp := common.CreateFailResponse("failed to get all film that available", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("get film success", http.StatusCreated, films)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetAllFilmsComingSoon(ctx *gin.Context) {
	films, err := fc.filmService.GetAllFilmByStatus(ctx,"Coming Soon")
	if err != nil {
		resp := common.CreateFailResponse("failed to get all film that coming soon", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("get film success", http.StatusCreated, films)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetAllFilmsNotPlaying(ctx *gin.Context) {
	films, err := fc.filmService.GetAllFilmByStatus(ctx,"Not Playing")
	if err != nil {
		resp := common.CreateFailResponse("failed to get all film that not playing", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("get film success", http.StatusCreated, films)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetFilmDetailBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	film, err := fc.filmService.GetFilmDetailBySlug(ctx, slug)
	if err != nil {
		resp := common.CreateFailResponse("failed to get film detail", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(film, entity.Film{}) {
		resp := common.CreateFailResponse("film with given slug not found", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("get film detail success", http.StatusCreated, film)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) UpdateFilm(ctx *gin.Context) {
	slug := ctx.Param("slug")
	var filmDTO dto.FilmRegisterRequest
	err := ctx.ShouldBind(&filmDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to update film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if slug != filmDTO.Slug {
		filmCheck, err := fc.filmService.GetFilmBySlug(ctx,filmDTO.Slug)
		if err != nil {
			resp := common.CreateFailResponse("failed to verify update", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
		if !(reflect.DeepEqual(filmCheck, entity.Film{})) {
			resp := common.CreateFailResponse("slug is already used", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
	}

	film , err := fc.filmService.GetFilmBySlug(ctx,slug)
	if err != nil {
		resp := common.CreateFailResponse("failed to update film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	if reflect.DeepEqual(film,entity.Film{}) {
		resp := common.CreateFailResponse("film not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	filmDTO.Status = filmDTO.StatusCode.String()
	if filmDTO.Status == "Unknown" {
		resp := common.CreateFailResponse("status is invalid", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	updatedFilm, err := fc.filmService.UpdateFilm(ctx,filmDTO,film)
	if err != nil {
		resp := common.CreateFailResponse("failed to update film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("update film success", http.StatusCreated, updatedFilm)
	ctx.JSON(http.StatusCreated, resp)
}

func (fc *filmController) DeleteFilm(ctx *gin.Context) {
	slug := ctx.Param("slug")

	film, err := fc.filmService.GetFilmBySlug(ctx, slug)
	if err != nil {
		resp := common.CreateFailResponse("failed to process film delete request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(film, entity.Film{}) {
		resp := common.CreateFailResponse("film with given slug not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	err = fc.filmService.DeleteFilm(ctx, slug)
	if err != nil {
		resp := common.CreateFailResponse("failed to delete film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("delete film success", http.StatusCreated, nil)
	ctx.JSON(http.StatusCreated, resp)
}
