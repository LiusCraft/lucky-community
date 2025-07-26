package services

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type PointsService struct{}

// getUserPointsDao 获取用户积分DAO实例
func (s *PointsService) getUserPointsDao() *dao.UserPointsDao {
	return &dao.UserPointsDao{}
}

// getPointRecordDao 获取积分记录DAO实例
func (s *PointsService) getPointRecordDao() *dao.PointRecordDao {
	return &dao.PointRecordDao{}
}

// GetUserPoints 获取用户积分信息
func (s *PointsService) GetUserPoints(userID int) (*model.UserPoints, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}
	
	userPointsDao := s.getUserPointsDao()
	return userPointsDao.GetUserPoints(userID)
}

// EarnPoints 用户获得积分
func (s *PointsService) EarnPoints(userID int, points int, sourceType string, description string) error {
	if userID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if points <= 0 {
		return fmt.Errorf("积分数量必须大于0")
	}
	
	if sourceType == "" {
		return fmt.Errorf("积分来源类型不能为空")
	}
	
	// 验证来源类型是否有效
	validSourceTypes := map[string]bool{
		model.SourceTypeInvite:   true,
		model.SourceTypeCourse:   true,
		model.SourceTypeContent:  true,
		model.SourceTypeDaily:    true,
		model.SourceTypeActivity: true,
		model.SourceTypeManual:   true,
		model.SourceTypeOther:    true,
	}
	
	if !validSourceTypes[sourceType] {
		return fmt.Errorf("无效的积分来源类型: %s", sourceType)
	}
	
	db := mysql.GetInstance()
	return db.Transaction(func(tx *gorm.DB) error {
		userPointsDao := s.getUserPointsDao()
		pointRecordDao := s.getPointRecordDao()
		
		// 增加用户积分
		_, err := userPointsDao.AddPointsWithTransaction(tx, userID, points)
		if err != nil {
			return err
		}
		
		// 创建积分记录
		record := &model.PointRecord{
			UserID:      userID,
			Type:        model.PointTypeEarn,
			Points:      points,
			SourceType:  &sourceType,
			Description: description,
		}
		
		return pointRecordDao.CreatePointRecordWithTransaction(tx, record)
	})
}

// SpendPoints 用户消费积分
func (s *PointsService) SpendPoints(userID int, points int, rewardType string, description string, exchangeRequestID *int64) error {
	if userID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if points <= 0 {
		return fmt.Errorf("积分数量必须大于0")
	}
	
	if rewardType == "" {
		return fmt.Errorf("兑换奖励类型不能为空")
	}
	
	// 验证奖励类型是否有效
	validRewardTypes := map[string]bool{
		model.RewardTypeCash:      true,
		model.RewardTypeService:   true,
		model.RewardTypeProduct:   true,
		model.RewardTypePrivilege: true,
		model.RewardTypeCoupon:    true,
		model.RewardTypeManual:    true,
		model.RewardTypeOther:     true,
	}
	
	if !validRewardTypes[rewardType] {
		return fmt.Errorf("无效的兑换奖励类型: %s", rewardType)
	}
	
	db := mysql.GetInstance()
	return db.Transaction(func(tx *gorm.DB) error {
		userPointsDao := s.getUserPointsDao()
		pointRecordDao := s.getPointRecordDao()
		
		// 扣减用户积分
		_, err := userPointsDao.DeductPointsWithTransaction(tx, userID, points)
		if err != nil {
			return err
		}
		
		// 创建积分记录
		record := &model.PointRecord{
			UserID:            userID,
			Type:              model.PointTypeSpend,
			Points:            points,
			RewardType:        &rewardType,
			Description:       description,
			ExchangeRequestID: exchangeRequestID,
		}
		
		return pointRecordDao.CreatePointRecordWithTransaction(tx, record)
	})
}

