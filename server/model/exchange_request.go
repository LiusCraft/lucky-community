package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// 兑换申请状态常量
const (
	ExchangeStatusPending   = "pending"   // 待处理
	ExchangeStatusProcessed = "processed" // 已处理
)

// ExchangeRequest 兑换申请模型
type ExchangeRequest struct {
	ID          int64              `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int                `json:"user_id" gorm:"not null;index"`                       // 用户ID
	ProductID   int64              `json:"product_id" gorm:"not null;index"`                    // 商品ID
	PointsCost  int                `json:"points_cost" gorm:"not null"`                         // 消费积分数量
	Quantity    int                `json:"quantity" gorm:"default:1"`                           // 兑换数量
	Status      string             `json:"status" gorm:"default:pending;size:20;index"`         // 状态
	UserInfo    string             `json:"user_info" gorm:"type:text"`                          // 用户填写的信息
	ProcessedAt *time.LocalTime    `json:"processed_at"`                                        // 处理完成时间
	CreatedAt   time.LocalTime     `json:"created_at" gorm:"index"`
	UpdatedAt   time.LocalTime     `json:"updated_at"`
	
	// 关联数据
	Product *PointProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`     // 关联商品
	User    *Users        `json:"user,omitempty" gorm:"foreignKey:UserID"`           // 关联用户
}

// TableName 指定表名
func (ExchangeRequest) TableName() string {
	return "exchange_requests"
}

// ExchangeRequestModel 获取兑换申请模型的数据库实例
func ExchangeRequestModel() *gorm.DB {
	return mysql.GetInstance().Model(&ExchangeRequest{})
}

// 验证方法

// IsPending 判断是否为待处理状态
func (er *ExchangeRequest) IsPending() bool {
	return er.Status == ExchangeStatusPending
}

// IsProcessed 判断是否为已处理状态
func (er *ExchangeRequest) IsProcessed() bool {
	return er.Status == ExchangeStatusProcessed
}

// CanProcess 判断是否可以处理
func (er *ExchangeRequest) CanProcess() bool {
	return er.IsPending()
}

// GetStatusDescription 获取状态描述
func (er *ExchangeRequest) GetStatusDescription() string {
	descriptions := map[string]string{
		ExchangeStatusPending:   "待处理",
		ExchangeStatusProcessed: "已处理",
	}
	
	if desc, exists := descriptions[er.Status]; exists {
		return desc
	}
	return er.Status
}

// GetTotalPointsCost 获取总积分消费
func (er *ExchangeRequest) GetTotalPointsCost() int {
	return er.PointsCost * er.Quantity
}