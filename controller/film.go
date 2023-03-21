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
	jwtService  service.JWTService
}

type FilmController interface {
	CreateFilm(ctx *gin.Context)
}

func NewfilmController(filmS service.FilmService, jwtS service.JWTService) FilmController {
	return &filmController{
		filmService: filmS,
		jwtService: jwtS,
	}
}
func (fc *filmController) CreateFilm(ctx *gin.Context) {
	var filmDTO dto.FilmRegisterRequest
	err := ctx.ShouldBind(&filmDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process create film request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	filmCheck, er := fc.filmService.GetFilmBySlug(ctx,filmDTO.Slug)
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
	newFilm, err :=fc.filmService.CreateNewFilm(ctx,filmDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process create film request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("Film succesfully created", http.StatusCreated, newFilm)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) GetAllFilms(ctx *gin.Context) {
	films, err := fc.filmService.GetAllFilm(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to get all film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("Get film success", http.StatusCreated, films)
	ctx.JSON(http.StatusCreated, resp)
}
func (fc *filmController) UpdateFilm(ctx *gin.Context)  {
	slug := ctx.Param("slug")
	var filmDTO dto.FilmUpdateRequest
	err := ctx.ShouldBind(&filmDTO)
	if  err != nil {
		resp := common.CreateFailResponse("failed to update film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	updatedFilm,err := fc.filmService.UpdateFilm(ctx,filmDTO, slug)
	if err != nil{
		resp := common.CreateFailResponse("failed to update film", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	resp := common.CreateSuccessResponse("update film success", http.StatusCreated, updatedFilm)
	ctx.JSON(http.StatusCreated, resp)
}