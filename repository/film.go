package repository

import (
	"context"
	"errors"
	"fp-rpl/dto"
	"fp-rpl/entity"

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
	UpdateFilmBySlug(ctx context.Context, tx *gorm.DB, slug string, filmDTO dto.FilmUpdateRequest) (dto.FilmUpdateRequest, error)
}
func NewFilmRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
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
func (filmR *filmRepository) UpdateFilmBySlug(ctx context.Context, tx *gorm.DB, slug string, filmDTO dto.FilmUpdateRequest) (dto.FilmUpdateRequest, error) {
	var err error
	if tx == nil {
		tx = filmR.db.WithContext(ctx).Debug().Model(entity.Film{}).Where("slug = $1", slug).Updates(map[string]interface{}{
			"Title": filmDTO.Title,
			"Synopsis": filmDTO.Synopsis,
			"Duration": filmDTO.Duration,
			"Genre": filmDTO.Genre,
			"Producer": filmDTO.Producer,
			"Director": filmDTO.Director,
			"Writer": filmDTO.Writer,
			"Production": filmDTO.Production,
			"Cast": filmDTO.Cast,
			"Trailer": filmDTO.Trailer,
			"Image": filmDTO.Image,
			"Status": filmDTO.Status,
		})
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Model(entity.Film{}).Where("slug = $1", slug).Updates(map[string]interface{}{
			"Title": filmDTO.Title,
			"Synopsis": filmDTO.Synopsis,
			"Duration": filmDTO.Duration,
			"Genre": filmDTO.Genre,
			"Producer": filmDTO.Producer,
			"Director": filmDTO.Director,
			"Writer": filmDTO.Writer,
			"Production": filmDTO.Production,
			"Cast": filmDTO.Cast,
			"Trailer": filmDTO.Trailer,
			"Image": filmDTO.Image,
			"Status": filmDTO.Status,
		}).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return filmDTO, err
	}
	return filmDTO, nil
}
