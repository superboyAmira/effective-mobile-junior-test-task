package model

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Group       string    `gorm:"type:varchar(1000);not null" json:"group"`
	Title       string    `gorm:"type:varchar(1000);not null" json:"song"`
	ReleaseDate time.Time `gorm:"type:timestamp;default:current_timestamp" json:"release_date"`
	Text        string    `gorm:"type:text" json:"text"`
	Link        string    `gorm:"type:varchar(500);unique;not null" json:"link"`
}

type SongDTO struct {
	Group string `json:"group"`
	Title string `json:"song"`
}

type SongFilter struct {
	Id          *uuid.UUID  `json:"id,omitempty"`
	Group       *string    `json:"group,omitempty"`
	Title       *string    `json:"song,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Text        *string    `json:"text,omitempty"`
	Link        *string    `json:"link,omitempty"`
}
