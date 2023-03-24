package repository

import (
	"context"
	"errors"
	"fp-rpl/entity"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

type SessionRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	GetSessionByTimeAndAreaID(ctx context.Context, tx *gorm.DB, time string, areaID uint64) (entity.Session, error)
	GetSessionByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Session, error)
	CreateNewSession(ctx context.Context, tx *gorm.DB, session entity.Session) (entity.Session, error)
	GetAllSessions(ctx context.Context, tx *gorm.DB) ([]entity.Session, error)
	DeleteSessionByID(ctx context.Context, tx *gorm.DB, id uint64) error
}

func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db: db}
}

func (sessionR *sessionRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := sessionR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (sessionR *sessionRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (sessionR *sessionRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (sessionR *sessionRepository) GetSessionByTimeAndAreaID(ctx context.Context, tx *gorm.DB, time string, areaID uint64) (entity.Session, error) {
	var err error
	var session entity.Session
	if tx == nil {
		tx = sessionR.db.WithContext(ctx).Debug().Where("time = $1 AND area_id = $2", time, areaID).Take(&session)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("time = $1 AND area_id = $2", time, areaID).Take(&session).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return session, err
	}
	return session, nil
}

func (sessionR *sessionRepository) GetSessionByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Session, error) {
	var err error
	var session entity.Session
	if tx == nil {
		tx = sessionR.db.WithContext(ctx).Debug().Where("id = $1", id).Take(&session)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&session).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return session, err
	}
	return session, nil
}

func (sessionR *sessionRepository) CreateNewSession(ctx context.Context, tx *gorm.DB, session entity.Session) (entity.Session, error) {
	var err error
	if tx == nil {
		tx = sessionR.db.WithContext(ctx).Debug().Create(&session)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&session).Error
	}

	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (sessionR *sessionRepository) GetAllSessions(ctx context.Context, tx *gorm.DB) ([]entity.Session, error) {
	var err error
	var sessions []entity.Session

	if tx == nil {
		tx = sessionR.db.WithContext(ctx).Debug().Preload("Spots").Find(&sessions)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Preload("Spots").Find(&sessions).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return sessions, err
	}
	return sessions, nil
}

func (sessionR *sessionRepository) DeleteSessionByID(ctx context.Context, tx *gorm.DB, id uint64) error {
	var err error
	if tx == nil {
		tx = sessionR.db.WithContext(ctx).Debug().Delete(&entity.Session{}, id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Session{}, id).Error
	}

	if err != nil {
		return err
	}
	return nil
}