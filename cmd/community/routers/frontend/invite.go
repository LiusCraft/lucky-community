package frontend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	inviteService services.InviteService
)

// 邀请相关的请求结构体
type InviteRecordsQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type CustomInviteCodeRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

// InitInviteRouters 初始化邀请相关路由
func InitInviteRouters(r *gin.Engine) {
	inviteGroup := r.Group("/community/user/invite")
	{
		inviteGroup.GET("/code", getUserInviteCode)          // 获取用户邀请码
		inviteGroup.GET("/stats", getInviterStatistics)      // 获取邀请统计
		inviteGroup.GET("/records", getInviteRecords)        // 获取邀请记录
		inviteGroup.GET("/ranking", getInviteRanking)        // 获取邀请排行榜
		inviteGroup.PUT("/code/custom", setCustomInviteCode) // 设置自定义邀请码
		inviteGroup.POST("/validate", validateInviteCode)    // 验证邀请码
	}
}

// getUserInviteCode 获取用户邀请码
// @Summary 获取用户邀请码
// @Description 获取当前用户的邀请码信息
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=model.UserInviteCode}
// @Failure 500 {object} result.Result
// @Router /api/user/invite/code [get]
func getUserInviteCode(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 获取用户邀请码，如果不存在则自动生成
	inviteCode, err := inviteService.GetOrCreateUserInviteCode(uid)
	if err != nil {
		result.Err("获取邀请码失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(inviteCode, "").Json(c)
}

// getInviterStatistics 获取邀请统计
// @Summary 获取邀请统计
// @Description 获取当前用户的邀请统计信息
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/user/invite/stats [get]
func getInviterStatistics(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 获取邀请统计
	stats, err := inviteService.GetInviterStatistics(uid)
	if err != nil {
		result.Err("获取邀请统计失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(stats, "").Json(c)
}

// getInviteRecords 获取邀请记录
// @Summary 获取邀请记录
// @Description 分页获取当前用户的邀请记录
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result{data=page.PageResult}
// @Failure 500 {object} result.Result
// @Router /api/user/invite/records [get]
func getInviteRecords(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 绑定查询参数
	var query InviteRecordsQuery
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

	// 获取邀请记录
	records, total, err := inviteService.GetInviteRelationsByInviter(uid, query.Page, query.PageSize)
	if err != nil {
		result.Err("获取邀请记录失败: " + err.Error()).Json(c)
		return
	}

	// 构造分页结果
	pageResult := map[string]interface{}{
		"list":     records,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	}

	result.Ok(pageResult, "").Json(c)
}

// getInviteRanking 获取邀请排行榜
// @Summary 获取邀请排行榜
// @Description 获取邀请系统的排行榜
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Param type query string false "排行榜类型 successful_invites|total_points_earned|total_invites"
// @Param limit query int false "排行榜数量，默认10"
// @Success 200 {object} result.Result{data=[]model.UserInviteCode}
// @Failure 500 {object} result.Result
// @Router /api/user/invite/ranking [get]
func getInviteRanking(c *gin.Context) {
	// 直接返回排行榜接口的简化版本 - 排行榜功能已被简化
	result.Ok([]interface{}{}, "排行榜功能已简化").Json(c)
}

// setCustomInviteCode 设置自定义邀请码
// @Summary 设置自定义邀请码
// @Description 设置用户的自定义邀请码
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Param request body CustomInviteCodeRequest true "自定义邀请码请求"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/user/invite/code/custom [put]
func setCustomInviteCode(c *gin.Context) {
	// 通过中间件获取用户ID
	_ = middleware.GetUserId(c)

	// 直接返回错误 - 自定义邀请码功能已被简化
	result.Err("自定义邀请码功能已简化").Json(c)
}

// validateInviteCode 验证邀请码
// @Summary 验证邀请码
// @Description 验证邀请码是否有效
// @Tags 邀请系统
// @Accept json
// @Produce json
// @Param invite_code query string true "邀请码"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/user/invite/validate [post]
func validateInviteCode(c *gin.Context) {
	// 获取邀请码参数
	inviteCode := c.Query("invite_code")
	if inviteCode == "" {
		result.Err("邀请码不能为空").Json(c)
		return
	}

	// 验证邀请码
	inviteCodeInfo, err := inviteService.ValidateInviteCode(inviteCode)
	if err != nil {
		result.Err("邀请码验证失败: " + err.Error()).Json(c)
		return
	}

	// 返回邀请码信息（隐藏敏感信息）
	response := map[string]interface{}{
		"valid":       true,
		"invite_code": inviteCodeInfo.InviteCode,
		"inviter_id":  inviteCodeInfo.UserID,
	}

	result.Ok(response, "").Json(c)
}
