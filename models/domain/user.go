package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	FullName     string    `gorm:"type:varchar(100);not null"`
	IsVerified   bool      `gorm:"default:false"`
	Role         string    `gorm:"type:varchar(50);default:'user'"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
