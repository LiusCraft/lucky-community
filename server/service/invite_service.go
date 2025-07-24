package services

import (
	"fmt"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type InviteService struct{}

// getUserInviteCodeDao 获取用户邀请码DAO实例
func (s *InviteService) getUserInviteCodeDao() *dao.UserInviteCodeDao {
	return &dao.UserInviteCodeDao{}
}

// getInviteRelationDao 获取邀请关系DAO实例
func (s *InviteService) getInviteRelationDao() *dao.InviteRelationDao {
	return &dao.InviteRelationDao{}
}

// getPointsService 获取积分服务实例
func (s *InviteService) getPointsService() *PointsService {
	return &PointsService{}
}

// GetUserInviteCode 获取用户的邀请码
func (s *InviteService) GetUserInviteCode(userID int) (*model.UserInviteCode, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}
	
	userInviteCodeDao := s.getUserInviteCodeDao()
	return userInviteCodeDao.GetUserInviteCode(userID)
}

// ValidateInviteCode 验证邀请码是否有效
func (s *InviteService) ValidateInviteCode(inviteCode string) (*model.UserInviteCode, error) {
	if inviteCode == "" {
		return nil, fmt.Errorf("邀请码不能为空")
	}
	
	// 验证邀请码格式
	if !model.ValidateInviteCode(inviteCode) {
		return nil, fmt.Errorf("邀请码格式不正确")
	}
	
	userInviteCodeDao := s.getUserInviteCodeDao()
	return userInviteCodeDao.GetInviteCodeByCode(inviteCode)
}

// CreateInviteRelation 创建邀请关系
func (s *InviteService) CreateInviteRelation(inviterID, inviteeID int, inviteCode string) error {
	if inviterID <= 0 || inviteeID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if inviterID == inviteeID {
		return fmt.Errorf("不能邀请自己")
	}
	
	// 验证邀请码是否属于邀请人
	inviteCodeInfo, err := s.ValidateInviteCode(inviteCode)
	if err != nil {
		return fmt.Errorf("邀请码验证失败: %v", err)
	}
	
	if inviteCodeInfo.UserID != inviterID {
		return fmt.Errorf("邀请码不属于该用户")
	}
	
	// 检查被邀请人是否已有邀请关系
	inviteRelationDao := s.getInviteRelationDao()
	existingRelation, err := inviteRelationDao.GetInviteRelationByInvitee(inviteeID)
	if err != nil {
		return fmt.Errorf("检查邀请关系失败: %v", err)
	}
	
	if existingRelation != nil {
		return fmt.Errorf("该用户已被他人邀请")
	}
	
	// 创建邀请关系
	relation := &model.InviteRelation{
		InviterID:  inviterID,
		InviteeID:  inviteeID,
		InviteCode: inviteCode,
	}
	
	if err := inviteRelationDao.CreateInviteRelation(relation); err != nil {
		return fmt.Errorf("创建邀请关系失败: %v", err)
	}
	
	
	return nil
}

// ProcessInviteReward 处理邀请奖励（当被邀请用户付费成功时调用）
func (s *InviteService) ProcessInviteReward(inviteeID int, rewardPoints int) error {
	if inviteeID <= 0 {
		return fmt.Errorf("无效的被邀请用户ID")
	}
	
	if rewardPoints <= 0 {
		return fmt.Errorf("奖励积分必须大于0")
	}
	
	// 查找邀请关系
	inviteRelationDao := s.getInviteRelationDao()
	relation, err := inviteRelationDao.GetInviteRelationByInvitee(inviteeID)
	if err != nil {
		return fmt.Errorf("查找邀请关系失败: %v", err)
	}
	
	if relation == nil {
		// 用户没有邀请关系，跳过奖励处理
		return nil
	}
	
	// 发放积分奖励给邀请人
	pointsService := s.getPointsService()
	description := fmt.Sprintf("邀请用户注册并付费奖励，被邀请用户ID: %d", inviteeID)
	
	if err := pointsService.EarnPoints(relation.InviterID, rewardPoints, model.SourceTypeInvite, description); err != nil {
		return fmt.Errorf("发放邀请奖励失败: %v", err)
	}
	
	
	return nil
}

