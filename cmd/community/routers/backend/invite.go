package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/dao"
	services "xhyovo.cn/community/server/service"
)

var (
	adminInviteService services.InviteService
)

// 管理员邀请管理相关的请求结构体
type AdminInviteStatisticsQuery struct {
	DateStart string `form:"date_start" json:"date_start"` // 开始日期
	DateEnd   string `form:"date_end" json:"date_end"`     // 结束日期
}

type AdminInviteRelationsQuery struct {
	Page      int `form:"page" json:"page"`
	PageSize  int `form:"pageSize" json:"pageSize"`
	InviterID int `form:"inviter_id" json:"inviter_id"` // 邀请者ID筛选
	InviteeID int `form:"invitee_id" json:"invitee_id"` // 被邀请者ID筛选
}

// InitAdminInviteRouters 初始化管理员邀请管理路由
func InitAdminInviteRouters(r *gin.Engine) {
	adminInviteGroup := r.Group("/community/admin/invite")
	{
		adminInviteGroup.GET("/statistics", getAdminInviteStatistics) // 获取邀请系统统计
		adminInviteGroup.GET("/relations", getAdminInviteRelations)   // 获取邀请关系列表
		adminInviteGroup.GET("/users/:id", getAdminUserInviteInfo)    // 获取指定用户邀请信息
	}
}

// getAdminInviteStatistics 获取邀请系统统计
// @Summary 获取邀请系统统计
// @Description 获取邀请系统的整体统计数据，包括总邀请数、成功邀请数、积分发放等
// @Tags 管理员-邀请系统
// @Accept json
// @Produce json
// @Param date_start query string false "开始日期 YYYY-MM-DD"
// @Param date_end query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/admin/invite/statistics [get]
func getAdminInviteStatistics(c *gin.Context) {
	// 绑定查询参数
	var query AdminInviteStatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 获取邀请系统统计
	statistics, err := adminInviteService.GetInviteStatistics()
	if err != nil {
		result.Err("获取邀请统计失败: " + err.Error()).Json(c)
		return
	}

	// 添加积分奖励统计（可以扩展为从积分记录中统计）
	// 使用邀请关系数量 * 10（每次邀请10积分）来估算
	totalInvites := int64(0)
	if val, ok := statistics["total_invite_relations"].(int64); ok {
		totalInvites = val
	}
	statistics["total_points_awarded"] = totalInvites * 10

	result.Ok(statistics, "").Json(c)
}

// getAdminInviteRelations 获取邀请关系列表
// @Summary 获取邀请关系列表
// @Description 分页获取所有邀请关系，支持邀请者和被邀请者筛选
// @Tags 管理员-邀请系统
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param inviter_id query int false "邀请者ID筛选"
// @Param invitee_id query int false "被邀请者ID筛选"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/admin/invite/relations [get]
func getAdminInviteRelations(c *gin.Context) {
	// 绑定查询参数
	var query AdminInviteRelationsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	// 使用新的DAO方法获取包含用户昵称的邀请关系列表
	inviteDao := &dao.InviteRelationDao{}
	relations, total, err := inviteDao.GetInviteRelationsWithUserInfo(query.Page, query.PageSize, query.InviterID, query.InviteeID)
	if err != nil {
		result.Err("获取邀请关系失败: " + err.Error()).Json(c)
		return
	}

	// 构造分页结果
	pageResult := map[string]interface{}{
		"list":     relations,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	}

	result.Ok(pageResult, "").Json(c)
}

// getAdminUserInviteInfo 获取指定用户邀请信息
// @Summary 获取指定用户邀请信息
// @Description 获取指定用户的邀请码、邀请统计和邀请记录
// @Tags 管理员-邀请系统
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/invite/users/{id} [get]
func getAdminUserInviteInfo(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		result.Err("用户ID格式错误").Json(c)
		return
	}

	// 获取用户邀请码
	inviteCode, err := adminInviteService.GetUserInviteCode(userID)
	if err != nil {
		result.Err("获取用户邀请码失败: " + err.Error()).Json(c)
		return
	}

	// 获取用户邀请统计
	statistics, err := adminInviteService.GetInviterStatistics(userID)
	if err != nil {
		result.Err("获取用户邀请统计失败: " + err.Error()).Json(c)
		return
	}

	// 获取用户邀请记录 (最近20条)
	inviteRecords, _, err := adminInviteService.GetInviteRelationsByInviter(userID, 1, 20)
	if err != nil {
		result.Err("获取用户邀请记录失败: " + err.Error()).Json(c)
		return
	}

	// 构造返回数据
	response := map[string]interface{}{
		"invite_code":    inviteCode,
		"statistics":     statistics,
		"invite_records": inviteRecords,
	}

	result.Ok(response, "").Json(c)
}
