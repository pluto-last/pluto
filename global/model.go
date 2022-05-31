package global

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type UUID struct {
	ID        string     `json:"id" form:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}

func (uuid *UUID) BeforeCreate(scope *gorm.Scope) error {
	var err error
	if uuid.ID == "" {
		err = scope.SetColumn("ID", RandUUID())
	}
	return err
}

func RandUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