// GetInviterStatistics 获取邀请人的统计信息
func (s *InviteService) GetInviterStatistics(inviterID int) (map[string]interface{}, error) {
	if inviterID <= 0 {
		return nil, fmt.Errorf("无效的邀请人ID")
	}
	
	// 获取邀请码信息
	userInviteCodeDao := s.getUserInviteCodeDao()
	inviteCodeInfo, err := userInviteCodeDao.GetUserInviteCode(inviterID)
	if err != nil {
		return nil, fmt.Errorf("获取邀请码信息失败: %v", err)
	}
	
	// 获取邀请关系统计
	inviteRelationDao := s.getInviteRelationDao()
	relationStats, err := inviteRelationDao.GetInviterStatistics(inviterID)
	if err != nil {
		return nil, fmt.Errorf("获取邀请关系统计失败: %v", err)
	}
	
	// 获取邀请获得的积分
	pointsService := s.getPointsService()
	invitePoints, err := pointsService.GetUserEarnPointsBySource(inviterID, model.SourceTypeInvite)
	if err != nil {
		return nil, fmt.Errorf("获取邀请积分统计失败: %v", err)
	}
	
	result := map[string]interface{}{
		"invite_code":                 inviteCodeInfo.InviteCode,
		"invite_points_from_records":  invitePoints, // 从积分记录中统计的邀请积分
	}
	
	// 合并邀请关系统计
	for k, v := range relationStats {
		result[k] = v
	}
	
	return result, nil
}

// GetInviteRelationsByInviter 获取邀请人的邀请关系列表
func (s *InviteService) GetInviteRelationsByInviter(inviterID int, page, pageSize int) ([]model.InviteRelation, int64, error) {
	if inviterID <= 0 {
		return nil, 0, fmt.Errorf("无效的邀请人ID")
	}
	
	inviteRelationDao := s.getInviteRelationDao()
	return inviteRelationDao.GetInviteRelationsByInviter(inviterID, page, pageSize)
}



// GetInviteStatistics 获取整体邀请统计
func (s *InviteService) GetInviteStatistics() (map[string]interface{}, error) {
	userInviteCodeDao := s.getUserInviteCodeDao()
	inviteRelationDao := s.getInviteRelationDao()
	
	// 获取邀请码统计
	codeStats, err := userInviteCodeDao.GetInviteCodeStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取邀请码统计失败: %v", err)
	}
	
	// 获取邀请关系统计
	relationStats, err := inviteRelationDao.GetInviteStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取邀请关系统计失败: %v", err)
	}
	
	// 合并统计信息
	result := make(map[string]interface{})
	for k, v := range codeStats {
		result[k] = v
	}
	for k, v := range relationStats {
		result[k] = v
	}
	
	return result, nil
}

// GetTopInviters 获取邀请排行榜
func (s *InviteService) GetTopInviters(limit int) ([]map[string]interface{}, error) {
	inviteRelationDao := s.getInviteRelationDao()
	return inviteRelationDao.GetTopInviters(limit)
}

// IsUserInvited 检查用户是否被邀请
func (s *InviteService) IsUserInvited(userID int) (bool, *model.InviteRelation, error) {
	if userID <= 0 {
		return false, nil, fmt.Errorf("无效的用户ID")
	}
	
	inviteRelationDao := s.getInviteRelationDao()
	relation, err := inviteRelationDao.GetInviteRelationByInvitee(userID)
	if err != nil {
		return false, nil, fmt.Errorf("检查邀请关系失败: %v", err)
	}
	
	if relation == nil {
		return false, nil, nil
	}
	
	return true, relation, nil
}

// ValidateInviteOperation 验证邀请操作的合法性
func (s *InviteService) ValidateInviteOperation(inviterID, inviteeID int, inviteCode string) error {
	if inviterID <= 0 || inviteeID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if inviterID == inviteeID {
		return fmt.Errorf("不能邀请自己")
	}
	
	// 验证邀请码
	if _, err := s.ValidateInviteCode(inviteCode); err != nil {
		return err
	}
	
	// 检查被邀请人是否已被邀请
	isInvited, _, err := s.IsUserInvited(inviteeID)
	if err != nil {
		return err
	}
	
	if isInvited {
		return fmt.Errorf("该用户已被邀请")
	}
	
	return nil
}