// GetUserPointRecords 获取用户积分记录
func (s *PointsService) GetUserPointRecords(userID int, page, pageSize int, recordType string) ([]model.PointRecord, int64, error) {
	if userID <= 0 {
		return nil, 0, fmt.Errorf("无效的用户ID")
	}
	
	pointRecordDao := s.getPointRecordDao()
	return pointRecordDao.GetUserPointRecords(userID, page, pageSize, recordType)
}

// GetPointsStatistics 获取积分系统统计信息
func (s *PointsService) GetPointsStatistics() (map[string]interface{}, error) {
	userPointsDao := s.getUserPointsDao()
	pointRecordDao := s.getPointRecordDao()
	
	// 获取用户积分统计
	userStats, err := userPointsDao.GetPointsStatistics()
	if err != nil {
		return nil, err
	}
	
	// 获取积分记录统计
	recordStats, err := pointRecordDao.GetPointRecordsStatistics()
	if err != nil {
		return nil, err
	}
	
	// 获取来源类型统计
	sourceStats, err := pointRecordDao.GetSourceTypeStatistics()
	if err != nil {
		return nil, err
	}
	
	// 获取奖励类型统计
	rewardStats, err := pointRecordDao.GetRewardTypeStatistics()
	if err != nil {
		return nil, err
	}
	
	// 合并统计信息
	result := make(map[string]interface{})
	for k, v := range userStats {
		result[k] = v
	}
	for k, v := range recordStats {
		result[k] = v
	}
	result["source_type_statistics"] = sourceStats
	result["reward_type_statistics"] = rewardStats
	
	return result, nil
}

// GetUsersPointsRanking 获取用户积分排行榜
func (s *PointsService) GetUsersPointsRanking(limit int) ([]model.UserPoints, error) {
	userPointsDao := s.getUserPointsDao()
	return userPointsDao.GetUsersPointsRanking(limit)
}

// GetUserEarnPointsBySource 获取用户特定来源的积分总数
func (s *PointsService) GetUserEarnPointsBySource(userID int, sourceType string) (int64, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("无效的用户ID")
	}
	
	pointRecordDao := s.getPointRecordDao()
	return pointRecordDao.GetUserEarnPointsBySource(userID, sourceType)
}

// TransferPoints 积分转账（预留功能）
func (s *PointsService) TransferPoints(fromUserID, toUserID int, points int, description string) error {
	if fromUserID <= 0 || toUserID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if fromUserID == toUserID {
		return fmt.Errorf("不能给自己转账")
	}
	
	if points <= 0 {
		return fmt.Errorf("转账积分必须大于0")
	}
	
	db := mysql.GetInstance()
	return db.Transaction(func(tx *gorm.DB) error {
		userPointsDao := s.getUserPointsDao()
		pointRecordDao := s.getPointRecordDao()
		
		// 扣减转出用户积分
		_, err := userPointsDao.DeductPointsWithTransaction(tx, fromUserID, points)
		if err != nil {
			return err
		}
		
		// 增加转入用户积分
		_, err = userPointsDao.AddPointsWithTransaction(tx, toUserID, points)
		if err != nil {
			return err
		}
		
		// 创建转出记录
		transferOutType := model.RewardTypeOther
		outRecord := &model.PointRecord{
			UserID:      fromUserID,
			Type:        model.PointTypeSpend,
			Points:      points,
			RewardType:  &transferOutType,
			Description: fmt.Sprintf("转账给用户%d: %s", toUserID, description),
		}
		
		if err := pointRecordDao.CreatePointRecordWithTransaction(tx, outRecord); err != nil {
			return err
		}
		
		// 创建转入记录
		transferInType := model.SourceTypeOther
		inRecord := &model.PointRecord{
			UserID:      toUserID,
			Type:        model.PointTypeEarn,
			Points:      points,
			SourceType:  &transferInType,
			Description: fmt.Sprintf("用户%d转账: %s", fromUserID, description),
		}
		
		return pointRecordDao.CreatePointRecordWithTransaction(tx, inRecord)
	})
}

