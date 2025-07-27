package model

import (
	"time"
)

// UpdateLog 更新日志模型
type UpdateLog struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement" db:"id"`
	Title       string    `json:"title" gorm:"not null;size:255" db:"title"`
	Description string    `json:"description" gorm:"type:text" db:"description"`
	Content     string    `json:"content" gorm:"type:longtext" db:"content"`
	Version     string    `json:"version" gorm:"size:50" db:"version"`
	Type        string    `json:"type" gorm:"size:50;default:feature" db:"type"`
	Status      string    `json:"status" gorm:"size:20;default:active" db:"status"`
	PublishDate time.Time `json:"publish_date" gorm:"default:CURRENT_TIMESTAMP" db:"publish_date"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
}

// TableName 指定表名
func (UpdateLog) TableName() string {
	return "update_logs"
}

// UpdateLogRequest 更新日志请求结构体
type UpdateLogRequest struct {
	Title       string `json:"title" binding:"required" validate:"required,max=255"`
	Description string `json:"description" binding:"required"`
	Content     string `json:"content"`
	Version     string `json:"version" validate:"max=50"`
	Type        string `json:"type" validate:"max=50"`
	Status      string `json:"status" validate:"max=20"`
	PublishDate string `json:"publish_date"` // 前端传递字符串格式的日期
}

// UpdateLogListRequest 更新日志列表请求结构体
type UpdateLogListRequest struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	Type     string `form:"type" json:"type"`
	Status   string `form:"status" json:"status"`
}

// UpdateLogListResponse 更新日志列表响应结构体
type UpdateLogListResponse struct {
	Items []UpdateLog `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// UpdateLogSimple 更新日志简化结构体（用于前端展示）
type UpdateLogSimple struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	Date        string `json:"date"` // 格式化后的日期字符串
}