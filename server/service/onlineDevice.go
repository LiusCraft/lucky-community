package services

import (
	"crypto/md5"
	"fmt"
	"time"

	xt "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

type OnlineDeviceService struct{}

// DeviceInfo 设备信息结构
type DeviceInfo struct {
	SessionID     string    `json:"sessionId"`
	DeviceInfo    string    `json:"deviceInfo"`
	IPAddress     string    `json:"ipAddress"`
	LoginTime     time.Time `json:"loginTime"`
	LastHeartbeat time.Time `json:"lastHeartbeat"`
}

// RegisterDevice 注册设备会话
func (s *OnlineDeviceService) RegisterDevice(userID int, sessionID, deviceInfo, ipAddress string) error {
	// 检查用户的最大设备数限制
	var user model.Users
	if err := model.User().Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 管理员无限制
	var userService UserService
	isAdmin, _ := userService.IsAdmin(userID)
	if !isAdmin {
		// 获取当前在线设备数量
		currentDevices := s.GetUserOnlineDevices(userID)
		if len(currentDevices) >= user.MaxConcurrentDevices {
			// 踢出最早登录的设备
			s.kickOldestDevice(userID, user.MaxConcurrentDevices-1)
		}
	}

	// 创建或更新会话记录
	session := model.UserOnlineSession{
		UserID:        userID,
		SessionID:     sessionID,
		DeviceInfo:    deviceInfo,
		IPAddress:     ipAddress,
		LoginTime:     xt.Now(),
		LastHeartbeat: xt.Now(),
		CreatedAt:     xt.Now(),
		UpdatedAt:     xt.Now(),
	}

	// 使用 UPSERT 操作
	return model.UserOnlineSessionModel().
		Where("session_id = ?", sessionID).
		Assign(session).
		FirstOrCreate(&session).Error
}

// UpdateHeartbeat 更新心跳时间
func (s *OnlineDeviceService) UpdateHeartbeat(sessionID, ipAddress string) error {
	return model.UserOnlineSessionModel().
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_heartbeat": xt.Now(),
			"ip_address":     ipAddress,
			"updated_at":     xt.Now(),
		}).Error
}

// GetUserOnlineDevices 获取用户在线设备列表
func (s *OnlineDeviceService) GetUserOnlineDevices(userID int) []DeviceInfo {
	var sessions []model.UserOnlineSession

	// 清理过期会话（超过5分钟没有心跳的会话）
	s.cleanExpiredSessions()

	model.UserOnlineSessionModel().
		Where("user_id = ?", userID).
		Order("login_time ASC").
		Find(&sessions)

	devices := make([]DeviceInfo, len(sessions))
	for i, session := range sessions {
		devices[i] = DeviceInfo{
			SessionID:     session.SessionID,
			DeviceInfo:    session.DeviceInfo,
			IPAddress:     session.IPAddress,
			LoginTime:     time.Time(session.LoginTime),
			LastHeartbeat: time.Time(session.LastHeartbeat),
		}
	}

	return devices
}

// KickDevice 踢出指定设备
func (s *OnlineDeviceService) KickDevice(userID int, sessionID string) error {
	// 验证会话属于该用户
	var session model.UserOnlineSession
	if err := model.UserOnlineSessionModel().
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		First(&session).Error; err != nil {
		return fmt.Errorf("会话不存在或无权限")
	}

	// 删除会话记录
	if err := model.UserOnlineSessionModel().
		Where("session_id = ?", sessionID).
		Delete(&model.UserOnlineSession{}).Error; err != nil {
		return err
	}

	// 将token加入黑名单
	var blackService BlacklistService
	blackService.AddBlackByToken(sessionID)

	return nil
}

// CheckDeviceLimit 检查设备数量限制
func (s *OnlineDeviceService) CheckDeviceLimit(userID int) (bool, int, int) {
	var user model.Users
	if err := model.User().Where("id = ?", userID).First(&user).Error; err != nil {
		return false, 0, 0
	}

	// 管理员无限制
	var userService UserService
	isAdmin, _ := userService.IsAdmin(userID)
	if isAdmin {
		return true, 999999, len(s.GetUserOnlineDevices(userID))
	}

	currentDevices := s.GetUserOnlineDevices(userID)
	return len(currentDevices) < user.MaxConcurrentDevices, user.MaxConcurrentDevices, len(currentDevices)
}

// kickOldestDevice 踢出最早登录的设备
func (s *OnlineDeviceService) kickOldestDevice(userID, keepCount int) {
	var sessions []model.UserOnlineSession
	model.UserOnlineSessionModel().
		Where("user_id = ?", userID).
		Order("login_time ASC").
		Find(&sessions)

	// 计算需要踢出的设备数量
	kickCount := len(sessions) - keepCount
	if kickCount <= 0 {
		return
	}

	// 踢出最早的设备
	for i := 0; i < kickCount && i < len(sessions); i++ {
		if err := s.KickDevice(userID, sessions[i].SessionID); err != nil {
			// 记录错误但继续处理其他设备
			fmt.Printf("踢出设备失败: %v\n", err)
		}
	}
}

// cleanExpiredSessions 清理过期会话
func (s *OnlineDeviceService) cleanExpiredSessions() {
	expiredTime := time.Now().Add(-5 * time.Minute)

	// 获取过期的会话
	var expiredSessions []model.UserOnlineSession
	model.UserOnlineSessionModel().
		Where("last_heartbeat < ?", expiredTime).
		Find(&expiredSessions)

	// 将过期的token加入黑名单
	var blackService BlacklistService
	for _, session := range expiredSessions {
		blackService.AddBlackByToken(session.SessionID)
	}

	// 删除过期会话
	model.UserOnlineSessionModel().
		Where("last_heartbeat < ?", expiredTime).
		Delete(&model.UserOnlineSession{})
}

// GenerateSessionID 生成会话ID（基于token的哈希）
func (s *OnlineDeviceService) GenerateSessionID(token string) string {
	hash := md5.Sum([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// IsSessionValid 检查会话是否有效
func (s *OnlineDeviceService) IsSessionValid(sessionID string) bool {
	var count int64
	model.UserOnlineSessionModel().
		Where("session_id = ?", sessionID).
		Count(&count)
	return count > 0
}

// GetUserMaxDevices 获取用户最大设备数限制
func (s *OnlineDeviceService) GetUserMaxDevices(userID int) int {
	var user model.Users
	if err := model.User().Where("id = ?", userID).First(&user).Error; err != nil {
		return 1 // 默认值
	}

	// 管理员无限制
	var userService UserService
	isAdmin, _ := userService.IsAdmin(userID)
	if isAdmin {
		return 999999
	}

	return user.MaxConcurrentDevices
}

// SetUserMaxDevices 设置用户最大设备数限制（管理员功能）
func (s *OnlineDeviceService) SetUserMaxDevices(userID, maxDevices int) error {
	return model.User().
		Where("id = ?", userID).
		Update("max_concurrent_devices", maxDevices).Error
}

// GetOnlineUsersCount 获取在线用户总数
func (s *OnlineDeviceService) GetOnlineUsersCount() int64 {
	s.cleanExpiredSessions()

	var count int64
	model.UserOnlineSessionModel().
		Select("DISTINCT user_id").
		Count(&count)
	return count
}
