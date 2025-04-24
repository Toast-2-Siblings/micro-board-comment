package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
    ID uint `gorm:"primaryKey"`
    BoardID uint `gorm:"index;not null"`  // 게시글 ID (Board 서비스의 Board ID)
    UserID uint `gorm:"not null"` // 작성자 ID
    Content string `gorm:"type:text;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
