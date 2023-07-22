package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// bases

type Base struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:uuid;not null;"`
	CreatedAt  time.Time  `gorm:"autoCreateTime;not null;default:current_timestamp;"`
	ModifiedAt time.Time  `gorm:"autoUpdateTime;not null;default:current_timestamp;"`
	DeletedAt  *time.Time `sql:"index"`
}

type BaseWithGeneratedUUID struct {
	Base
}

func (base *BaseWithGeneratedUUID) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID, _ = uuid.NewV4()
	return
}

// models

type Account struct {
	Base
	PasswordHash string `gorm:"not null"`
	Verified     bool   `gorm:"default:false;"`
}

type MFA struct {
	BaseWithGeneratedUUID
	UserID      uuid.UUID `gorm:"type:uuid;column:user;not null;unique;"`
	Type        string    `gorm:"size:4;not null;"`
	Secret      string    `gorm:"size:64;"`
	LastUsedOTP string    `gorm:"size:6;"`
}

type Session struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;not null;default:gen_random_uuid();"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;default:current_timestamp;"`
	UserID    uuid.UUID `gorm:"type: uuid;column:user;not null; onDelete: cascade; onUpdate: cascade;"`
	IP        *string   `gorm:"type:inet;" sql:"index"`
}
