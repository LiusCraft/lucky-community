package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type PointRecordDao struct{}

// CreatePointRecord 创建积分记录
func (d *PointRecordDao) CreatePointRecord(record *model.PointRecord) error {
	err := model.PointRecordModel().Create(record).Error
	if err != nil {
		return fmt.Errorf("创建积分记录失败: %v", err)
	}
	return nil
}

// CreatePointRecordWithTransaction 在事务中创建积分记录
func (d *PointRecordDao) CreatePointRecordWithTransaction(tx *gorm.DB, record *model.PointRecord) error {
	err := tx.Model(&model.PointRecord{}).Create(record).Error
	if err != nil {
		return fmt.Errorf("创建积分记录失败: %v", err)
	}
	return nil
}

// GetUserPointRecords 获取用户积分记录列表（分页）
func (d *PointRecordDao) GetUserPointRecords(userID int, page, pageSize int, recordType string) ([]model.PointRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := model.PointRecordModel().Where("user_id = ?", userID)
	
	// 根据类型筛选
	if recordType != "" && (recordType == model.PointTypeEarn || recordType == model.PointTypeSpend) {
		query = query.Where("type = ?", recordType)
	}

	var total int64
	var records []model.PointRecord

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询积分记录总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户积分记录失败: %v", err)
	}

	return records, total, nil
}

// GetPointRecordsBySourceType 根据来源类型查询积分记录
func (d *PointRecordDao) GetPointRecordsBySourceType(sourceType string, page, pageSize int) ([]model.PointRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	var records []model.PointRecord

	query := model.PointRecordModel().Where("source_type = ?", sourceType)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询积分记录总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询积分记录失败: %v", err)
	}

	return records, total, nil
}

// GetPointRecordsByExchangeRequest 根据兑换申请ID查询相关积分记录
func (d *PointRecordDao) GetPointRecordsByExchangeRequest(exchangeRequestID int64) ([]model.PointRecord, error) {
	var records []model.PointRecord
	err := model.PointRecordModel().
		Where("exchange_request_id = ?", exchangeRequestID).
		Order("created_at DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("查询兑换相关积分记录失败: %v", err)
	}

	return records, nil
}

// GetPointRecordsStatistics 获取积分记录统计信息
func (d *PointRecordDao) GetPointRecordsStatistics() (map[string]interface{}, error) {
	var earnStats struct {
		TotalEarnRecords int64 `json:"total_earn_records"`
		TotalEarnPoints  int64 `json:"total_earn_points"`
	}

	var spendStats struct {
		TotalSpendRecords int64 `json:"total_spend_records"`
		TotalSpendPoints  int64 `json:"total_spend_points"`
	}

	// 获取积分获得统计
	if err := model.PointRecordModel().
		Where("type = ?", model.PointTypeEarn).
		Select("COUNT(*) as total_earn_records, SUM(points) as total_earn_points").
		Scan(&earnStats).Error; err != nil {
		return nil, fmt.Errorf("查询积分获得统计失败: %v", err)
	}

	// 获取积分消费统计
	if err := model.PointRecordModel().
		Where("type = ?", model.PointTypeSpend).
		Select("COUNT(*) as total_spend_records, SUM(points) as total_spend_points").
		Scan(&spendStats).Error; err != nil {
		return nil, fmt.Errorf("查询积分消费统计失败: %v", err)
	}

	result := map[string]interface{}{
		"total_earn_records": earnStats.TotalEarnRecords,
		"total_earn_points":  earnStats.TotalEarnPoints,
		"total_spend_records": spendStats.TotalSpendRecords,
		"total_spend_points":  spendStats.TotalSpendPoints,
	}

	return result, nil
}

// GetSourceTypeStatistics 获取按来源类型的积分统计
func (d *PointRecordDao) GetSourceTypeStatistics() ([]map[string]interface{}, error) {
	var results []struct {
		SourceType   string `json:"source_type"`
		RecordCount  int64  `json:"record_count"`
		TotalPoints  int64  `json:"total_points"`
	}

	err := model.PointRecordModel().
		Where("type = ? AND source_type IS NOT NULL", model.PointTypeEarn).
		Select("source_type, COUNT(*) as record_count, SUM(points) as total_points").
		Group("source_type").
		Order("total_points DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询来源类型统计失败: %v", err)
	}

	// 转换为 map 切片
	var statistics []map[string]interface{}
	for _, result := range results {
		statistics = append(statistics, map[string]interface{}{
			"source_type":   result.SourceType,
			"record_count":  result.RecordCount,
			"total_points":  result.TotalPoints,
		})
	}

	return statistics, nil
}

