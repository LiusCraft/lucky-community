package services

import (
	"fmt"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type PointConfigService struct{}

// getPointConfigDao 获取积分配置DAO实例
func (s *PointConfigService) getPointConfigDao() *dao.PointConfigDao {
	return &dao.PointConfigDao{}
}

// GetPointConfig 获取积分配置
func (s *PointConfigService) GetPointConfig() (*model.PointConfig, error) {
	pointConfigDao := s.getPointConfigDao()
	return pointConfigDao.GetPointConfig()
}

// UpdatePointConfig 更新积分配置
func (s *PointConfigService) UpdatePointConfig(rulesDescription string, inviteRewardPoints int) error {
	if rulesDescription == "" {
		return fmt.Errorf("积分规则说明不能为空")
	}
	
	if inviteRewardPoints < 0 {
		return fmt.Errorf("邀请奖励积分不能为负数")
	}
	
	pointConfigDao := s.getPointConfigDao()
	
	// 先获取当前配置
	currentConfig, err := pointConfigDao.GetPointConfig()
	if err != nil {
		return fmt.Errorf("获取当前配置失败: %v", err)
	}
	
	// 更新配置
	currentConfig.RulesDescription = rulesDescription
	currentConfig.InviteRewardPoints = inviteRewardPoints
	
	return pointConfigDao.UpdatePointConfig(currentConfig)
}

// GetInviteRewardPoints 获取邀请奖励积分
func (s *PointConfigService) GetInviteRewardPoints() (int, error) {
	pointConfigDao := s.getPointConfigDao()
	return pointConfigDao.GetInviteRewardPoints()
}