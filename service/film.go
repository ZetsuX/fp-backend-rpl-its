package service

import (
	"context"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/repository"

	"github.com/jinzhu/copier"
)

type filmService struct {
	filmRepository repository.FilmRepository
}

type FilmService interface {
	CreateNewFilm(ctx context.Context, filmDTO dto.FilmRegisterRequest) (entity.Film, error)
	GetFilmBySlug(ctx context.Context, slug string) (entity.Film, error)
	GetFilmByID(ctx context.Context, id uint64) (entity.Film, error)
	GetFilmDetailBySlug(ctx context.Context, slug string) (entity.Film, error)
	GetAllFilm(ctx context.Context) ([]entity.Film, error)
	UpdateFilm(ctx context.Context, filmDTO dto.FilmUpdateRequest, slug string) (dto.FilmUpdateRequest, error)
	DeleteFilm(ctx context.Context, slug string) error
	GetAllFilmAvailable(ctx context.Context) ([]entity.Film, error)
}

func NewFilmService(filmR repository.FilmRepository) FilmService {
	return &filmService{filmRepository: filmR}
}
func (fs *filmService) CreateNewFilm(ctx context.Context, filmDTO dto.FilmRegisterRequest) (entity.Film, error) {
	filmDTO.Status = "playing"
	var film entity.Film
	copier.Copy(&film, &filmDTO)

	NewFilm, err := fs.filmRepository.CreateNewFilm(ctx, nil, film)
	if err != nil {
		return entity.Film{}, err
	}
	return NewFilm, nil
}
func (fs *filmService) GetFilmBySlug(ctx context.Context, slug string) (entity.Film, error) {
	film, err := fs.filmRepository.GetFilmBySlug(ctx, nil, slug)
	if err != nil {
		return entity.Film{}, err
	}
	return film, nil
}

func (fs *filmService) GetFilmByID(ctx context.Context, id uint64) (entity.Film, error) {
	film, err := fs.filmRepository.GetFilmByID(ctx, nil, id)
	if err != nil {
		return entity.Film{}, err
	}
	return film, nil
}

func (fs *filmService) GetFilmDetailBySlug(ctx context.Context, slug string) (entity.Film, error) {
	film, err := fs.filmRepository.GetFilmDetailBySlug(ctx, nil, slug)
	if err != nil {
		return entity.Film{}, err
	}
	return film, nil
}

func (fs *filmService) GetAllFilm(ctx context.Context) ([]entity.Film, error) {
	films, err := fs.filmRepository.GetAllFilms(ctx, nil)
	if err != nil {
		return []entity.Film{}, err
	}
	return films, nil
}
func (fs *filmService) GetAllFilmAvailable(ctx context.Context) ([]entity.Film, error) {
	films, err := fs.filmRepository.GetAllFilmsByStatus(ctx,nil,"playing")
	if err != nil {
		return []entity.Film{}, err
	}
	return films, nil
}

func (fs *filmService) UpdateFilm(ctx context.Context, filmDTO dto.FilmUpdateRequest, slug string) (dto.FilmUpdateRequest, error) {
	film, err := fs.filmRepository.UpdateFilmBySlug(ctx, nil, slug, filmDTO)
	if err != nil {
		return dto.FilmUpdateRequest{}, err
	}
	return film, nil
}

func (fs *filmService) DeleteFilm(ctx context.Context, slug string) error {
	err := fs.filmRepository.DeleteFilm(ctx, nil, slug)
	if err != nil {
		return err
	}
	return nil
}
