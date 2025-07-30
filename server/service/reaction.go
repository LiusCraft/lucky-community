package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type ReactionService struct {
	ctx *gin.Context
}

// NewReactionService 创建通用表情回复服务实例
func NewReactionService(ctx *gin.Context) *ReactionService {
	return &ReactionService{ctx: ctx}
}

// ToggleReaction 切换表情回复状态
func (s *ReactionService) ToggleReaction(businessType, businessId, userId int, reactionType string) (bool, error) {

	reactionDao := dao.GetReactionDao()

	// 检查用户是否已经添加过此表情
	exists, err := reactionDao.CheckUserReaction(businessType, businessId, userId, reactionType)
	if err != nil {
		return false, err
	}

	if exists {
		// 已存在，移除表情回复
		err = reactionDao.RemoveReaction(businessType, businessId, userId, reactionType)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		// 不存在，添加表情回复
		reaction := &model.Reaction{
			BusinessType: businessType,
			BusinessId:   businessId,
			UserId:       userId,
			ReactionType: reactionType,
		}
		err = reactionDao.AddReaction(reaction)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

// GetReactionSummary 获取单个业务的表情统计
func (s *ReactionService) GetReactionSummary(businessType, businessId, currentUserId int) ([]model.ReactionSummary, error) {
	log.Infof("获取表情统计，业务类型: %d, 业务ID: %d, 当前用户ID: %d", businessType, businessId, currentUserId)

	reactionDao := dao.GetReactionDao()
	summaries, err := reactionDao.GetReactionSummaryByBusiness(businessType, businessId, currentUserId)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

// GetReactionSummaryBatch 批量获取多个业务的表情统计
func (s *ReactionService) GetReactionSummaryBatch(businessType int, businessIds []int, currentUserId int) (map[int][]model.ReactionSummary, error) {

	if len(businessIds) == 0 {
		return make(map[int][]model.ReactionSummary), nil
	}

	reactionDao := dao.GetReactionDao()
	summaryMap, err := reactionDao.GetReactionSummaryByBusinessBatch(businessType, businessIds, currentUserId)
	if err != nil {
		return nil, err
	}

	return summaryMap, nil
}

// ValidateBusinessType 验证业务类型是否有效
func (s *ReactionService) ValidateBusinessType(businessType int) bool {
	return businessType >= model.BusinessTypeArticle && businessType <= model.BusinessTypeAINews
}

// ValidateReactionType 验证表情类型是否有效
func (s *ReactionService) ValidateReactionType(reactionType string) bool {
	// 先进行基本的非空验证
	if reactionType == "" {
		return false
	}

	// 从数据库查询表情类型是否存在且启用
	var count int64
	err := model.ReactionDB().Model(&model.ExpressionType{}).
		Where("code = ? AND is_active = ?", reactionType, true).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}

// GetAllExpressionTypes 获取所有表情类型
func (s *ReactionService) GetAllExpressionTypes() ([]model.ExpressionType, error) {

	var types []model.ExpressionType
	err := model.ReactionDB().
		Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&types).Error

	if err != nil {
		return nil, err
	}

	return types, nil
}

// PageExpressionTypes 分页获取表情类型（管理后台用）
func (s *ReactionService) PageExpressionTypes(page, limit int) ([]model.ExpressionType, int64, error) {

	var types []model.ExpressionType
	var total int64

	db := model.ReactionDB().Model(&model.ExpressionType{})

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = db.Order("sort_order ASC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&types).Error

	if err != nil {
		return nil, 0, err
	}

	return types, total, nil
}

// CreateExpressionType 创建表情类型
func (s *ReactionService) CreateExpressionType(expression *model.ExpressionType) (*model.ExpressionType, error) {

	// 检查code是否已存在
	var count int64
	err := model.ReactionDB().Model(&model.ExpressionType{}).
		Where("code = ?", expression.Code).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
	}

	err = model.ReactionDB().Create(expression).Error
	if err != nil {
		return nil, err
	}

	return expression, nil
}

// UpdateExpressionType 更新表情类型
func (s *ReactionService) UpdateExpressionType(expression *model.ExpressionType) error {

	// 检查表情是否存在
	var existing model.ExpressionType
	err := model.ReactionDB().First(&existing, expression.ID).Error
	if err != nil {
		return err
	}

	// 如果更改了code，检查新code是否已存在
	if existing.Code != expression.Code {
		var count int64
		err := model.ReactionDB().Model(&model.ExpressionType{}).
			Where("code = ? AND id != ?", expression.Code, expression.ID).
			Count(&count).Error

		if err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("表情代码已存在")
		}
	}

	err = model.ReactionDB().Save(expression).Error
	if err != nil {
		log.Errorf("更新表情类型失败: %v", err)
		return err
	}

	log.Infof("更新表情类型成功: %d", expression.ID)
	return nil
}

// DeleteExpressionType 删除表情类型
func (s *ReactionService) DeleteExpressionType(id int) error {

	err := model.ReactionDB().Delete(&model.ExpressionType{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

// ToggleExpressionStatus 切换表情启用状态
func (s *ReactionService) ToggleExpressionStatus(id int) (bool, error) {

	var expression model.ExpressionType
	err := model.ReactionDB().First(&expression, id).Error
	if err != nil {
		return false, err
	}

	expression.IsActive = !expression.IsActive
	err = model.ReactionDB().Save(&expression).Error
	if err != nil {
		return false, err
	}

	return expression.IsActive, nil
}

// CheckExpressionInUse 检查表情是否被使用
func (s *ReactionService) CheckExpressionInUse(expressionId int) (bool, error) {

	// 先获取表情代码
	var expression model.ExpressionType
	err := model.ReactionDB().First(&expression, expressionId).Error
	if err != nil {
		return false, err
	}

	// 检查是否有使用该表情代码的回复
	var count int64
	err = model.ReactionDB().Model(&model.Reaction{}).
		Where("reaction_type = ?", expression.Code).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	inUse := count > 0
	return inUse, nil
}
