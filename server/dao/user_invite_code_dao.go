package dao

import (
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type UserInviteCodeDao struct{}

// GetUserInviteCode 获取用户的邀请码，如果不存在则创建
func (d *UserInviteCodeDao) GetUserInviteCode(userID int) (*model.UserInviteCode, error) {
	var inviteCode model.UserInviteCode
	err := model.UserInviteCodeModel().Where("user_id = ?", userID).First(&inviteCode).Error
	
	if err == gorm.ErrRecordNotFound {
		// 用户邀请码不存在，创建新的邀请码
		newCode := model.GenerateInviteCode(userID)
		inviteCode = model.UserInviteCode{
			UserID:     userID,
			InviteCode: newCode,
		}
		
		if createErr := model.UserInviteCodeModel().Create(&inviteCode).Error; createErr != nil {
			return nil, fmt.Errorf("创建用户邀请码失败: %v", createErr)
		}
		
		return &inviteCode, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询用户邀请码失败: %v", err)
	}
	
	return &inviteCode, nil
}

// GetInviteCodeByCode 根据邀请码查询邀请码信息
func (d *UserInviteCodeDao) GetInviteCodeByCode(code string) (*model.UserInviteCode, error) {
	var inviteCode model.UserInviteCode
	err := model.UserInviteCodeModel().Where("invite_code = ?", code).First(&inviteCode).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("邀请码不存在")
	}
	
	if err != nil {
		return nil, fmt.Errorf("查询邀请码失败: %v", err)
	}
	
	return &inviteCode, nil
}





// GetInviteCodeStatistics 获取邀请码系统统计
func (d *UserInviteCodeDao) GetInviteCodeStatistics() (map[string]interface{}, error) {
	var totalCodes int64
	
	// 获取邀请码总数
	if err := model.UserInviteCodeModel().Count(&totalCodes).Error; err != nil {
		return nil, fmt.Errorf("查询邀请码总数失败: %v", err)
	}
	
	result := map[string]interface{}{
		"total_codes": totalCodes,
	}
	
	return result, nil
}

// GetInviteCodesWithPagination 分页获取邀请码列表（管理员用）
func (d *UserInviteCodeDao) GetInviteCodesWithPagination(page, pageSize int, userID int) ([]model.UserInviteCode, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	
	query := model.UserInviteCodeModel()
	
	// 根据userID筛选
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	
	var total int64
	var inviteCodes []model.UserInviteCode
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询邀请码总数失败: %v", err)
	}
	
	// 分页查询，预加载用户信息
	offset := (page - 1) * pageSize
	err := query.
		Preload("User").  // 预加载用户信息
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&inviteCodes).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("分页查询邀请码失败: %v", err)
	}
	
	return inviteCodes, total, nil
}

// IsInviteCodeExists 检查邀请码是否存在
func (d *UserInviteCodeDao) IsInviteCodeExists(code string) (bool, error) {
	var count int64
	err := model.UserInviteCodeModel().Where("invite_code = ?", code).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查邀请码是否存在失败: %v", err)
	}
	return count > 0, nil
}

// DeleteInviteCode 删除邀请码（管理员操作）
func (d *UserInviteCodeDao) DeleteInviteCode(userID int) error {
	err := model.UserInviteCodeModel().Where("user_id = ?", userID).Delete(&model.UserInviteCode{}).Error
	if err != nil {
		return fmt.Errorf("删除邀请码失败: %v", err)
	}
	return nil
}

// BatchCreateInviteCodes 批量创建邀请码
func (d *UserInviteCodeDao) BatchCreateInviteCodes(userIDs []int) error {
	var inviteCodes []model.UserInviteCode
	
	for _, userID := range userIDs {
		inviteCodes = append(inviteCodes, model.UserInviteCode{
			UserID:     userID,
			InviteCode: model.GenerateInviteCode(userID),
		})
	}
	
	if err := model.UserInviteCodeModel().CreateInBatches(inviteCodes, 100).Error; err != nil {
		return fmt.Errorf("批量创建邀请码失败: %v", err)
	}
	
	return nil
}