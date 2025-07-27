package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// 积分操作类型常量
const (
	PointTypeEarn  = "earn"  // 积分获得
	PointTypeSpend = "spend" // 积分消费
)

// 积分来源类型常量
const (
	SourceTypeInvite   = "INVITE"   // 邀请注册
	SourceTypeCourse   = "COURSE"   // 课程相关
	SourceTypeContent  = "CONTENT"  // 内容贡献
	SourceTypeDaily    = "DAILY"    // 日常活跃
	SourceTypeActivity = "ACTIVITY" // 运营活动
	SourceTypeManual   = "MANUAL"   // 手动发放
	SourceTypeOther    = "OTHER"    // 其他类型
)

// 兑换奖励类型常量
const (
	RewardTypeCash      = "CASH"      // 现金兑换
	RewardTypeService   = "SERVICE"   // 虚拟服务
	RewardTypeProduct   = "PRODUCT"   // 实物商品
	RewardTypePrivilege = "PRIVILEGE" // 社区特权
	RewardTypeCoupon    = "COUPON"    // 优惠券类
	RewardTypeManual    = "MANUAL"    // 手动处理
	RewardTypeOther     = "OTHER"     // 其他类型
)

// PointRecord 积分记录模型
type PointRecord struct {
	ID                int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            int             `json:"user_id" gorm:"not null;index"`
	Type              string          `json:"type" gorm:"not null;size:20"`                    // 积分操作类型：earn获得,spend消费
	Points            int             `json:"points" gorm:"not null"`                          // 积分数量
	SourceType        *string         `json:"source_type" gorm:"size:20;index"`                // 积分来源类型
	RewardType        *string         `json:"reward_type" gorm:"size:20;index"`                // 兑换奖励类型
	Description       string          `json:"description" gorm:"type:text"`                    // 操作描述
	ExchangeRequestID *int64          `json:"exchange_request_id" gorm:"index"`                // 关联的兑换申请ID
	CreatedAt         time.LocalTime  `json:"created_at" gorm:"index"`
}

// TableName 指定表名
func (PointRecord) TableName() string {
	return "point_records"
}

// PointRecordModel 获取积分记录模型的数据库实例
func PointRecordModel() *gorm.DB {
	return mysql.GetInstance().Model(&PointRecord{})
}

// 验证方法

// IsEarnType 判断是否为积分获得类型
func (pr *PointRecord) IsEarnType() bool {
	return pr.Type == PointTypeEarn
}

// IsSpendType 判断是否为积分消费类型
func (pr *PointRecord) IsSpendType() bool {
	return pr.Type == PointTypeSpend
}

// GetSourceTypeDescription 获取来源类型描述
func (pr *PointRecord) GetSourceTypeDescription() string {
	if pr.SourceType == nil {
		return ""
	}
	
	descriptions := map[string]string{
		SourceTypeInvite:   "邀请注册",
		SourceTypeCourse:   "课程相关",
		SourceTypeContent:  "内容贡献",
		SourceTypeDaily:    "日常活跃",
		SourceTypeActivity: "运营活动",
		SourceTypeManual:   "手动发放",
		SourceTypeOther:    "其他",
	}
	
	if desc, exists := descriptions[*pr.SourceType]; exists {
		return desc
	}
	return *pr.SourceType
}

// GetRewardTypeDescription 获取奖励类型描述
func (pr *PointRecord) GetRewardTypeDescription() string {
	if pr.RewardType == nil {
		return ""
	}
	
	descriptions := map[string]string{
		RewardTypeCash:      "现金兑换",
		RewardTypeService:   "虚拟服务",
		RewardTypeProduct:   "实物商品",
		RewardTypePrivilege: "社区特权",
		RewardTypeCoupon:    "优惠券",
		RewardTypeManual:    "手动处理",
		RewardTypeOther:     "其他",
	}
	
	if desc, exists := descriptions[*pr.RewardType]; exists {
		return desc
	}
	return *pr.RewardType
}