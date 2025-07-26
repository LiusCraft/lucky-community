package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// PointConfig 积分配置模型
type PointConfig struct {
	ID                  int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	RulesDescription    string          `json:"rules_description" gorm:"type:text;not null"`         // 积分规则说明
	InviteRewardPoints  int             `json:"invite_reward_points" gorm:"not null;default:10"`     // 邀请用户注册积分奖励
	CreatedAt          time.LocalTime  `json:"created_at"`
	UpdatedAt          time.LocalTime  `json:"updated_at"`
}

// TableName 指定表名
func (PointConfig) TableName() string {
	return "point_configs"
}

// PointConfigModel 获取积分配置模型的数据库实例
func PointConfigModel() *gorm.DB {
	return mysql.GetInstance().Model(&PointConfig{})
}