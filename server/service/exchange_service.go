package services

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type ExchangeService struct{}

// getPointProductDao 获取积分商品DAO实例
func (s *ExchangeService) getPointProductDao() *dao.PointProductDao {
	return &dao.PointProductDao{}
}

// getExchangeRequestDao 获取兑换申请DAO实例
func (s *ExchangeService) getExchangeRequestDao() *dao.ExchangeRequestDao {
	return &dao.ExchangeRequestDao{}
}

// getPointsService 获取积分服务实例
func (s *ExchangeService) getPointsService() *PointsService {
	return &PointsService{}
}

// GetActiveProducts 获取所有启用的积分商品
func (s *ExchangeService) GetActiveProducts() ([]model.PointProduct, error) {
	pointProductDao := s.getPointProductDao()
	return pointProductDao.GetActiveProducts()
}

// GetProductByID 根据ID获取积分商品
func (s *ExchangeService) GetProductByID(id int64) (*model.PointProduct, error) {
	if id <= 0 {
		return nil, fmt.Errorf("无效的商品ID")
	}
	
	pointProductDao := s.getPointProductDao()
	return pointProductDao.GetProductByID(id)
}

// GetProductsByRewardType 根据奖励类型获取商品
func (s *ExchangeService) GetProductsByRewardType(rewardType string) ([]model.PointProduct, error) {
	if rewardType == "" {
		return nil, fmt.Errorf("奖励类型不能为空")
	}
	
	pointProductDao := s.getPointProductDao()
	return pointProductDao.GetProductsByRewardType(rewardType)
}

