package model

import (
	"crypto/rand"
	"fmt"
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
	localTime "xhyovo.cn/community/pkg/time"
)

// UserInviteCode 用户邀请码模型
type UserInviteCode struct {
	ID         int64               `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int                 `json:"user_id" gorm:"not null;uniqueIndex"`             // 用户ID
	InviteCode string              `json:"invite_code" gorm:"not null;uniqueIndex;size:20"` // 邀请码
	CreatedAt  localTime.LocalTime `json:"created_at"`
	UpdatedAt  localTime.LocalTime `json:"updated_at"`

	// 关联数据
	User *Users `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (UserInviteCode) TableName() string {
	return "user_invite_codes"
}

// UserInviteCodeModel 获取用户邀请码模型的数据库实例
func UserInviteCodeModel() *gorm.DB {
	return mysql.GetInstance().Model(&UserInviteCode{})
}

// 验证方法

// 静态方法

// GenerateInviteCode 生成邀请码
func GenerateInviteCode(userID int) string {
	// 生成格式：USR{userID}{随机字符串}
	randomStr := generateRandomString(6)
	return fmt.Sprintf("USR%d%s", userID, randomStr)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("%d", time.Now().Unix())[:length]
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}
