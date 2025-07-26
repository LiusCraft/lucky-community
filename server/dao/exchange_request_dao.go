package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type ExchangeRequestDao struct{}

// ExchangeRequestWithUser 包含用户信息的兑换申请
type ExchangeRequestWithUser struct {
	model.ExchangeRequest
	UserName    string `json:"user_name" gorm:"column:user_name"`
	ProductName string `json:"product_name" gorm:"column:product_name"`
}

// CreateExchangeRequest 创建兑换申请
func (d *ExchangeRequestDao) CreateExchangeRequest(request *model.ExchangeRequest) error {
	err := model.ExchangeRequestModel().Create(request).Error
	if err != nil {
		return fmt.Errorf("创建兑换申请失败: %v", err)
	}
	return nil
}

// CreateExchangeRequestWithTransaction 在事务中创建兑换申请
func (d *ExchangeRequestDao) CreateExchangeRequestWithTransaction(tx *gorm.DB, request *model.ExchangeRequest) error {
	err := tx.Model(&model.ExchangeRequest{}).Create(request).Error
	if err != nil {
		return fmt.Errorf("创建兑换申请失败: %v", err)
	}
	return nil
}

// GetExchangeRequestByID 根据ID获取兑换申请
func (d *ExchangeRequestDao) GetExchangeRequestByID(id int64) (*model.ExchangeRequest, error) {
	var request model.ExchangeRequest
	err := model.ExchangeRequestModel().Where("id = ?", id).First(&request).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("兑换申请不存在")
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询兑换申请失败: %v", err)
	}
	
	return &request, nil
}

// GetUserExchangeRequests 获取用户的兑换申请列表
func (d *ExchangeRequestDao) GetUserExchangeRequests(userID int, page, pageSize int, status string) ([]model.ExchangeRequest, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := model.ExchangeRequestModel().Where("user_id = ?", userID)
	
	// 根据状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var requests []model.ExchangeRequest

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询兑换申请总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询用户兑换申请失败: %v", err)
	}

	return requests, total, nil
}

// GetExchangeRequestsWithPagination 分页获取所有兑换申请（管理员用）
func (d *ExchangeRequestDao) GetExchangeRequestsWithPagination(page, pageSize int, status string) ([]ExchangeRequestWithUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := model.ExchangeRequestModel().
		Select("exchange_requests.*, users.name as user_name, point_products.name as product_name").
		Joins("LEFT JOIN users ON exchange_requests.user_id = users.id").
		Joins("LEFT JOIN point_products ON exchange_requests.product_id = point_products.id")
	
	// 根据状态筛选
	if status != "" {
		query = query.Where("exchange_requests.status = ?", status)
	}

	var total int64
	var requests []ExchangeRequestWithUser

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询兑换申请总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("exchange_requests.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	if err != nil {
		return nil, 0, fmt.Errorf("分页查询兑换申请失败: %v", err)
	}

	return requests, total, nil
}

// UpdateExchangeRequestStatus 更新兑换申请状态
func (d *ExchangeRequestDao) UpdateExchangeRequestStatus(id int64, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	// 如果状态为已处理，设置处理时间
	if status == model.ExchangeStatusProcessed {
		updates["processed_at"] = gorm.Expr("NOW()")
	}

	err := model.ExchangeRequestModel().Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新兑换申请状态失败: %v", err)
	}
	
	return nil
}

// UpdateExchangeRequest 更新兑换申请
func (d *ExchangeRequestDao) UpdateExchangeRequest(request *model.ExchangeRequest) error {
	err := model.ExchangeRequestModel().Save(request).Error
	if err != nil {
		return fmt.Errorf("更新兑换申请失败: %v", err)
	}
	return nil
}

