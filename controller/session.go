package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type sessionController struct {
	sessionService service.SessionService
	areaService service.AreaService
	filmService service.FilmService
}

type SessionController interface {
	CreateSession(ctx *gin.Context)
	GetAllSessions(ctx *gin.Context)
}

func NewSessionController(sessionS service.SessionService, areaS service.AreaService, filmS service.FilmService) SessionController {
	return &sessionController{
		sessionService: sessionS,
		areaService: areaS,
		filmService: filmS,
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

	newSession, err := sessionC.sessionService.CreateNewSession(ctx, sessionDTO, area.SpotCount, area.SpotPerRow)
	if err != nil {
		resp := common.CreateFailResponse("failed to process session create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully created session", http.StatusCreated, newSession)
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
