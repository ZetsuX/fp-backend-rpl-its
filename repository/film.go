package repository

import (
	"context"
	"errors"
	"fp-rpl/dto"
	"fp-rpl/entity"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type filmRepository struct {
	db *gorm.DB
}

type FilmRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	CreateNewFilm(ctx context.Context, tx *gorm.DB, user entity.Film) (entity.Film, error)
	GetAllFilms(ctx context.Context, tx *gorm.DB) ([]entity.Film, error)
	GetFilmBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Film, error)
	GetFilmByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Film, error)
	GetFilmDetailBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Film, error)
	UpdateFilmBySlug(ctx context.Context, tx *gorm.DB, filmDTO dto.FilmRegisterRequest, film entity.Film) (entity.Film, error)
	DeleteFilm(ctx context.Context, tx *gorm.DB, slug string) error
	GetAllFilmsByStatus(ctx context.Context, tx *gorm.DB,status string) ([]entity.Film, error)
}

func NewFilmRepository(db *gorm.DB) FilmRepository {
	return &filmRepository{db: db}
}
func (filmR *filmRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := filmR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (filmR *filmRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (filmR *filmRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (filmR *filmRepository) CreateNewFilm(ctx context.Context, tx *gorm.DB, film entity.Film) (entity.Film, error) {
	var err error
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Create(&film)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&film).Error
	}

	if err != nil {
		return entity.Film{}, err
	}
	return film, nil
}

func (filmR *filmRepository) GetAllFilms(ctx context.Context, tx *gorm.DB) ([]entity.Film, error) {
	var err error
	var films []entity.Film

	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Find(&films)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Find(&films).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return films, err
	}
	return films, nil
}
func (filmR *filmRepository) GetAllFilmsByStatus(ctx context.Context, tx *gorm.DB,status string) ([]entity.Film, error) {
	var err error
	var films []entity.Film

	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Where("status = ?",status).Find(&films)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("status = ?",status).Find(&films).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return films, err
	}
	return films, nil
}
func (filmR *filmRepository) GetFilmBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Film, error) {
	var err error
	var film entity.Film
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Where("slug = $1", slug).Take(&film)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("slug = $1", slug).Take(&film).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return film, err
	}
	return film, nil
}

func (filmR *filmRepository) GetFilmByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Film, error) {
	var err error
	var film entity.Film
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Where("id = $1", id).Take(&film)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&film).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return film, err
	}
	return film, nil
}

func (filmR *filmRepository) GetFilmDetailBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Film, error) {
	var err error
	var film entity.Film
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Where("slug = $1", slug).Preload("Sessions").Take(&film)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("slug = $1", slug).Preload("Sessions").Take(&film).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return film, err
	}
	return film, nil
}

func (filmR *filmRepository) UpdateFilmBySlug(ctx context.Context, tx *gorm.DB, filmDTO dto.FilmRegisterRequest, film entity.Film) (entity.Film, error) {
	var err error
	filmUpdate := film
	copier.Copy(&filmUpdate,&filmDTO)

	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug()
	}

	tx = tx.Save(&filmUpdate)
	err = tx.Error

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return filmUpdate, err
	}

	return filmUpdate, nil
}

func (filmR *filmRepository) DeleteFilm(ctx context.Context, tx *gorm.DB, slug string) error {
	var err error
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Where("slug = ?", slug).Delete(&entity.Film{})
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("slug = ?", slug).Delete(&entity.Film{}).Error
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return err
	}
	return nil
}
