package service

import (
	"context"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/repository"

	"github.com/jinzhu/copier"
)

type areaService struct {
	areaRepository repository.AreaRepository
}

type AreaService interface {
	GetAreaByName(ctx context.Context, name string) (entity.Area, error)
	CreateNewArea(ctx context.Context, areaDTO dto.AreaCreateRequest) (entity.Area, error)
	GetAllAreas(ctx context.Context) ([]entity.Area, error)
	GetAreaByID(ctx context.Context, id uint64) (entity.Area, error)
	UpdateArea(ctx context.Context, areaDTO dto.AreaCreateRequest, area entity.Area) (entity.Area, error)
}

func NewAreaService(areaR repository.AreaRepository) AreaService {
	return &areaService{areaRepository: areaR}
}

func (areaS *areaService) GetAreaByName(ctx context.Context, name string) (entity.Area, error) {
	area, err := areaS.areaRepository.GetAreaByName(ctx, nil, name)
	if err != nil {
		return entity.Area{}, err
	}
	return area, nil
}

func (areaS *areaService) CreateNewArea(ctx context.Context, areaDTO dto.AreaCreateRequest) (entity.Area, error) {
	// Copy AreaDTO to empty newly created area var
	var area entity.Area
	copier.Copy(&area, &areaDTO)

	// create new area
	newArea, err := areaS.areaRepository.CreateNewArea(ctx, nil, area)
	if err != nil {
		return entity.Area{}, err
	}
	return newArea, nil
}

func (areaS *areaService) GetAllAreas(ctx context.Context) ([]entity.Area, error) {
	areas, err := areaS.areaRepository.GetAllAreas(ctx, nil)
	if err != nil {
		return []entity.Area{}, err
	}
	return areas, nil
}

func (areaS *areaService) GetAreaByID(ctx context.Context, id uint64) (entity.Area, error) {
	area, err := areaS.areaRepository.GetAreaByID(ctx, nil, id)
	if err != nil {
		return entity.Area{}, err
	}
	return area, nil
}

func (areaS *areaService) UpdateArea(ctx context.Context, areaDTO dto.AreaCreateRequest, area entity.Area) (entity.Area, error) {
	area, err := areaS.areaRepository.UpdateArea(ctx, nil, areaDTO, area)
	if err != nil {
		return entity.Area{}, err
	}
	return area, nil
}