// GetRewardTypeStatistics 获取按奖励类型的积分统计
func (d *PointRecordDao) GetRewardTypeStatistics() ([]map[string]interface{}, error) {
	var results []struct {
		RewardType   string `json:"reward_type"`
		RecordCount  int64  `json:"record_count"`
		TotalPoints  int64  `json:"total_points"`
	}

	err := model.PointRecordModel().
		Where("type = ? AND reward_type IS NOT NULL", model.PointTypeSpend).
		Select("reward_type, COUNT(*) as record_count, SUM(points) as total_points").
		Group("reward_type").
		Order("total_points DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询奖励类型统计失败: %v", err)
	}

	// 转换为 map 切片
	var statistics []map[string]interface{}
	for _, result := range results {
		statistics = append(statistics, map[string]interface{}{
			"reward_type":   result.RewardType,
			"record_count":  result.RecordCount,
			"total_points":  result.TotalPoints,
		})
	}

	return statistics, nil
}

// GetUserEarnPointsBySource 获取用户特定来源的积分总数
func (d *PointRecordDao) GetUserEarnPointsBySource(userID int, sourceType string) (int64, error) {
	var totalPoints int64
	err := model.PointRecordModel().
		Where("user_id = ? AND type = ? AND source_type = ?", userID, model.PointTypeEarn, sourceType).
		Select("COALESCE(SUM(points), 0)").
		Scan(&totalPoints).Error

	if err != nil {
		return 0, fmt.Errorf("查询用户特定来源积分失败: %v", err)
	}

	return totalPoints, nil
}

// DeletePointRecordsByExchangeRequest 删除兑换申请相关的积分记录（用于取消兑换）
func (d *PointRecordDao) DeletePointRecordsByExchangeRequest(exchangeRequestID int64) error {
	err := model.PointRecordModel().
		Where("exchange_request_id = ?", exchangeRequestID).
		Delete(&model.PointRecord{}).Error

	if err != nil {
		return fmt.Errorf("删除兑换相关积分记录失败: %v", err)
	}

	return nil
}

// GetPointsSumByType 获取指定类型积分总数
func (d *PointRecordDao) GetPointsSumByType(recordType, dateStart, dateEnd string) (int64, error) {
	query := model.PointRecordModel().Where("type = ?", recordType)
	
	// 按日期范围筛选
	if dateStart != "" {
		query = query.Where("created_at >= ?", dateStart)
	}
	if dateEnd != "" {
		query = query.Where("created_at <= ?", dateEnd + " 23:59:59")
	}
	
	var total int64
	err := query.Select("COALESCE(SUM(points), 0)").Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("查询积分总数失败: %v", err)
	}
	
	return total, nil
}

// GetActiveUsersCount 获取活跃用户数
func (d *PointRecordDao) GetActiveUsersCount(dateStart, dateEnd string) (int64, error) {
	query := model.PointRecordModel()
	
	// 按日期范围筛选
	if dateStart != "" {
		query = query.Where("created_at >= ?", dateStart)
	}
	if dateEnd != "" {
		query = query.Where("created_at <= ?", dateEnd + " 23:59:59")
	}
	
	var count int64
	err := query.Distinct("user_id").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("查询活跃用户数失败: %v", err)
	}
	
	return count, nil
}

// GetAllPointRecords 获取所有积分记录（管理员使用）
func (d *PointRecordDao) GetAllPointRecords(page, pageSize, userID int, recordType, dateStart, dateEnd string) ([]model.PointRecord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	
	query := model.PointRecordModel()
	
	// 按用户ID筛选
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	
	// 按类型筛选
	if recordType != "" && (recordType == model.PointTypeEarn || recordType == model.PointTypeSpend) {
		query = query.Where("type = ?", recordType)
	}
	
	// 按日期范围筛选
	if dateStart != "" {
		query = query.Where("created_at >= ?", dateStart)
	}
	if dateEnd != "" {
		query = query.Where("created_at <= ?", dateEnd + " 23:59:59")
	}
	
	var total int64
	var records []model.PointRecord
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询积分记录总数失败: %v", err)
	}
	
	// 获取列表数据（按创建时间降序）
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询积分记录列表失败: %v", err)
	}
	
	return records, total, nil
}