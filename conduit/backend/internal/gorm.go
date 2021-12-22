package internal

import (
	"gorm.io/gorm"
	"time"
)

// CustomGromModel is simple struct to house the Creation/Update/Deleted information
type CustomGromModel struct {
	//ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
