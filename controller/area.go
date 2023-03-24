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
	GetAllAreas(ctx *gin.Context)
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
	if !(reflect.DeepEqual(areaCheck, entity.Area{})) {
		resp := common.CreateFailResponse("name has already been used by another area", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	newArea, err := areaC.areaService.CreateNewArea(ctx, areaDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area create request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully created area", http.StatusCreated, newArea)
	ctx.JSON(http.StatusCreated, resp)
}

func (areaC *areaController) GetAllAreas(ctx *gin.Context) {
	areas, err := areaC.areaService.GetAllAreas(ctx)
	if err != nil {
		resp := common.CreateFailResponse("failed to fetch all areas", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if len(areas) == 0 {
		resp = common.CreateSuccessResponse("no area found", http.StatusOK, areas)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched all areas", http.StatusOK, areas)
	}
	ctx.JSON(http.StatusOK, resp)
}