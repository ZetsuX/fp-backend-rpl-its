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

type areaController struct {
	areaService service.AreaService
	jwtService  service.JWTService
}

type AreaController interface {
	CreateArea(ctx *gin.Context)
}

func NewAreaController(areaS service.AreaService, jwtS service.JWTService) AreaController {
	return &areaController{
		areaService: areaS,
		jwtService:  jwtS,
	}
}

func (areaC *areaController) CreateArea(ctx *gin.Context) {
	var areaDTO dto.AreaCreateRequest
	err := ctx.ShouldBind(&areaDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check for duplicate Area Name
	areaCheck, err := areaC.areaService.GetAreaByName(ctx, areaDTO.Name)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	// Check if duplicate is found
	if !(reflect.DeepEqual(areaCheck, entity.User{})) {
		resp := common.CreateFailResponse("name has already been used by another area", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newUser, err := areaC.areaService.CreateNewArea(ctx, areaDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully created area", http.StatusCreated, newUser)
	ctx.JSON(http.StatusCreated, resp)
}