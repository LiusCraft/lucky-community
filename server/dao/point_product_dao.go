package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type PointProductDao struct{}

// GetActiveProducts 获取所有启用的积分商品
func (d *PointProductDao) GetActiveProducts() ([]model.PointProduct, error) {
	var products []model.PointProduct
	err := model.PointProductModel().
		Where("status = ?", model.ProductStatusActive).
		Order("sort_order DESC, created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, fmt.Errorf("查询积分商品失败: %v", err)
	}

	return products, nil
}

// GetProductByID 根据ID获取积分商品
func (d *PointProductDao) GetProductByID(id int64) (*model.PointProduct, error) {
	var product model.PointProduct
	err := model.PointProductModel().Where("id = ?", id).First(&product).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("积分商品不存在")
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询积分商品失败: %v", err)
	}
	
	return &product, nil
}

// GetProductsWithPagination 分页获取积分商品
func (d *PointProductDao) GetProductsWithPagination(page, pageSize int, status, rewardType string) ([]model.PointProduct, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := model.PointProductModel()
	
	// 根据状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 根据奖励类型筛选
	if rewardType != "" {
		query = query.Where("reward_type = ?", rewardType)
	}

	var total int64
	var products []model.PointProduct

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询积分商品总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort_order DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	if err != nil {
		return nil, 0, fmt.Errorf("分页查询积分商品失败: %v", err)
	}

	return products, total, nil
}

// CreateProduct 创建积分商品
func (d *PointProductDao) CreateProduct(product *model.PointProduct) error {
	err := model.PointProductModel().Create(product).Error
	if err != nil {
		return fmt.Errorf("创建积分商品失败: %v", err)
	}
	return nil
}

// UpdateProduct 更新积分商品
func (d *PointProductDao) UpdateProduct(product *model.PointProduct) error {
	if product.ID <= 0 {
		return fmt.Errorf("商品ID无效")
	}
	
	err := model.PointProductModel().Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return fmt.Errorf("更新积分商品失败: %v", err)
	}
	return nil
}

// UpdateProductByID 根据ID更新积分商品指定字段
func (d *PointProductDao) UpdateProductByID(productID int64, updates map[string]interface{}) error {
	if productID <= 0 {
		return fmt.Errorf("商品ID无效")
	}
	
	if len(updates) == 0 {
		return nil
	}
	
	err := model.PointProductModel().Where("id = ?", productID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新积分商品失败: %v", err)
	}
	return nil
}

// UpdateProductStock 更新商品库存
func (d *PointProductDao) UpdateProductStock(productID int64, quantity int) error {
	// 检查商品是否存在且有足够库存
	var product model.PointProduct
	if err := model.PointProductModel().Where("id = ?", productID).First(&product).Error; err != nil {
		return fmt.Errorf("商品不存在: %v", err)
	}

	// 无限库存商品不需要更新
	if product.HasUnlimitedStock() {
		return nil
	}

	// 检查库存是否充足
	if product.Stock < quantity {
		return fmt.Errorf("库存不足：需要%d，当前库存%d", quantity, product.Stock)
	}

	// 扣减库存
	err := model.PointProductModel().
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error

	if err != nil {
		return fmt.Errorf("更新商品库存失败: %v", err)
	}

	return nil
}