// CreateExchangeRequest 创建兑换申请
func (s *ExchangeService) CreateExchangeRequest(userID int, productID int64, quantity int, userInfo string) (*model.ExchangeRequest, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}
	
	if productID <= 0 {
		return nil, fmt.Errorf("无效的商品ID")
	}
	
	if quantity <= 0 {
		return nil, fmt.Errorf("兑换数量必须大于0")
	}
	
	// 获取商品信息
	product, err := s.GetProductByID(productID)
	if err != nil {
		return nil, fmt.Errorf("获取商品信息失败: %v", err)
	}
	
	// 检查商品是否可购买
	if !product.CanPurchase(quantity) {
		if !product.IsActive() {
			return nil, fmt.Errorf("商品已下架")
		}
		if !product.IsInStock(quantity) {
			return nil, fmt.Errorf("商品库存不足")
		}
	}
	
	// 计算所需积分
	totalPoints := product.PricePoints * quantity
	
	// 验证用户积分是否充足
	pointsService := s.getPointsService()
	if err := pointsService.ValidatePointsOperation(userID, totalPoints, model.PointTypeSpend); err != nil {
		return nil, fmt.Errorf("积分验证失败: %v", err)
	}
	
	// 在事务中创建兑换申请并扣除积分
	db := mysql.GetInstance()
	var exchangeRequest *model.ExchangeRequest
	
	err = db.Transaction(func(tx *gorm.DB) error {
		pointProductDao := s.getPointProductDao()
		exchangeRequestDao := s.getExchangeRequestDao()
		
		// 更新商品库存
		if err := pointProductDao.UpdateProductStockWithTransaction(tx, productID, quantity); err != nil {
			return fmt.Errorf("更新商品库存失败: %v", err)
		}
		
		// 创建兑换申请
		exchangeRequest = &model.ExchangeRequest{
			UserID:     userID,
			ProductID:  productID,
			PointsCost: product.PricePoints,
			Quantity:   quantity,
			Status:     model.ExchangeStatusPending,
			UserInfo:   userInfo,
		}
		
		if err := exchangeRequestDao.CreateExchangeRequestWithTransaction(tx, exchangeRequest); err != nil {
			return fmt.Errorf("创建兑换申请失败: %v", err)
		}
		
		// 扣除用户积分
		description := fmt.Sprintf("兑换商品: %s (数量: %d)", product.Name, quantity)
		if err := pointsService.SpendPoints(userID, totalPoints, product.RewardType, description, &exchangeRequest.ID); err != nil {
			return fmt.Errorf("扣除积分失败: %v", err)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return exchangeRequest, nil
}

// GetUserExchangeRequests 获取用户的兑换申请列表（包含商品名称）
func (s *ExchangeService) GetUserExchangeRequests(userID int, page, pageSize int, status string) ([]dao.ExchangeRequestWithUser, int64, error) {
	if userID <= 0 {
		return nil, 0, fmt.Errorf("无效的用户ID")
	}
	
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.GetUserExchangeRequests(userID, page, pageSize, status)
}

// GetExchangeRequestByID 根据ID获取兑换申请
func (s *ExchangeService) GetExchangeRequestByID(id int64) (*model.ExchangeRequest, error) {
	if id <= 0 {
		return nil, fmt.Errorf("无效的申请ID")
	}
	
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.GetExchangeRequestByID(id)
}

// ProcessExchangeRequest 处理兑换申请（管理员操作）
func (s *ExchangeService) ProcessExchangeRequest(requestID int64) error {
	if requestID <= 0 {
		return fmt.Errorf("无效的申请ID")
	}
	
	// 获取兑换申请
	request, err := s.GetExchangeRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("获取兑换申请失败: %v", err)
	}
	
	// 检查申请状态
	if !request.CanProcess() {
		return fmt.Errorf("申请状态不允许处理，当前状态: %s", request.GetStatusDescription())
	}
	
	// 更新申请状态
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.UpdateExchangeRequestStatus(requestID, model.ExchangeStatusProcessed)
}



// GetExchangeRequestsWithPagination 分页获取所有兑换申请（管理员用）
func (s *ExchangeService) GetExchangeRequestsWithPagination(page, pageSize int, status string) ([]dao.ExchangeRequestWithUser, int64, error) {
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.GetExchangeRequestsWithPagination(page, pageSize, status)
}

// GetExchangeStatistics 获取兑换统计信息
func (s *ExchangeService) GetExchangeStatistics() (map[string]interface{}, error) {
	exchangeRequestDao := s.getExchangeRequestDao()
	pointProductDao := s.getPointProductDao()
	
	// 获取兑换申请统计
	requestStats, err := exchangeRequestDao.GetExchangeRequestStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取兑换申请统计失败: %v", err)
	}
	
	// 获取商品统计
	productStats, err := pointProductDao.GetProductStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取商品统计失败: %v", err)
	}
	
	// 获取按商品的兑换统计
	productExchangeStats, err := exchangeRequestDao.GetExchangeRequestsByProduct()
	if err != nil {
		return nil, fmt.Errorf("获取商品兑换统计失败: %v", err)
	}
	
	// 获取按奖励类型的商品统计
	rewardTypeStats, err := pointProductDao.GetRewardTypeStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取奖励类型统计失败: %v", err)
	}
	
	// 合并统计信息
	result := make(map[string]interface{})
	for k, v := range requestStats {
		result[k] = v
	}
	for k, v := range productStats {
		result[k] = v
	}
	result["product_exchange_statistics"] = productExchangeStats
	result["reward_type_statistics"] = rewardTypeStats
	
	return result, nil
}

// BatchProcessExchangeRequests 批量处理兑换申请
func (s *ExchangeService) BatchProcessExchangeRequests(requestIDs []int64, newStatus string) error {
	if len(requestIDs) == 0 {
		return fmt.Errorf("申请ID列表不能为空")
	}
	
	// 验证状态
	validStatuses := map[string]bool{
		model.ExchangeStatusProcessed: true,
	}
	
	if !validStatuses[newStatus] {
		return fmt.Errorf("无效的状态: %s", newStatus)
	}
	
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.BatchUpdateExchangeRequestStatus(requestIDs, newStatus)
}

// ValidateExchangeOperation 验证兑换操作的合法性
func (s *ExchangeService) ValidateExchangeOperation(userID int, productID int64, quantity int) error {
	if userID <= 0 {
		return fmt.Errorf("无效的用户ID")
	}
	
	if productID <= 0 {
		return fmt.Errorf("无效的商品ID")
	}
	
	if quantity <= 0 {
		return fmt.Errorf("兑换数量必须大于0")
	}
	
	// 获取商品信息
	product, err := s.GetProductByID(productID)
	if err != nil {
		return fmt.Errorf("商品不存在: %v", err)
	}
	
	// 检查商品状态
	if !product.IsActive() {
		return fmt.Errorf("商品已下架")
	}
	
	// 检查库存
	if !product.IsInStock(quantity) {
		return fmt.Errorf("商品库存不足")
	}
	
	// 检查用户积分
	totalPoints := product.PricePoints * quantity
	pointsService := s.getPointsService()
	if err := pointsService.ValidatePointsOperation(userID, totalPoints, model.PointTypeSpend); err != nil {
		return fmt.Errorf("积分不足: %v", err)
	}
	
	return nil
}

