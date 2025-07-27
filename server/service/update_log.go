package services

import (
	"fmt"
	"time"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type UpdateLogService struct {
}

func (u *UpdateLogService) getUpdateLogDao() (*dao.UpdateLogDao, error) {
	db := mysql.GetInstance()
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败，请检查数据库配置")
	}
	return dao.NewUpdateLogDao(db), nil
}

// GetUpdateLogList 获取更新日志列表
func (u *UpdateLogService) GetUpdateLogList(req model.UpdateLogListRequest) (*model.UpdateLogListResponse, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	items, total, err := dao.GetUpdateLogList(req)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	return &model.UpdateLogListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetAdminUpdateLogList 获取更新日志列表（管理员专用）
func (u *UpdateLogService) GetAdminUpdateLogList(req model.UpdateLogListRequest) (*model.UpdateLogListResponse, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	items, total, err := dao.GetAdminUpdateLogList(req)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	return &model.UpdateLogListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetUpdateLogByID 根据ID获取更新日志详情
func (u *UpdateLogService) GetUpdateLogByID(id int64) (*model.UpdateLog, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	return dao.GetUpdateLogByID(id)
}

// CreateUpdateLog 创建更新日志
func (u *UpdateLogService) CreateUpdateLog(req model.UpdateLogRequest) (*model.UpdateLog, error) {
	// 验证数据
	if err := u.validateUpdateLogRequest(req); err != nil {
		return nil, err
	}

	// 解析发布日期
	var publishDate time.Time
	if req.PublishDate != "" {
		var err error
		publishDate, err = time.Parse("2006-01-02", req.PublishDate)
		if err != nil {
			return nil, fmt.Errorf("发布日期格式不正确，请使用YYYY-MM-DD格式")
		}
	} else {
		publishDate = time.Now()
	}

	// 构建模型
	item := &model.UpdateLog{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Version:     req.Version,
		Type:        req.Type,
		Status:      req.Status,
		PublishDate: publishDate,
	}

	// 设置默认值
	if item.Status == "" {
		item.Status = "active"
	}
	if item.Type == "" {
		item.Type = "feature"
	}

	// 创建更新日志
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	if err := dao.CreateUpdateLog(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateUpdateLog 更新更新日志
func (u *UpdateLogService) UpdateUpdateLog(id int64, req model.UpdateLogRequest) (*model.UpdateLog, error) {
	// 验证更新日志是否存在
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	existingItem, err := dao.GetUpdateLogByID(id)
	if err != nil {
		return nil, err
	}

	// 验证数据
	if err := u.validateUpdateLogRequest(req); err != nil {
		return nil, err
	}

	// 解析发布日期
	if req.PublishDate != "" {
		publishDate, err := time.Parse("2006-01-02", req.PublishDate)
		if err != nil {
			return nil, fmt.Errorf("发布日期格式不正确，请使用YYYY-MM-DD格式")
		}
		existingItem.PublishDate = publishDate
	}

	// 更新字段
	existingItem.Title = req.Title
	existingItem.Description = req.Description
	existingItem.Content = req.Content
	existingItem.Version = req.Version

	if req.Status != "" {
		existingItem.Status = req.Status
	}
	if req.Type != "" {
		existingItem.Type = req.Type
	}

	// 更新更新日志
	if err := dao.UpdateUpdateLog(id, existingItem); err != nil {
		return nil, err
	}

	return existingItem, nil
}

// DeleteUpdateLog 删除更新日志
func (u *UpdateLogService) DeleteUpdateLog(id int64) error {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return err
	}
	// 验证更新日志是否存在
	if _, err := dao.GetUpdateLogByID(id); err != nil {
		return err
	}

	return dao.DeleteUpdateLog(id)
}

// UpdateUpdateLogStatus 更新更新日志状态
func (u *UpdateLogService) UpdateUpdateLogStatus(id int64, status string) error {
	// 验证状态值长度
	if len(status) > 20 {
		return fmt.Errorf("状态值长度不能超过20个字符")
	}

	dao, err := u.getUpdateLogDao()
	if err != nil {
		return err
	}
	// 验证更新日志是否存在
	if _, err := dao.GetUpdateLogByID(id); err != nil {
		return err
	}

	return dao.UpdateUpdateLogStatus(id, status)
}

// GetActiveUpdateLogList 获取活跃更新日志列表（用于前端展示）
func (u *UpdateLogService) GetActiveUpdateLogList() ([]model.UpdateLog, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	return dao.GetActiveUpdateLogList()
}

// GetRecentUpdateLogs 获取最近的更新日志（用于首页展示）
func (u *UpdateLogService) GetRecentUpdateLogs(limit int) ([]model.UpdateLogSimple, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	
	items, err := dao.GetRecentUpdateLogs(limit)
	if err != nil {
		return nil, err
	}

	// 转换为简化结构
	var simpleItems []model.UpdateLogSimple
	for _, item := range items {
		simpleItems = append(simpleItems, model.UpdateLogSimple{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Content:     item.Content,
			Version:     item.Version,
			Type:        item.Type,
			Date:        item.PublishDate.Format("2006-01-02"),
		})
	}

	return simpleItems, nil
}

// GetUpdateLogListByType 根据类型获取更新日志列表
func (u *UpdateLogService) GetUpdateLogListByType(logType string) ([]model.UpdateLog, error) {
	dao, err := u.getUpdateLogDao()
	if err != nil {
		return nil, err
	}
	return dao.GetUpdateLogListByType(logType)
}

// validateUpdateLogRequest 验证更新日志请求数据
func (u *UpdateLogService) validateUpdateLogRequest(req model.UpdateLogRequest) error {
	if req.Title == "" {
		return fmt.Errorf("更新日志标题不能为空")
	}
	if len(req.Title) > 255 {
		return fmt.Errorf("更新日志标题长度不能超过255个字符")
	}
	if req.Description == "" {
		return fmt.Errorf("更新日志描述不能为空")
	}
	if len(req.Version) > 50 {
		return fmt.Errorf("版本号长度不能超过50个字符")
	}
	if len(req.Type) > 50 {
		return fmt.Errorf("更新日志类型长度不能超过50个字符")
	}
	if len(req.Status) > 20 {
		return fmt.Errorf("状态值长度不能超过20个字符")
	}
	
	// 验证发布日期格式
	if req.PublishDate != "" {
		if _, err := time.Parse("2006-01-02", req.PublishDate); err != nil {
			return fmt.Errorf("发布日期格式不正确，请使用YYYY-MM-DD格式")
		}
	}
	
	return nil
}