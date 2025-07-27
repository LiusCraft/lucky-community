package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// InviteRelation 邀请关系模型
type InviteRelation struct {
	ID         int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	InviterID  int             `json:"inviter_id" gorm:"not null;index"`       // 邀请人用户ID
	InviteeID  int             `json:"invitee_id" gorm:"not null;index"`       // 被邀请人用户ID
	InviteCode string          `json:"invite_code" gorm:"not null;size:20;index"` // 使用的邀请码
	CreatedAt  time.LocalTime  `json:"created_at" gorm:"index"`
}

// TableName 指定表名
func (InviteRelation) TableName() string {
	return "invite_relations"
}

// InviteRelationModel 获取邀请关系模型的数据库实例
func InviteRelationModel() *gorm.DB {
	return mysql.GetInstance().Model(&InviteRelation{})
}