package repository

import (
	"context"
	"errors"
	"fp-rpl/dto"
	"fp-rpl/entity"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type areaRepository struct {
	db *gorm.DB
}

type AreaRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	GetAreaByName(ctx context.Context, tx *gorm.DB, name string) (entity.Area, error)
	CreateNewArea(ctx context.Context, tx *gorm.DB, area entity.Area) (entity.Area, error)
	GetAllAreas(ctx context.Context, tx *gorm.DB) ([]entity.Area, error)
	GetAreaByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Area, error)
	UpdateArea(ctx context.Context, tx *gorm.DB, areaDTO dto.AreaCreateRequest, area entity.Area) (entity.Area, error)
	DeleteAreaByID(ctx context.Context, tx *gorm.DB, id uint64) error
}

func NewAreaRepository(db *gorm.DB) *areaRepository {
	return &areaRepository{db: db}
}

func (areaR *areaRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := areaR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (areaR *areaRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (areaR *areaRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (areaR *areaRepository) GetAreaByName(ctx context.Context, tx *gorm.DB, name string) (entity.Area, error) {
	var err error
	var area entity.Area
	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Where("name = $1", name).Take(&area)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("name = $1", name).Take(&area).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return area, err
	}
	return area, nil
}

func (areaR *areaRepository) CreateNewArea(ctx context.Context, tx *gorm.DB, area entity.Area) (entity.Area, error) {
	var err error
	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Create(&area)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&area).Error
	}

	if err != nil {
		return entity.Area{}, err
	}
	return area, nil
}

func (areaR *areaRepository) GetAllAreas(ctx context.Context, tx *gorm.DB) ([]entity.Area, error) {
	var err error
	var areas []entity.Area

	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Find(&areas)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Find(&areas).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return areas, err
	}
	return areas, nil
}

func (areaR *areaRepository) GetAreaByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Area, error) {
	var err error
	var area entity.Area
	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Where("id = $1", id).Preload("Sessions").Take(&area)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Preload("Sessions").Take(&area).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return area, err
	}
	return area, nil
}

func (areaR *areaRepository) UpdateArea(ctx context.Context, tx *gorm.DB, areaDTO dto.AreaCreateRequest, area entity.Area) (entity.Area, error) {
	var err error
	areaUpdate := area
	copier.Copy(&areaUpdate, &areaDTO)
	
	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Save(&areaUpdate)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Save(&areaUpdate).Error
	}

	if err != nil {
		return areaUpdate, err
	}
	return areaUpdate, nil
}

func (areaR *areaRepository) DeleteAreaByID(ctx context.Context, tx *gorm.DB, id uint64) error {
	var err error
	if tx == nil {
		tx = areaR.db.WithContext(ctx).Debug().Delete(&entity.Area{}, id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Area{}, id).Error
	}

	if err != nil {
		return err
	}
	return nil
}
