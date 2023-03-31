package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type sessionController struct {
	sessionService service.SessionService
	areaService    service.AreaService
	filmService    service.FilmService
}

type SessionController interface {
	CreateSession(ctx *gin.Context)
	GetAllSessions(ctx *gin.Context)
	GetSessionsByFilmSlug(ctx *gin.Context)
	DeleteSessionByID(ctx *gin.Context)
	GetSessionDetailByID(ctx *gin.Context)
}

func NewSessionController(sessionS service.SessionService, areaS service.AreaService, filmS service.FilmService) SessionController {
	return &sessionController{
		sessionService: sessionS,
		areaService:    areaS,
		filmService:    filmS,
	}
}

func (sessionC *sessionController) CreateSession(ctx *gin.Context) {
	var sessionDTO dto.SessionCreateRequest
	err := ctx.ShouldBind(&sessionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	checkTime, err := time.Parse(time.RFC3339, sessionDTO.Time)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session time", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if checkTime.Before(time.Now()) {
		resp := common.CreateFailResponse("invalid session time", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check for duplicate Session
	sessionCheck, err := sessionC.sessionService.GetSessionByTimeAndPlace(ctx, sessionDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check if duplicate is found
	if !(reflect.DeepEqual(sessionCheck, entity.Session{})) {
		resp := common.CreateFailResponse("session with the exact same attributes already exists", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check Film by ID
	film, err := sessionC.filmService.GetFilmByID(ctx, sessionDTO.FilmID)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(film, entity.Film{}) {
		resp := common.CreateFailResponse("film not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if film.Status != "Now Playing" {
		resp := common.CreateFailResponse("Film is not currently playing", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	// Check Area by ID
	area, err := sessionC.areaService.GetAreaByID(ctx, sessionDTO.AreaID)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(area, entity.Area{}) {
		resp := common.CreateFailResponse("area not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	_, err = sessionC.sessionService.CreateNewSession(ctx, sessionDTO, area.SpotCount, area.SpotPerRow)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateEmptySuccessResponse("successfully created session", http.StatusCreated)
	ctx.JSON(http.StatusCreated, resp)
}

func (sessionC *sessionController) GetAllSessions(ctx *gin.Context) {
	sessions, err := sessionC.sessionService.GetAllSessions(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch all sessions", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(sessions) == 0 {
		resp = common.CreateSuccessResponse("no session found", http.StatusOK, sessions)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched all sessions", http.StatusOK, sessions)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (sessionC *sessionController) GetSessionsByFilmSlug(ctx *gin.Context) {
	filmSlug := ctx.Param("filmslug")

	film, err := sessionC.filmService.GetFilmDetailBySlug(ctx, filmSlug)
	if err != nil {
		resp := common.CreateFailResponse("failed to get session", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(film, entity.Film{}) {
		resp := common.CreateFailResponse("film with given slug not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully fetched sessions", http.StatusOK, film.Sessions)
	ctx.JSON(http.StatusOK, resp)
}

func (sessionC *sessionController) DeleteSessionByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id of delete session request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	session, err := sessionC.sessionService.GetSessionByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session delete request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(session, entity.Session{}) {
		resp := common.CreateFailResponse("session with given id not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	err = sessionC.sessionService.DeleteSessionByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session delete request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully deleted session", http.StatusOK, nil)
	ctx.JSON(http.StatusOK, resp)
}

func (sessionC *sessionController) GetSessionDetailByID(ctx *gin.Context) {
	filmSlug := ctx.Param("filmslug")

	film, err := sessionC.filmService.GetFilmBySlug(ctx, filmSlug)
	if err != nil {
		resp := common.CreateFailResponse("failed to get film by filmslug of get session request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(film, entity.Film{}) {
		resp := common.CreateFailResponse("film with given slug not found", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id of get session request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	session, err := sessionC.sessionService.GetSessionByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to get session by id of get session request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(session, entity.Session{}) {
		resp := common.CreateFailResponse("session with given id not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if film.ID != session.FilmID {
		resp := common.CreateFailResponse("film and session not match", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	session, err = sessionC.sessionService.GetSessionDetailByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session get request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully fetched session", http.StatusOK, session)
	ctx.JSON(http.StatusOK, resp)
}
