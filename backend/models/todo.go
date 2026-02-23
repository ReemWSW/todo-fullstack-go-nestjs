package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo represents a todo item
type Todo struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
}

// BeforeCreate will set a UUID for the Todo ID
func (todo *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	todo.ID = uuid.New()
	return
}
