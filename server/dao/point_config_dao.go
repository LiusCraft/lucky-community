package dao

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type PointConfigDao struct{}

// GetPointConfig 获取积分配置（获取第一条记录，如果不存在则创建默认配置）
func (d *PointConfigDao) GetPointConfig() (*model.PointConfig, error) {
	var config model.PointConfig
	
	err := model.PointConfigModel().First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有配置记录，创建默认配置
			defaultConfig := &model.PointConfig{
				RulesDescription: `# 积分规则说明

## 如何获得积分

- **邀请好友注册**：获得 **10** 积分
- **参与社区活动**：根据活动规则获得积分

## 积分使用说明

- 积分可用于兑换商城内的各种商品
- 兑换后积分立即扣除，不可退还
- 积分不可转让给其他用户

## 其他注意事项

- 积分有效期为获得后**1年**
- 商品数量有限，**先兑先得**
- 如有疑问请联系客服

---

> 💡 **提示**：积分是社区的通用货币，合理使用可以获得更多优质服务！`,
				InviteRewardPoints: 10,
			}
			
			if createErr := model.PointConfigModel().Create(defaultConfig).Error; createErr != nil {
				return nil, createErr
			}
			return defaultConfig, nil
		}
		return nil, err
	}
	
	return &config, nil
}

// UpdatePointConfig 更新积分配置
func (d *PointConfigDao) UpdatePointConfig(config *model.PointConfig) error {
	if config.ID <= 0 {
		return gorm.ErrInvalidField
	}
	
	// 只更新指定字段
	return model.PointConfigModel().
		Where("id = ?", config.ID).
		Select("rules_description", "invite_reward_points").
		Updates(config).Error
}

// GetInviteRewardPoints 获取邀请奖励积分
func (d *PointConfigDao) GetInviteRewardPoints() (int, error) {
	config, err := d.GetPointConfig()
	if err != nil {
		return 0, err
	}
	return config.InviteRewardPoints, nil
}