// GetExchangeRequestStatistics 获取兑换申请统计信息
func (d *ExchangeRequestDao) GetExchangeRequestStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalRequests     int64 `json:"total_requests"`
		PendingRequests   int64 `json:"pending_requests"`
		ProcessedRequests int64 `json:"processed_requests"`
		TotalPointsSpent  int64 `json:"total_points_spent"`
	}

	// 获取各状态的申请数量
	if err := model.ExchangeRequestModel().Count(&stats.TotalRequests).Error; err != nil {
		return nil, fmt.Errorf("查询兑换申请总数失败: %v", err)
	}

	if err := model.ExchangeRequestModel().Where("status = ?", model.ExchangeStatusPending).Count(&stats.PendingRequests).Error; err != nil {
		return nil, fmt.Errorf("查询待处理申请数失败: %v", err)
	}

	if err := model.ExchangeRequestModel().Where("status = ?", model.ExchangeStatusProcessed).Count(&stats.ProcessedRequests).Error; err != nil {
		return nil, fmt.Errorf("查询已处理申请数失败: %v", err)
	}

	// 获取总消费积分
	row := model.ExchangeRequestModel().
		Where("status = ?", model.ExchangeStatusProcessed).
		Select("SUM(points_cost * quantity)").
		Row()
	
	if err := row.Scan(&stats.TotalPointsSpent); err != nil {
		return nil, fmt.Errorf("查询总消费积分失败: %v", err)
	}

	result := map[string]interface{}{
		"total_requests":     stats.TotalRequests,
		"pending_requests":   stats.PendingRequests,
		"processed_requests": stats.ProcessedRequests,
		"total_points_spent": stats.TotalPointsSpent,
		"completion_rate": func() float64 {
			if stats.TotalRequests > 0 {
				return float64(stats.ProcessedRequests) / float64(stats.TotalRequests) * 100
			}
			return 0
		}(),
	}

	return result, nil
}

// GetExchangeRequestsByProduct 根据商品获取兑换申请统计
func (d *ExchangeRequestDao) GetExchangeRequestsByProduct() ([]map[string]interface{}, error) {
	var results []struct {
		ProductID     int64 `json:"product_id"`
		ProductName   string `json:"product_name"`
		RequestCount  int64  `json:"request_count"`
		TotalQuantity int64  `json:"total_quantity"`
		TotalPoints   int64  `json:"total_points"`
	}

	err := model.ExchangeRequestModel().
		Select("er.product_id, pp.name as product_name, COUNT(*) as request_count, SUM(er.quantity) as total_quantity, SUM(er.points_cost * er.quantity) as total_points").
		Table("exchange_requests er").
		Joins("LEFT JOIN point_products pp ON er.product_id = pp.id").
		Where("er.status = ?", model.ExchangeStatusProcessed).
		Group("er.product_id, pp.name").
		Order("total_points DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询商品兑换统计失败: %v", err)
	}

	// 转换为 map 切片
	var statistics []map[string]interface{}
	for _, result := range results {
		statistics = append(statistics, map[string]interface{}{
			"product_id":     result.ProductID,
			"product_name":   result.ProductName,
			"request_count":  result.RequestCount,
			"total_quantity": result.TotalQuantity,
			"total_points":   result.TotalPoints,
		})
	}

	return statistics, nil
}

// BatchUpdateExchangeRequestStatus 批量更新兑换申请状态
func (d *ExchangeRequestDao) BatchUpdateExchangeRequestStatus(ids []int64, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	// 如果状态为已处理，设置处理时间
	if status == model.ExchangeStatusProcessed {
		updates["processed_at"] = gorm.Expr("NOW()")
	}

	err := model.ExchangeRequestModel().Where("id IN ?", ids).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("批量更新兑换申请状态失败: %v", err)
	}
	
	return nil
}

// DeleteExchangeRequest 删除兑换申请（管理员操作）
func (d *ExchangeRequestDao) DeleteExchangeRequest(id int64) error {
	err := model.ExchangeRequestModel().Where("id = ?", id).Delete(&model.ExchangeRequest{}).Error
	if err != nil {
		return fmt.Errorf("删除兑换申请失败: %v", err)
	}
	return nil
}