// GetAllProducts 获取所有积分商品（包括已禁用的）
func (s *ExchangeService) GetAllProducts() ([]model.PointProduct, error) {
	pointProductDao := s.getPointProductDao()
	return pointProductDao.GetAllProducts()
}

// GetAllExchangeRequests 获取所有兑换申请（管理员使用）
func (s *ExchangeService) GetAllExchangeRequests(page, pageSize int, status string) ([]dao.ExchangeRequestWithUser, int64, error) {
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.GetExchangeRequestsWithPagination(page, pageSize, status)
}

// ProcessExchangeRequestStatus 处理兑换申请（管理员操作）
func (s *ExchangeService) ProcessExchangeRequestStatus(requestID int64, status string) error {
	if requestID <= 0 {
		return fmt.Errorf("无效的申请ID")
	}
	
	// 验证状态
	validStatuses := map[string]bool{
		model.ExchangeStatusProcessed: true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("无效的状态: %s", status)
	}
	
	// 获取兑换申请
	request, err := s.GetExchangeRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("获取兑换申请失败: %v", err)
	}
	
	// 检查申请状态
	if request.Status != model.ExchangeStatusPending {
		return fmt.Errorf("申请状态不允许处理，当前状态: %s", request.Status)
	}
	
	// 更新申请状态
	exchangeRequestDao := s.getExchangeRequestDao()
	return exchangeRequestDao.UpdateExchangeRequestStatus(requestID, status)
}


// CreateProduct 创建积分商品
func (s *ExchangeService) CreateProduct(name, description string, pointsCost int, rewardType, rewardValue string, stock, sortOrder int, isEnabled bool) (*model.PointProduct, error) {
	if name == "" {
		return nil, fmt.Errorf("商品名称不能为空")
	}
	if pointsCost <= 0 {
		return nil, fmt.Errorf("积分成本必须大于0")
	}
	if rewardType == "" {
		return nil, fmt.Errorf("奖励类型不能为空")
	}
	
	status := model.ProductStatusActive
	if !isEnabled {
		status = model.ProductStatusInactive
	}
	
	product := &model.PointProduct{
		Name:        name,
		Description: description,
		PricePoints: pointsCost,
		RewardType:  rewardType,
		Stock:       stock,
		SortOrder:   sortOrder,
		Status:      status,
	}
	
	pointProductDao := s.getPointProductDao()
	if err := pointProductDao.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct 更新积分商品
func (s *ExchangeService) UpdateProduct(productID int64, name, description string, pointsCost int, rewardType, rewardValue string, stock, sortOrder int, isEnabled *bool) error {
	if productID <= 0 {
		return fmt.Errorf("无效的商品ID")
	}
	
	// 直接创建要更新的结构体，前端传什么就是什么
	product := &model.PointProduct{
		ID:          productID,
		Name:        name,
		Description: description,
		PricePoints: pointsCost,
		RewardType:  rewardType,
		Stock:       stock,
		SortOrder:   sortOrder,
	}
	
	if isEnabled != nil {
		if *isEnabled {
			product.Status = model.ProductStatusActive
		} else {
			product.Status = model.ProductStatusInactive
		}
	}
	
	pointProductDao := s.getPointProductDao()
	// 直接用GORM的Updates方法更新结构体
	return pointProductDao.UpdateProduct(product)
}

// DeleteProduct 删除积分商品（物理删除）
func (s *ExchangeService) DeleteProduct(productID int64) error {
	if productID <= 0 {
		return fmt.Errorf("无效的商品ID")
	}
	
	pointProductDao := s.getPointProductDao()
	return pointProductDao.DeleteProduct(productID)
}