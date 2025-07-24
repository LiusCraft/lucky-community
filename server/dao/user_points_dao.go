package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type UserPointsDao struct{}

// GetUserPoints 获取用户积分信息，如果不存在则创建
func (d *UserPointsDao) GetUserPoints(userID int) (*model.UserPoints, error) {
	var userPoints model.UserPoints
	err := model.UserPoint().Where("user_id = ?", userID).First(&userPoints).Error
	
	if err == gorm.ErrRecordNotFound {
		// 用户积分记录不存在，创建新记录
		userPoints = model.UserPoints{
			UserID:          userID,
			TotalEarned:     0,
			AvailablePoints: 0,
			TotalSpent:      0,
		}
		if createErr := model.UserPoint().Create(&userPoints).Error; createErr != nil {
			return nil, fmt.Errorf("创建用户积分记录失败: %v", createErr)
		}
		return &userPoints, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询用户积分失败: %v", err)
	}
	
	return &userPoints, nil
}

// UpdateUserPoints 更新用户积分
func (d *UserPointsDao) UpdateUserPoints(userPoints *model.UserPoints) error {
	err := model.UserPoint().Save(userPoints).Error
	if err != nil {
		return fmt.Errorf("更新用户积分失败: %v", err)
	}
	return nil
}

// AddPointsWithTransaction 在事务中增加用户积分
func (d *UserPointsDao) AddPointsWithTransaction(tx *gorm.DB, userID int, points int) (*model.UserPoints, error) {
	// 查询用户积分记录
	var userPoints model.UserPoints
	err := tx.Model(&model.UserPoints{}).Where("user_id = ?", userID).First(&userPoints).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新的积分记录
		userPoints = model.UserPoints{
			UserID:          userID,
			TotalEarned:     points,
			AvailablePoints: points,
			TotalSpent:      0,
		}
		if createErr := tx.Create(&userPoints).Error; createErr != nil {
			return nil, fmt.Errorf("创建用户积分记录失败: %v", createErr)
		}
	} else if err != nil {
		return nil, fmt.Errorf("查询用户积分失败: %v", err)
	} else {
		// 更新积分
		userPoints.TotalEarned += points
		userPoints.AvailablePoints += points
		
		if updateErr := tx.Save(&userPoints).Error; updateErr != nil {
			return nil, fmt.Errorf("更新用户积分失败: %v", updateErr)
		}
	}
	
	return &userPoints, nil
}

// DeductPointsWithTransaction 在事务中扣减用户积分
func (d *UserPointsDao) DeductPointsWithTransaction(tx *gorm.DB, userID int, points int) (*model.UserPoints, error) {
	// 查询用户积分记录
	var userPoints model.UserPoints
	err := tx.Model(&model.UserPoints{}).Where("user_id = ?", userID).First(&userPoints).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("用户积分记录不存在")
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询用户积分失败: %v", err)
	}
	
	// 检查积分是否充足
	if userPoints.AvailablePoints < points {
		return nil, fmt.Errorf("积分不足：需要%d积分，当前可用%d积分", points, userPoints.AvailablePoints)
	}
	
	// 扣减积分
	userPoints.AvailablePoints -= points
	userPoints.TotalSpent += points
	
	if updateErr := tx.Save(&userPoints).Error; updateErr != nil {
		return nil, fmt.Errorf("扣减用户积分失败: %v", updateErr)
	}
	
	return &userPoints, nil
}

// GetUsersPointsRanking 获取用户积分排行榜
func (d *UserPointsDao) GetUsersPointsRanking(limit int) ([]model.UserPoints, error) {
	if limit <= 0 {
		limit = 10
	}
	
	var rankings []model.UserPoints
	err := model.UserPoint().
		Order("total_earned DESC").
		Limit(limit).
		Find(&rankings).Error
	
	if err != nil {
		return nil, fmt.Errorf("查询积分排行榜失败: %v", err)
	}
	
	return rankings, nil
}

// GetUsersPointsWithPagination 分页获取用户积分列表
func (d *UserPointsDao) GetUsersPointsWithPagination(page, pageSize int, orderBy string) ([]model.UserPoints, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	
	// 验证排序字段
	validOrderFields := map[string]bool{
		"total_earned":     true,
		"available_points": true,
		"total_spent":      true,
		"created_at":       true,
		"updated_at":       true,
	}
	
	if orderBy == "" || !validOrderFields[orderBy] {
		orderBy = "total_earned DESC"
	} else {
		orderBy = orderBy + " DESC"
	}
	
	var total int64
	var userPoints []model.UserPoints
	
	// 获取总数
	if err := model.UserPoint().Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %v", err)
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := model.UserPoint().
		Order(orderBy).
		Offset(offset).
		Limit(pageSize).
		Find(&userPoints).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("分页查询用户积分失败: %v", err)
	}
	
	return userPoints, total, nil
}

// GetPointsStatistics 获取积分系统统计信息
func (d *UserPointsDao) GetPointsStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalUsers         int64 `json:"total_users"`
		TotalPointsEarned  int64 `json:"total_points_earned"`
		TotalPointsSpent   int64 `json:"total_points_spent"`
		ActivePointsAmount int64 `json:"active_points_amount"`
	}
	
	// 获取用户总数
	if err := model.UserPoint().Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %v", err)
	}
	
	// 获取积分统计
	row := model.UserPoint().
		Select("SUM(total_earned) as total_points_earned, SUM(total_spent) as total_points_spent, SUM(available_points) as active_points_amount").
		Row()
	
	if err := row.Scan(&stats.TotalPointsEarned, &stats.TotalPointsSpent, &stats.ActivePointsAmount); err != nil {
		return nil, fmt.Errorf("查询积分统计失败: %v", err)
	}
	
	result := map[string]interface{}{
		"total_users":           stats.TotalUsers,
		"total_points_earned":   stats.TotalPointsEarned,
		"total_points_spent":    stats.TotalPointsSpent,
		"active_points_amount":  stats.ActivePointsAmount,
		"average_points_per_user": func() float64 {
			if stats.TotalUsers > 0 {
				return float64(stats.TotalPointsEarned) / float64(stats.TotalUsers)
			}
			return 0
		}(),
	}
	
	return result, nil
}

// BatchCreateUserPoints 批量创建用户积分记录
func (d *UserPointsDao) BatchCreateUserPoints(userIDs []int) error {
	var userPoints []model.UserPoints
	
	for _, userID := range userIDs {
		userPoints = append(userPoints, model.UserPoints{
			UserID:          userID,
			TotalEarned:     0,
			AvailablePoints: 0,
			TotalSpent:      0,
		})
	}
	
	if err := model.UserPoint().CreateInBatches(userPoints, 100).Error; err != nil {
		return fmt.Errorf("批量创建用户积分记录失败: %v", err)
	}
	
	return nil
}

// GetTotalUsersCount 获取用户总数
func (d *UserPointsDao) GetTotalUsersCount() (int64, error) {
	var count int64
	err := model.UserPoint().Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("查询用户总数失败: %v", err)
	}
	return count, nil
}

// GetAveragePointsBalance 获取平均积分余额
func (d *UserPointsDao) GetAveragePointsBalance() (float64, error) {
	var avgBalance float64
	err := model.UserPoint().Select("COALESCE(AVG(available_points), 0)").Row().Scan(&avgBalance)
	if err != nil {
		return 0, fmt.Errorf("查询平均积分余额失败: %v", err)
	}
	return avgBalance, nil
}