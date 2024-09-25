package repository

import (
	"context"
	"log/slog"
	"online-song-library/internal/model"
	"online-song-library/pkg/storage/postgresql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(conn *gorm.DB) *SongRepository {
	return &SongRepository{
		db: conn,
	}
}

// http запросы на обогащение будут проводиться в service
func (r *SongRepository) Create(ctx context.Context, log *slog.Logger, song model.Song) (uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	default:
	}

	if err := postgresql.TxSaveExecutor(r.db, func(d *gorm.DB) error {
		log.Debug("Create sql query:", 
			slog.String("id", song.Id.String()), 
			slog.String("gruop", song.Group),
			slog.String("title", song.Title),
			slog.Time("release_date", song.ReleaseDate),
			slog.String("link", song.Link))

		if result := d.Create(song); result.Error != nil {
			return result.Error
		}
		return nil
	}); err != nil {
		return uuid.Nil, err
	}
	return song.Id, nil
}

func (r *SongRepository) Update(ctx context.Context, log *slog.Logger, song model.Song) (model.Song, error) {
	select {
	case <-ctx.Done():
		return model.Song{}, ctx.Err()

	default:
	}

	var oldModel model.Song
	if err := postgresql.TxSaveExecutor(r.db, func(d *gorm.DB) error {
		log.Debug("Update sql query:", 
			slog.String("id", song.Id.String()), 
			slog.String("gruop", song.Group),
			slog.String("title", song.Title),
			slog.Time("release_date", song.ReleaseDate),
			slog.String("link", song.Link))

		if result := d.First(&oldModel, "id = ?", song.Id); result.Error != nil {
			return result.Error
		}

		if result := d.Model(&oldModel).Updates(song); result.Error != nil {
			return result.Error
		}
		return nil
	}); err != nil {
		return model.Song{}, err
	}
	return oldModel, nil
}

func (r *SongRepository) Delete(ctx context.Context, log *slog.Logger, songUUID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := postgresql.TxSaveExecutor(r.db, func(d *gorm.DB) error {
		var model model.Song

		log.Debug("Delete sql query:", 
			slog.String("id", songUUID.String()))

		if result := d.First(&model, "id = ?", songUUID); result.Error != nil {
			return result.Error
		}
		if err := d.Delete(&model).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *SongRepository) GetAll(ctx context.Context, log *slog.Logger, limit int, offset int, filter model.SongFilter) ([]model.Song, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var models []model.Song
	if err := postgresql.TxSaveExecutor(r.db, func(d *gorm.DB) error {
		query := d.Limit(limit).Offset(offset)

		log.Debug("GetAll sql query:", slog.Int("limit", limit), slog.Int("offset", offset))

		if filter.Id != nil {
			query = query.Where("id = ?", *filter.Id)
			log.Debug("filter detected", slog.String("filter_id", (*filter.Id).String()))
		}
		if filter.Group != nil {
			query = query.Where("group = ?", *filter.Group)
			log.Debug("filter detected", slog.String("filter_group", (*filter.Group)))
		}
		if filter.Title != nil {
			query = query.Where("title = ?", *filter.Title)
			log.Debug("filter detected", slog.String("filter_title", (*filter.Group)))
		}
		if filter.ReleaseDate != nil {
			query = query.Where("release_date = ?", *filter.ReleaseDate)
			log.Debug("filter detected", slog.String("filter_release_date", (*filter.ReleaseDate).String()))
		}
		if filter.Text != nil {
			query = query.Where("text = ?", *filter.Text)
			log.Debug("filter detected", slog.String("filter_text", (*filter.Text)))

		}
		if filter.Link != nil {
			query = query.Where("link = ?", *filter.Link)
			log.Debug("filter detected", slog.String("filter_link", (*filter.Link)))
		}

		res := query.Find(&models)
		if res.Error != nil {
			return res.Error
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return models, nil
}


func (r *SongRepository) GetVerses(ctx context.Context, log *slog.Logger, songUUID uuid.UUID) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	var verses string
	if err := postgresql.TxSaveExecutor(r.db, func(d *gorm.DB) error {
		log.Debug("GetVerses sql query:", 
			slog.String("id", songUUID.String()))

		res := d.Select("text").Where("id = ?", songUUID).First(&verses)
		if res.Error != nil {
			return res.Error
		}
		return nil
	}); err != nil {
		return "", err
	}
	return verses, nil
}
