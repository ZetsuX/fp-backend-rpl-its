package controller

import (
	"fp-rpl/common"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/service"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

type areaController struct {
	areaService service.AreaService
	jwtService  service.JWTService
}

type AreaController interface {
	CreateArea(ctx *gin.Context)
	GetAllAreas(ctx *gin.Context)
	GetAreaByID(ctx *gin.Context)
	UpdateAreaByID(ctx *gin.Context)
	DeleteAreaByID(ctx *gin.Context)
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

func (areaC *areaController) GetAreaByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id of get area request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	area, err := areaC.areaService.GetAreaByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(area, entity.Area{}) {
		resp = common.CreateSuccessResponse("area not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("successfully fetched area", http.StatusOK, area)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (areaC *areaController) UpdateAreaByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id of get area request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var areaDTO dto.AreaCreateRequest
	err = ctx.ShouldBind(&areaDTO)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area update request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	area, err := areaC.areaService.GetAreaByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area update request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(area, entity.Area{}) {
		resp := common.CreateFailResponse("area with given id not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if (area.Name != areaDTO.Name) {
		// Check for duplicate Area Name
		areaCheck, err := areaC.areaService.GetAreaByName(ctx, areaDTO.Name)
		if err != nil {
			resp := common.CreateFailResponse("failed to process area update request", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Check if duplicate is found
		if !(reflect.DeepEqual(areaCheck, entity.Area{})) {
			resp := common.CreateFailResponse("name has already been used by another area", http.StatusBadRequest)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
	}

	area, err = areaC.areaService.UpdateArea(ctx, areaDTO, area)
	if err != nil {
		resp := common.CreateFailResponse(err.Error(), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	var resp common.Response
	if reflect.DeepEqual(area, entity.Area{}) {
		resp = common.CreateSuccessResponse("area not found", http.StatusOK, nil)
	} else {
		resp = common.CreateSuccessResponse("successfully updated area", http.StatusOK, area)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (areaC *areaController) DeleteAreaByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		resp := common.CreateFailResponse("failed to process id of delete area request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	area, err := areaC.areaService.GetAreaByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area delete request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if reflect.DeepEqual(area, entity.Area{}) {
		resp := common.CreateFailResponse("area with given id not found", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	err = areaC.areaService.DeleteAreaByID(ctx, id)
	if err != nil {
		resp := common.CreateFailResponse("failed to process area delete request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	resp := common.CreateSuccessResponse("successfully deleted area", http.StatusOK, nil)
	ctx.JSON(http.StatusOK, resp)
}