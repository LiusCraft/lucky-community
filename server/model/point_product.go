package model

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// 商品状态常量
const (
	ProductStatusActive   = "active"   // 启用
	ProductStatusInactive = "inactive" // 禁用
)

// PointProduct 积分商品模型
type PointProduct struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"not null;size:100"`                       // 商品名称
	Description string          `json:"description" gorm:"type:text"`                        // 商品描述
	RewardType  string          `json:"reward_type" gorm:"not null;size:20;index"`           // 奖励类型
	PricePoints int             `json:"price_points" gorm:"not null;index"`                  // 所需积分
	Stock       int             `json:"stock" gorm:"default:-1"`                             // 库存数量，-1表示无限制
	Status      string          `json:"status" gorm:"default:active;size:20;index"`          // 状态
	SortOrder   int             `json:"sort_order" gorm:"default:0;index"`                   // 排序权重
	ImageURL    string          `json:"image_url" gorm:"size:255"`                           // 商品图片URL
	CreatedAt   time.LocalTime  `json:"created_at"`
	UpdatedAt   time.LocalTime  `json:"updated_at"`
}

// TableName 指定表名
func (PointProduct) TableName() string {
	return "point_products"
}

// PointProductModel 获取积分商品模型的数据库实例
func PointProductModel() *gorm.DB {
	return mysql.GetInstance().Model(&PointProduct{})
}

// 验证方法

// IsActive 判断商品是否启用
func (pp *PointProduct) IsActive() bool {
	return pp.Status == ProductStatusActive
}

// IsInactive 判断商品是否禁用
func (pp *PointProduct) IsInactive() bool {
	return pp.Status == ProductStatusInactive
}

// HasUnlimitedStock 判断是否无限库存
func (pp *PointProduct) HasUnlimitedStock() bool {
	return pp.Stock == -1
}

// IsInStock 判断是否有库存
func (pp *PointProduct) IsInStock(quantity int) bool {
	if pp.HasUnlimitedStock() {
		return true
	}
	return pp.Stock >= quantity
}

// CanPurchase 判断是否可以购买
func (pp *PointProduct) CanPurchase(quantity int) bool {
	return pp.IsActive() && pp.IsInStock(quantity)
}

// GetRewardTypeDescription 获取奖励类型描述
func (pp *PointProduct) GetRewardTypeDescription() string {
	descriptions := map[string]string{
		RewardTypeCash:      "现金兑换",
		RewardTypeService:   "虚拟服务",
		RewardTypeProduct:   "实物商品",
		RewardTypePrivilege: "社区特权",
		RewardTypeCoupon:    "优惠券",
		RewardTypeManual:    "手动处理",
		RewardTypeOther:     "其他",
	}
	
	if desc, exists := descriptions[pp.RewardType]; exists {
		return desc
	}
	return pp.RewardType
}

// GetStatusDescription 获取状态描述
func (pp *PointProduct) GetStatusDescription() string {
	descriptions := map[string]string{
		ProductStatusActive:   "启用",
		ProductStatusInactive: "禁用",
	}
	
	if desc, exists := descriptions[pp.Status]; exists {
		return desc
	}
	return pp.Status
}

// GetStockDescription 获取库存描述
func (pp *PointProduct) GetStockDescription() string {
	if pp.HasUnlimitedStock() {
		return "无限制"
	}
	if pp.Stock <= 0 {
		return "缺货"
	}
	return fmt.Sprintf("%d件", pp.Stock)
}

// DeductStock 扣减库存
func (pp *PointProduct) DeductStock(quantity int) error {
	if !pp.HasUnlimitedStock() {
		if pp.Stock < quantity {
			return fmt.Errorf("库存不足")
		}
		pp.Stock -= quantity
	}
	return nil
}