// ValidatePointsOperation 验证积分操作的合法性
func (s *PointsService) ValidatePointsOperation(userID int, points int, operationType string) error {
	if userID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if points <= 0 {
		return fmt.Errorf("积分数量必须大于0")
	}
	
	// 如果是消费操作，检查用户积分是否充足
	if operationType == model.PointTypeSpend {
		userPoints, err := s.GetUserPoints(userID)
		if err != nil {
			return fmt.Errorf("获取用户积分失败: %v", err)
		}
		
		if !userPoints.HasSufficientPoints(points) {
			return fmt.Errorf("积分不足：需要%d积分，当前可用%d积分", points, userPoints.AvailablePoints)
		}
	}
	
	return nil
}

// BatchEarnPoints 批量发放积分
func (s *PointsService) BatchEarnPoints(operations []struct {
	UserID      int
	Points      int
	SourceType  string
	Description string
}) error {
	if len(operations) == 0 {
		return fmt.Errorf("操作列表不能为空")
	}
	
	db := mysql.GetInstance()
	return db.Transaction(func(tx *gorm.DB) error {
		userPointsDao := s.getUserPointsDao()
		pointRecordDao := s.getPointRecordDao()
		
		for _, op := range operations {
			// 验证参数
			if err := s.ValidatePointsOperation(op.UserID, op.Points, model.PointTypeEarn); err != nil {
				return fmt.Errorf("用户%d操作验证失败: %v", op.UserID, err)
			}
			
			// 增加用户积分
			_, err := userPointsDao.AddPointsWithTransaction(tx, op.UserID, op.Points)
			if err != nil {
				return fmt.Errorf("用户%d积分增加失败: %v", op.UserID, err)
			}
			
			// 创建积分记录
			record := &model.PointRecord{
				UserID:      op.UserID,
				Type:        model.PointTypeEarn,
				Points:      op.Points,
				SourceType:  &op.SourceType,
				Description: op.Description,
			}
			
			if err := pointRecordDao.CreatePointRecordWithTransaction(tx, record); err != nil {
				return fmt.Errorf("用户%d积分记录创建失败: %v", op.UserID, err)
			}
		}
		
		return nil
	})
}

// GetSystemPointsStatistics 获取积分系统统计数据
func (s *PointsService) GetSystemPointsStatistics(dateStart, dateEnd string) (map[string]interface{}, error) {
	pointRecordDao := s.getPointRecordDao()
	
	// 获取总积分发放和消费
	totalEarn, err := pointRecordDao.GetPointsSumByType("earn", dateStart, dateEnd)
	if err != nil {
		return nil, fmt.Errorf("获取积分发放统计失败: %v", err)
	}
	
	totalSpend, err := pointRecordDao.GetPointsSumByType("spend", dateStart, dateEnd)
	if err != nil {
		return nil, fmt.Errorf("获取积分消费统计失败: %v", err)
	}
	
	statistics := map[string]interface{}{
		"total_earn":  totalEarn,
		"total_spend": totalSpend,
	}
	
	return statistics, nil
}

// GetAllPointRecords 获取所有积分记录（管理员使用）
func (s *PointsService) GetAllPointRecords(page, pageSize, userID int, recordType, dateStart, dateEnd string) ([]dao.PointRecordWithUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	
	pointRecordDao := s.getPointRecordDao()
	return pointRecordDao.GetAllPointRecords(page, pageSize, userID, recordType, dateStart, dateEnd)
}

// ManualGrantPoints 管理员手动发放积分
func (s *PointsService) ManualGrantPoints(userID int, points int, description string) error {
	if userID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if points <= 0 {
		return fmt.Errorf("发放积分数量必须大于0")
	}
	
	if description == "" {
		return fmt.Errorf("发放原因不能为空")
	}
	
	// 使用手动发放类型
	return s.EarnPoints(userID, points, model.SourceTypeManual, description)
}