package service

import (
	"context"
	"fp-rpl/entity"
	"fp-rpl/repository"
)

type spotService struct {
	spotRepository repository.SpotRepository
}

type SpotService interface {
	GetSpotBySessionIDAndAttributes(ctx context.Context, sessionID uint64, row string, number int) (entity.Spot, error)
	UpdateSpot(ctx context.Context, spot entity.Spot) (entity.Spot, error)
}

func NewSpotService(spotR repository.SpotRepository) SpotService {
	return &spotService{spotRepository: spotR}
}

func (spotS *spotService) GetSpotBySessionIDAndAttributes(ctx context.Context, sessionID uint64, row string, number int) (entity.Spot, error) {
	spot, err := spotS.spotRepository.GetSpotBySessionIDAndAttributes(ctx, nil, sessionID, row, number)
	if err != nil {
		return entity.Spot{}, err
	}
	return spot, nil
}

func (spotS *spotService) UpdateSpot(ctx context.Context, spot entity.Spot) (entity.Spot, error) {
	spot, err := spotS.spotRepository.UpdateSpot(ctx, nil, spot)
	if err != nil {
		return entity.Spot{}, err
	}
	return spot, nil
}