// UpdateProductStockWithTransaction 在事务中更新商品库存
func (d *PointProductDao) UpdateProductStockWithTransaction(tx *gorm.DB, productID int64, quantity int) error {
	// 检查商品是否存在且有足够库存
	var product model.PointProduct
	if err := tx.Model(&model.PointProduct{}).Where("id = ?", productID).First(&product).Error; err != nil {
		return fmt.Errorf("商品不存在: %v", err)
	}

	// 无限库存商品不需要更新
	if product.HasUnlimitedStock() {
		return nil
	}

	// 检查库存是否充足
	if product.Stock < quantity {
		return fmt.Errorf("库存不足：需要%d，当前库存%d", quantity, product.Stock)
	}

	// 扣减库存
	err := tx.Model(&model.PointProduct{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error

	if err != nil {
		return fmt.Errorf("更新商品库存失败: %v", err)
	}

	return nil
}

// RestoreProductStock 恢复商品库存（用于取消兑换）
func (d *PointProductDao) RestoreProductStock(productID int64, quantity int) error {
	// 检查商品是否存在
	var product model.PointProduct
	if err := model.PointProductModel().Where("id = ?", productID).First(&product).Error; err != nil {
		return fmt.Errorf("商品不存在: %v", err)
	}

	// 无限库存商品不需要恢复
	if product.HasUnlimitedStock() {
		return nil
	}

	// 恢复库存
	err := model.PointProductModel().
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error

	if err != nil {
		return fmt.Errorf("恢复商品库存失败: %v", err)
	}

	return nil
}

// UpdateProductStatus 更新商品状态
func (d *PointProductDao) UpdateProductStatus(id int64, status string) error {
	err := model.PointProductModel().Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("更新商品状态失败: %v", err)
	}
	return nil
}

// DeleteProduct 删除积分商品
func (d *PointProductDao) DeleteProduct(id int64) error {
	err := model.PointProductModel().Where("id = ?", id).Delete(&model.PointProduct{}).Error
	if err != nil {
		return fmt.Errorf("删除积分商品失败: %v", err)
	}
	return nil
}

// GetProductsByRewardType 根据奖励类型获取商品
func (d *PointProductDao) GetProductsByRewardType(rewardType string) ([]model.PointProduct, error) {
	var products []model.PointProduct
	err := model.PointProductModel().
		Where("reward_type = ? AND status = ?", rewardType, model.ProductStatusActive).
		Order("sort_order DESC, price_points ASC").
		Find(&products).Error

	if err != nil {
		return nil, fmt.Errorf("根据奖励类型查询商品失败: %v", err)
	}

	return products, nil
}

// GetProductStatistics 获取商品统计信息
func (d *PointProductDao) GetProductStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalProducts    int64 `json:"total_products"`
		ActiveProducts   int64 `json:"active_products"`
		InactiveProducts int64 `json:"inactive_products"`
	}

	// 获取商品总数
	if err := model.PointProductModel().Count(&stats.TotalProducts).Error; err != nil {
		return nil, fmt.Errorf("查询商品总数失败: %v", err)
	}

	// 获取启用商品数
	if err := model.PointProductModel().Where("status = ?", model.ProductStatusActive).Count(&stats.ActiveProducts).Error; err != nil {
		return nil, fmt.Errorf("查询启用商品数失败: %v", err)
	}

	// 获取禁用商品数
	if err := model.PointProductModel().Where("status = ?", model.ProductStatusInactive).Count(&stats.InactiveProducts).Error; err != nil {
		return nil, fmt.Errorf("查询禁用商品数失败: %v", err)
	}

	result := map[string]interface{}{
		"total_products":    stats.TotalProducts,
		"active_products":   stats.ActiveProducts,
		"inactive_products": stats.InactiveProducts,
	}

	return result, nil
}

// GetRewardTypeStatistics 获取按奖励类型的商品统计
func (d *PointProductDao) GetRewardTypeStatistics() ([]map[string]interface{}, error) {
	var results []struct {
		RewardType    string `json:"reward_type"`
		ProductCount  int64  `json:"product_count"`
	}

	err := model.PointProductModel().
		Where("status = ?", model.ProductStatusActive).
		Select("reward_type, COUNT(*) as product_count").
		Group("reward_type").
		Order("product_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询奖励类型统计失败: %v", err)
	}

	// 转换为 map 切片
	var statistics []map[string]interface{}
	for _, result := range results {
		statistics = append(statistics, map[string]interface{}{
			"reward_type":    result.RewardType,
			"product_count":  result.ProductCount,
		})
	}

	return statistics, nil
}

// BatchUpdateProductStatus 批量更新商品状态
func (d *PointProductDao) BatchUpdateProductStatus(ids []int64, status string) error {
	err := model.PointProductModel().Where("id IN ?", ids).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("批量更新商品状态失败: %v", err)
	}
	return nil
}