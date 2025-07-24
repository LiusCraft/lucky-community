package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// UserPoints 用户积分账户模型
type UserPoints struct {
	ID              int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int             `json:"user_id" gorm:"not null;uniqueIndex"`
	TotalEarned     int             `json:"total_earned" gorm:"default:0"`      // 累计获得积分
	AvailablePoints int             `json:"available_points" gorm:"default:0"`  // 当前可用积分
	TotalSpent      int             `json:"total_spent" gorm:"default:0"`       // 累计消费积分
	CreatedAt       time.LocalTime  `json:"created_at"`
	UpdatedAt       time.LocalTime  `json:"updated_at"`
}

// TableName 指定表名
func (UserPoints) TableName() string {
	return "user_points"
}

// UserPoint 获取用户积分模型的数据库实例
func UserPoint() *gorm.DB {
	return mysql.GetInstance().Model(&UserPoints{})
}

// 验证方法

// HasSufficientPoints 检查用户是否有足够的积分
func (up *UserPoints) HasSufficientPoints(points int) bool {
	return up.AvailablePoints >= points
}

// GetPointsBalance 获取积分余额信息
func (up *UserPoints) GetPointsBalance() map[string]int {
	return map[string]int{
		"total_earned":     up.TotalEarned,
		"available_points": up.AvailablePoints,
		"total_spent":      up.TotalSpent,
	}
}