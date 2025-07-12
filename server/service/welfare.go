package services

import (
	"fmt"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type WelfareService struct {
}

func (w *WelfareService) getWelfareDao() (*dao.WelfareDao, error) {
	db := mysql.GetInstance()
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败，请检查数据库配置")
	}
	return dao.NewWelfareDao(db), nil
}

// GetWelfareList 获取福利列表
func (w *WelfareService) GetWelfareList(req model.WelfareListRequest) (*model.WelfareListResponse, error) {
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	items, total, err := dao.GetWelfareList(req)
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

	return &model.WelfareListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetAdminWelfareList 获取福利列表（管理员专用）
func (w *WelfareService) GetAdminWelfareList(req model.WelfareListRequest) (*model.WelfareListResponse, error) {
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	items, total, err := dao.GetAdminWelfareList(req)
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

	return &model.WelfareListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// GetWelfareByID 根据ID获取福利详情
func (w *WelfareService) GetWelfareByID(id int64) (*model.WelfareItem, error) {
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	return dao.GetWelfareByID(id)
}

// CreateWelfare 创建福利
func (w *WelfareService) CreateWelfare(req model.WelfareItemRequest) (*model.WelfareItem, error) {
	// 验证数据
	if err := w.validateWelfareRequest(req); err != nil {
		return nil, err
	}

	// 构建模型
	item := &model.WelfareItem{
		Title:         req.Title,
		Description:   req.Description,
		DetailContent: req.DetailContent,
		Tag:           req.Tag,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		DiscountText:  req.DiscountText,
		ActionText:    req.ActionText,
		Status:        req.Status,
		SortOrder:     req.SortOrder,
	}

	// 设置默认值
	if item.Status == "" {
		item.Status = "active"
	}
	if item.ActionText == "" {
		item.ActionText = "立即查看"
	}

	// 创建福利
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	if err := dao.CreateWelfare(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateWelfare 更新福利
func (w *WelfareService) UpdateWelfare(id int64, req model.WelfareItemRequest) (*model.WelfareItem, error) {
	// 验证福利是否存在
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	existingItem, err := dao.GetWelfareByID(id)
	if err != nil {
		return nil, err
	}

	// 验证数据
	if err := w.validateWelfareRequest(req); err != nil {
		return nil, err
	}

	// 更新字段
	existingItem.Title = req.Title
	existingItem.Description = req.Description
	existingItem.DetailContent = req.DetailContent
	existingItem.Tag = req.Tag
	existingItem.Price = req.Price
	existingItem.OriginalPrice = req.OriginalPrice
	existingItem.DiscountText = req.DiscountText
	existingItem.ActionText = req.ActionText
	existingItem.SortOrder = req.SortOrder

	if req.Status != "" {
		existingItem.Status = req.Status
	}
	if existingItem.ActionText == "" {
		existingItem.ActionText = "立即查看"
	}

	// 更新福利
	if err := dao.UpdateWelfare(id, existingItem); err != nil {
		return nil, err
	}

	return existingItem, nil
}

// DeleteWelfare 删除福利
func (w *WelfareService) DeleteWelfare(id int64) error {
	dao, err := w.getWelfareDao()
	if err != nil {
		return err
	}
	// 验证福利是否存在
	if _, err := dao.GetWelfareByID(id); err != nil {
		return err
	}

	return dao.DeleteWelfare(id)
}

// UpdateWelfareStatus 更新福利状态
func (w *WelfareService) UpdateWelfareStatus(id int64, status string) error {
	// 验证状态值
	if status != "active" && status != "inactive" {
		return fmt.Errorf("无效的状态值: %s", status)
	}

	dao, err := w.getWelfareDao()
	if err != nil {
		return err
	}
	// 验证福利是否存在
	if _, err := dao.GetWelfareByID(id); err != nil {
		return err
	}

	return dao.UpdateWelfareStatus(id, status)
}

// GetActiveWelfareList 获取活跃福利列表（用于前端展示）
func (w *WelfareService) GetActiveWelfareList() ([]model.WelfareItem, error) {
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	return dao.GetActiveWelfareList()
}

// GetWelfareListByTag 根据标签获取福利列表
func (w *WelfareService) GetWelfareListByTag(tag string) ([]model.WelfareItem, error) {
	dao, err := w.getWelfareDao()
	if err != nil {
		return nil, err
	}
	return dao.GetWelfareListByTag(tag)
}

// validateWelfareRequest 验证福利请求数据
func (w *WelfareService) validateWelfareRequest(req model.WelfareItemRequest) error {
	if req.Title == "" {
		return fmt.Errorf("福利标题不能为空")
	}
	if len(req.Title) > 255 {
		return fmt.Errorf("福利标题长度不能超过255个字符")
	}
	if req.Tag == "" {
		return fmt.Errorf("福利标签不能为空")
	}
	if len(req.Tag) > 50 {
		return fmt.Errorf("福利标签长度不能超过50个字符")
	}
	if req.Price < 0 {
		return fmt.Errorf("价格不能为负数")
	}
	if req.OriginalPrice < 0 {
		return fmt.Errorf("原价不能为负数")
	}
	if req.OriginalPrice > 0 && req.Price > req.OriginalPrice {
		return fmt.Errorf("现价不能高于原价")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return fmt.Errorf("无效的状态值: %s", req.Status)
	}
	return nil
}