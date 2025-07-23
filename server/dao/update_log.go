package dao

import (
	"fmt"
	"time"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type UpdateLogDao struct {
	db *gorm.DB
}

func NewUpdateLogDao(db *gorm.DB) *UpdateLogDao {
	return &UpdateLogDao{db: db}
}

// GetUpdateLogList 获取更新日志列表
func (u *UpdateLogDao) GetUpdateLogList(req model.UpdateLogListRequest) ([]model.UpdateLog, int64, error) {
	var items []model.UpdateLog
	var total int64

	query := u.db.Model(&model.UpdateLog{})

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		// 默认只查询活跃状态（用于前端用户）
		query = query.Where("status = ?", "active")
	}

	// 类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取更新日志总数失败: %v", err)
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
	if err := query.Order("publish_date DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("查询更新日志列表失败: %v", err)
	}

	return items, total, nil
}

// GetAdminUpdateLogList 获取更新日志列表（管理员专用，显示所有状态）
func (u *UpdateLogDao) GetAdminUpdateLogList(req model.UpdateLogListRequest) ([]model.UpdateLog, int64, error) {
	var items []model.UpdateLog
	var total int64

	query := u.db.Model(&model.UpdateLog{})

	// 状态筛选（管理员可以查看所有状态）
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	// 不设置默认状态筛选，允许查询所有状态

	// 类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取更新日志总数失败: %v", err)
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
	if err := query.Order("publish_date DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("查询更新日志列表失败: %v", err)
	}

	return items, total, nil
}

// GetUpdateLogByID 根据ID获取更新日志详情
func (u *UpdateLogDao) GetUpdateLogByID(id int64) (*model.UpdateLog, error) {
	var item model.UpdateLog
	if err := u.db.Where("id = ?", id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("更新日志不存在")
		}
		return nil, fmt.Errorf("查询更新日志失败: %v", err)
	}
	return &item, nil
}

// CreateUpdateLog 创建更新日志
func (u *UpdateLogDao) CreateUpdateLog(item *model.UpdateLog) error {
	if err := u.db.Create(item).Error; err != nil {
		return fmt.Errorf("创建更新日志失败: %v", err)
	}
	return nil
}

// UpdateUpdateLog 更新更新日志
func (u *UpdateLogDao) UpdateUpdateLog(id int64, item *model.UpdateLog) error {
	result := u.db.Where("id = ?", id).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("更新更新日志失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("更新日志不存在")
	}
	return nil
}

// DeleteUpdateLog 删除更新日志
func (u *UpdateLogDao) DeleteUpdateLog(id int64) error {
	result := u.db.Where("id = ?", id).Delete(&model.UpdateLog{})
	if result.Error != nil {
		return fmt.Errorf("删除更新日志失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("更新日志不存在")
	}
	return nil
}

// UpdateUpdateLogStatus 更新更新日志状态
func (u *UpdateLogDao) UpdateUpdateLogStatus(id int64, status string) error {
	result := u.db.Model(&model.UpdateLog{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("更新更新日志状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("更新日志不存在")
	}
	return nil
}

// GetActiveUpdateLogList 获取活跃更新日志列表（用于前端展示）
func (u *UpdateLogDao) GetActiveUpdateLogList() ([]model.UpdateLog, error) {
	var items []model.UpdateLog
	if err := u.db.Where("status = ?", "active").
		Where("publish_date <= ?", time.Now()).
		Order("publish_date DESC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("查询活跃更新日志列表失败: %v", err)
	}
	return items, nil
}

// GetRecentUpdateLogs 获取最近的更新日志（用于首页展示）
func (u *UpdateLogDao) GetRecentUpdateLogs(limit int) ([]model.UpdateLog, error) {
	if limit <= 0 {
		limit = 5
	}

	var items []model.UpdateLog
	if err := u.db.Where("status = ?", "active").
		Where("publish_date <= ?", time.Now()).
		Order("publish_date DESC, created_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("查询最近更新日志失败: %v", err)
	}
	return items, nil
}

// GetUpdateLogListByType 根据类型获取更新日志列表
func (u *UpdateLogDao) GetUpdateLogListByType(logType string) ([]model.UpdateLog, error) {
	var items []model.UpdateLog
	query := u.db.Where("status = ?", "active").
		Where("publish_date <= ?", time.Now())
	
	if logType != "" {
		query = query.Where("type = ?", logType)
	}
	
	if err := query.Order("publish_date DESC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("根据类型查询更新日志列表失败: %v", err)
	}
	return items, nil
}