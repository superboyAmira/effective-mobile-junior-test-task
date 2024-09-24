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
