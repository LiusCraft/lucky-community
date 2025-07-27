package model

import (
	"time"
)

// WelfareItem 福利项目模型
type WelfareItem struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement" db:"id"`
	Title         string    `json:"title" gorm:"not null;size:255" db:"title"`
	Description   string    `json:"description" gorm:"type:text" db:"description"`
	DetailContent string    `json:"detail_content" gorm:"type:text" db:"detail_content"`
	Tag           string    `json:"tag" gorm:"not null;size:50" db:"tag"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);default:0" db:"price"`
	OriginalPrice float64   `json:"original_price" gorm:"type:decimal(10,2);default:0" db:"original_price"`
	DiscountText  string    `json:"discount_text" gorm:"size:50" db:"discount_text"`
	ActionText    string    `json:"action_text" gorm:"size:50;default:立即查看" db:"action_text"`
	Status        string    `json:"status" gorm:"type:enum('active','inactive');default:active" db:"status"`
	SortOrder     int       `json:"sort_order" gorm:"default:0" db:"sort_order"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
}

// TableName 指定表名
func (WelfareItem) TableName() string {
	return "welfare_items"
}

// WelfareItemRequest 福利项目请求结构体
type WelfareItemRequest struct {
	Title         string  `json:"title" binding:"required" validate:"required,max=255"`
	Description   string  `json:"description"`
	DetailContent string  `json:"detail_content"`
	Tag           string  `json:"tag" binding:"required" validate:"required,max=50"`
	Price         float64 `json:"price" validate:"min=0"`
	OriginalPrice float64 `json:"original_price" validate:"min=0"`
	DiscountText  string  `json:"discount_text" validate:"max=50"`
	ActionText    string  `json:"action_text" validate:"max=50"`
	Status        string  `json:"status" validate:"omitempty,oneof=active inactive"`
	SortOrder     int     `json:"sort_order"`
}

// WelfareListRequest 福利列表请求结构体
type WelfareListRequest struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	Tag      string `form:"tag" json:"tag"`
	Status   string `form:"status" json:"status" validate:"omitempty,oneof=active inactive"`
}

// WelfareListResponse 福利列表响应结构体
type WelfareListResponse struct {
	Items []WelfareItem `json:"items"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}
