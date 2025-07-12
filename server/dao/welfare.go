package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type WelfareDao struct {
	db *gorm.DB
}

func NewWelfareDao(db *gorm.DB) *WelfareDao {
	return &WelfareDao{db: db}
}

// GetWelfareList 获取福利列表
func (w *WelfareDao) GetWelfareList(req model.WelfareListRequest) ([]model.WelfareItem, int64, error) {
	var items []model.WelfareItem
	var total int64

	query := w.db.Model(&model.WelfareItem{})

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		// 默认只查询活跃状态（用于前端用户）
		query = query.Where("status = ?", "active")
	}

	// 标签筛选
	if req.Tag != "" {
		query = query.Where("tag = ?", req.Tag)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取福利总数失败: %v", err)
	}

	// 分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据
	if err := query.Order("sort_order DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("查询福利列表失败: %v", err)
	}

	return items, total, nil
}

// GetAdminWelfareList 获取福利列表（管理员专用，显示所有状态）
func (w *WelfareDao) GetAdminWelfareList(req model.WelfareListRequest) ([]model.WelfareItem, int64, error) {
	var items []model.WelfareItem
	var total int64

	query := w.db.Model(&model.WelfareItem{})

	// 状态筛选（管理员可以查看所有状态）
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	// 不设置默认状态筛选，允许查询所有状态

	// 标签筛选
	if req.Tag != "" {
		query = query.Where("tag = ?", req.Tag)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取福利总数失败: %v", err)
	}

	// 分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据
	if err := query.Order("sort_order DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("查询福利列表失败: %v", err)
	}

	return items, total, nil
}

// GetWelfareByID 根据ID获取福利详情
func (w *WelfareDao) GetWelfareByID(id int64) (*model.WelfareItem, error) {
	var item model.WelfareItem
	if err := w.db.Where("id = ?", id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("福利不存在")
		}
		return nil, fmt.Errorf("查询福利失败: %v", err)
	}
	return &item, nil
}

// CreateWelfare 创建福利
func (w *WelfareDao) CreateWelfare(item *model.WelfareItem) error {
	if err := w.db.Create(item).Error; err != nil {
		return fmt.Errorf("创建福利失败: %v", err)
	}
	return nil
}

// UpdateWelfare 更新福利
func (w *WelfareDao) UpdateWelfare(id int64, item *model.WelfareItem) error {
	result := w.db.Where("id = ?", id).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("更新福利失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("福利不存在")
	}
	return nil
}

// DeleteWelfare 删除福利
func (w *WelfareDao) DeleteWelfare(id int64) error {
	result := w.db.Where("id = ?", id).Delete(&model.WelfareItem{})
	if result.Error != nil {
		return fmt.Errorf("删除福利失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("福利不存在")
	}
	return nil
}

// UpdateWelfareStatus 更新福利状态
func (w *WelfareDao) UpdateWelfareStatus(id int64, status string) error {
	result := w.db.Model(&model.WelfareItem{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("更新福利状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("福利不存在")
	}
	return nil
}

// GetActiveWelfareList 获取活跃福利列表（用于前端展示）
func (w *WelfareDao) GetActiveWelfareList() ([]model.WelfareItem, error) {
	var items []model.WelfareItem
	if err := w.db.Where("status = ?", "active").
		Order("sort_order DESC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("查询活跃福利列表失败: %v", err)
	}
	return items, nil
}

// GetWelfareListByTag 根据标签获取福利列表
func (w *WelfareDao) GetWelfareListByTag(tag string) ([]model.WelfareItem, error) {
	var items []model.WelfareItem
	query := w.db.Where("status = ?", "active")
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if err := query.Order("sort_order DESC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("根据标签查询福利列表失败: %v", err)
	}
	return items, nil
}