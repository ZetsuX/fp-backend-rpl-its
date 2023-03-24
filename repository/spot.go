package repository

import (
	"context"
	"fp-rpl/entity"

	"gorm.io/gorm"
)

type spotRepository struct {
	db *gorm.DB
}

type SpotRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	CreateNewSpot(ctx context.Context, tx *gorm.DB, spot entity.Spot) (entity.Spot, error)
	DeleteSpotsBySessionID(ctx context.Context, tx *gorm.DB, spotID uint64) error
}

func NewSpotRepository(db *gorm.DB) *spotRepository {
	return &spotRepository{db: db}
}

func (spotR *spotRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := spotR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (spotR *spotRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (spotR *spotRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (spotR *spotRepository) CreateNewSpot(ctx context.Context, tx *gorm.DB, spot entity.Spot) (entity.Spot, error) {
	var err error
	if tx == nil {
		tx = spotR.db.WithContext(ctx).Debug().Create(&spot)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&spot).Error
	}

	if err != nil {
		return entity.Spot{}, err
	}
	return spot, nil
}

func (spotR *spotRepository) DeleteSpotsBySessionID(ctx context.Context, tx *gorm.DB, spotID uint64) error {
	var err error
	if tx == nil {
		tx = spotR.db.WithContext(ctx).Debug().Where("session_id = $1", spotID).Unscoped().Delete(&entity.Spot{})
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("session_id = $1", spotID).Unscoped().Delete(&entity.Spot{}).Error
	}

	if err != nil {
		return err
	}
	return nil
}