package service

import (
	"context"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/repository"
	"fp-rpl/utils"

	"github.com/jinzhu/copier"
)

type sessionService struct {
	sessionRepository repository.SessionRepository
	spotRepository    repository.SpotRepository
}

type SessionService interface {
	GetSessionByTimeAndPlace(ctx context.Context, sessionDTO dto.SessionCreateRequest) (entity.Session, error)
	GetSessionByID(ctx context.Context, id uint64) (entity.Session, error)
	CreateNewSession(ctx context.Context, sessionDTO dto.SessionCreateRequest, spotCount int, spotPerRow int) (entity.Session, error)
	GetAllSessions(ctx context.Context) ([]entity.Session, error)
	DeleteSessionByID(ctx context.Context, id uint64) error
	GetSessionDetailByID(ctx context.Context, id uint64) (entity.Session, error)
}

func NewSessionService(sessionR repository.SessionRepository, spotR repository.SpotRepository) SessionService {
	return &sessionService{
		sessionRepository: sessionR,
		spotRepository:    spotR,
	}
}

func (sessionS *sessionService) GetSessionByTimeAndPlace(ctx context.Context, sessionDTO dto.SessionCreateRequest) (entity.Session, error) {
	session, err := sessionS.sessionRepository.GetSessionByTimeAndAreaID(ctx, nil, sessionDTO.Time, sessionDTO.AreaID)
	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (sessionS *sessionService) GetSessionByID(ctx context.Context, id uint64) (entity.Session, error) {
	session, err := sessionS.sessionRepository.GetSessionByID(ctx, nil, id)
	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (sessionS *sessionService) CreateNewSession(ctx context.Context, sessionDTO dto.SessionCreateRequest, spotCount int, spotPerRow int) (entity.Session, error) {
	var session entity.Session
	copier.Copy(&session, &sessionDTO)

	newSession, err := sessionS.sessionRepository.CreateNewSession(ctx, nil, session)
	if err != nil {
		return entity.Session{}, err
	}

	i, j := 1, 1
	rowCount := spotCount / spotPerRow

	// Create Spots according to spot_count and spot_per_row
	for i <= rowCount {
		j = 1
		for j <= spotPerRow {
			spot := entity.Spot{
				Row:       string(utils.IntToChar(i)),
				Number:    j,
				SessionID: newSession.ID,
			}

			_, err := sessionS.spotRepository.CreateNewSpot(ctx, nil, spot)
			if err != nil {
				return entity.Session{}, err
			}
			j++
		}
		i++
	}

	return newSession, nil
}

func (sessionS *sessionService) GetAllSessions(ctx context.Context) ([]entity.Session, error) {
	sessions, err := sessionS.sessionRepository.GetAllSessions(ctx, nil)
	if err != nil {
		return []entity.Session{}, err
	}
	return sessions, nil
}

func (sessionS *sessionService) DeleteSessionByID(ctx context.Context, id uint64) error {
	err := sessionS.sessionRepository.DeleteSessionByID(ctx, nil, id)
	if err != nil {
		return err
	}

	err = sessionS.spotRepository.DeleteSpotsBySessionID(ctx, nil, id)
	if err != nil {
		return err
	}

	return nil
}

func (sessionS *sessionService) GetSessionDetailByID(ctx context.Context, id uint64) (entity.Session, error) {
	session, err := sessionS.sessionRepository.GetSessionDetailByID(ctx, nil, id)
	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}
