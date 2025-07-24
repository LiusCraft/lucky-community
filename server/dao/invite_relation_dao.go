package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type InviteRelationDao struct{}

// CreateInviteRelation 创建邀请关系
func (d *InviteRelationDao) CreateInviteRelation(relation *model.InviteRelation) error {
	err := model.InviteRelationModel().Create(relation).Error
	if err != nil {
		return fmt.Errorf("创建邀请关系失败: %v", err)
	}
	return nil
}

// GetInviteRelationByInvitee 根据被邀请人ID查询邀请关系
func (d *InviteRelationDao) GetInviteRelationByInvitee(inviteeID int) (*model.InviteRelation, error) {
	var relation model.InviteRelation
	err := model.InviteRelationModel().Where("invitee_id = ?", inviteeID).First(&relation).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, nil // 返回nil表示没有找到邀请关系
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询邀请关系失败: %v", err)
	}
	
	return &relation, nil
}

// GetInviteRelationsByInviter 根据邀请人ID查询邀请关系列表
func (d *InviteRelationDao) GetInviteRelationsByInviter(inviterID int, page, pageSize int) ([]model.InviteRelation, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	var relations []model.InviteRelation

	query := model.InviteRelationModel().Where("inviter_id = ?", inviterID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询邀请关系总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&relations).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询邀请关系列表失败: %v", err)
	}

	return relations, total, nil
}

// GetInviteRelationByCode 根据邀请码查询邀请关系
func (d *InviteRelationDao) GetInviteRelationByCode(inviteCode string) ([]model.InviteRelation, error) {
	var relations []model.InviteRelation
	err := model.InviteRelationModel().
		Where("invite_code = ?", inviteCode).
		Order("created_at DESC").
		Find(&relations).Error

	if err != nil {
		return nil, fmt.Errorf("根据邀请码查询邀请关系失败: %v", err)
	}

	return relations, nil
}

// GetInviterStatistics 获取邀请人的统计信息
func (d *InviteRelationDao) GetInviterStatistics(inviterID int) (map[string]interface{}, error) {
	var stats struct {
		TotalInvites int64 `json:"total_invites"`
	}

	// 获取邀请总数
	if err := model.InviteRelationModel().
		Where("inviter_id = ?", inviterID).
		Count(&stats.TotalInvites).Error; err != nil {
		return nil, fmt.Errorf("查询邀请统计失败: %v", err)
	}

	result := map[string]interface{}{
		"total_invites": stats.TotalInvites,
	}

	return result, nil
}

// GetInviteStatistics 获取整体邀请统计
func (d *InviteRelationDao) GetInviteStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalInviteRelations int64 `json:"total_invite_relations"`
		TotalInviters        int64 `json:"total_inviters"`
	}

	// 获取邀请关系总数
	if err := model.InviteRelationModel().Count(&stats.TotalInviteRelations).Error; err != nil {
		return nil, fmt.Errorf("查询邀请关系总数失败: %v", err)
	}

	// 获取邀请人总数（去重）
	if err := model.InviteRelationModel().
		Select("COUNT(DISTINCT inviter_id)").
		Scan(&stats.TotalInviters).Error; err != nil {
		return nil, fmt.Errorf("查询邀请人总数失败: %v", err)
	}

	result := map[string]interface{}{
		"total_invite_relations": stats.TotalInviteRelations,
		"total_inviters":         stats.TotalInviters,
		"average_invites_per_inviter": func() float64 {
			if stats.TotalInviters > 0 {
				return float64(stats.TotalInviteRelations) / float64(stats.TotalInviters)
			}
			return 0
		}(),
	}

	return result, nil
}

// GetTopInviters 获取邀请排行榜
func (d *InviteRelationDao) GetTopInviters(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var results []struct {
		InviterID    int   `json:"inviter_id"`
		InviteCount  int64 `json:"invite_count"`
	}

	err := model.InviteRelationModel().
		Select("inviter_id, COUNT(*) as invite_count").
		Group("inviter_id").
		Order("invite_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("查询邀请排行榜失败: %v", err)
	}

	// 转换为 map 切片
	var rankings []map[string]interface{}
	for i, result := range results {
		rankings = append(rankings, map[string]interface{}{
			"rank":         i + 1,
			"inviter_id":   result.InviterID,
			"invite_count": result.InviteCount,
		})
	}

	return rankings, nil
}

// IsInviteeExists 检查被邀请人是否已存在邀请关系
func (d *InviteRelationDao) IsInviteeExists(inviteeID int) (bool, error) {
	var count int64
	err := model.InviteRelationModel().Where("invitee_id = ?", inviteeID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查被邀请人是否存在失败: %v", err)
	}
	return count > 0, nil
}

// GetInviteRelationsWithPagination 分页获取所有邀请关系（管理员用）
func (d *InviteRelationDao) GetInviteRelationsWithPagination(page, pageSize int) ([]model.InviteRelation, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	var relations []model.InviteRelation

	// 获取总数
	if err := model.InviteRelationModel().Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询邀请关系总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := model.InviteRelationModel().
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&relations).Error

	if err != nil {
		return nil, 0, fmt.Errorf("分页查询邀请关系失败: %v", err)
	}

	return relations, total, nil
}

// DeleteInviteRelation 删除邀请关系（管理员操作）
func (d *InviteRelationDao) DeleteInviteRelation(id int64) error {
	err := model.InviteRelationModel().Where("id = ?", id).Delete(&model.InviteRelation{}).Error
	if err != nil {
		return fmt.Errorf("删除邀请关系失败: %v", err)
	}
	return nil
}