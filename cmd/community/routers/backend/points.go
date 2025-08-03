package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	adminPointsService      services.PointsService
	adminPointConfigService services.PointConfigService
)

// 管理员积分查询相关的请求结构体
type AdminPointRecordsQuery struct {
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
	UserID    int    `form:"user_id" json:"user_id"`       // 按用户ID筛选
	Type      string `form:"type" json:"type"`             // earn 或 spend
	DateStart string `form:"date_start" json:"date_start"` // 开始日期
	DateEnd   string `form:"date_end" json:"date_end"`     // 结束日期
}

type AdminPointsStatisticsQuery struct {
	DateStart string `form:"date_start" json:"date_start"` // 开始日期
	DateEnd   string `form:"date_end" json:"date_end"`     // 结束日期
}

// 积分配置更新请求
type UpdatePointConfigRequest struct {
	RulesDescription   string `json:"rules_description" binding:"required"` // 积分规则说明
	InviteRewardPoints int    `json:"invite_reward_points" binding:"min=0"` // 邀请奖励积分
}

// 手动发放积分请求
type ManualGrantPointsRequest struct {
	UserID      int    `json:"user_id" binding:"required,min=1"`     // 用户ID
	Points      int    `json:"points" binding:"required,min=1"`      // 发放积分数量
	Description string `json:"description" binding:"required,min=1"` // 发放原因
}

// InitAdminPointsRouters 初始化管理员积分相关路由
func InitAdminPointsRouters(r *gin.Engine) {
	adminPointsGroup := r.Group("/community/admin/points")
	{
		adminPointsGroup.GET("/statistics", getAdminPointsStatistics) // 获取积分系统统计
		adminPointsGroup.GET("/records", getAdminPointRecords)        // 获取所有积分记录
		adminPointsGroup.GET("/users/:id", getAdminUserPoints)        // 获取指定用户积分详情

		// 积分配置管理
		adminPointsGroup.GET("/config", getPointConfig)    // 获取积分配置
		adminPointsGroup.PUT("/config", updatePointConfig) // 更新积分配置

		// 手动发放积分
		adminPointsGroup.POST("/manual-grant", manualGrantPoints) // 手动发放积分
	}
}

// getAdminPointsStatistics 获取积分系统统计
// @Summary 获取积分系统统计
// @Description 获取积分系统的整体统计数据，包括总发放积分、总消费积分、活跃用户等
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Param date_start query string false "开始日期 YYYY-MM-DD"
// @Param date_end query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/admin/points/statistics [get]
func getAdminPointsStatistics(c *gin.Context) {
	// 绑定查询参数
	var query AdminPointsStatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 获取积分系统统计
	statistics, err := adminPointsService.GetSystemPointsStatistics(query.DateStart, query.DateEnd)
	if err != nil {
		result.Err("获取积分统计失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(statistics, "").Json(c)
}

// getAdminPointRecords 获取所有积分记录
// @Summary 获取所有积分记录
// @Description 分页获取所有用户的积分变动记录，支持多种筛选条件
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param user_id query int false "用户ID筛选"
// @Param type query string false "记录类型 earn|spend"
// @Param date_start query string false "开始日期 YYYY-MM-DD"
// @Param date_end query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/admin/points/records [get]
func getAdminPointRecords(c *gin.Context) {
	// 绑定查询参数
	var query AdminPointRecordsQuery
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

	// 获取积分记录
	records, total, err := adminPointsService.GetAllPointRecords(
		query.Page, query.PageSize, query.UserID, query.Type, query.DateStart, query.DateEnd)
	if err != nil {
		result.Err("获取积分记录失败: " + err.Error()).Json(c)
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

// getAdminUserPoints 获取指定用户积分详情
// @Summary 获取指定用户积分详情
// @Description 获取指定用户的积分账户信息和最近的积分记录
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/points/users/{id} [get]
func getAdminUserPoints(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		result.Err("用户ID格式错误").Json(c)
		return
	}

	// 获取用户积分信息
	userPoints, err := adminPointsService.GetUserPoints(userID)
	if err != nil {
		result.Err("获取用户积分失败: " + err.Error()).Json(c)
		return
	}

	// 获取用户最近的积分记录 (最近20条)
	recentRecords, _, err := adminPointsService.GetUserPointRecords(userID, 1, 20, "")
	if err != nil {
		result.Err("获取用户积分记录失败: " + err.Error()).Json(c)
		return
	}

	// 构造返回数据
	response := map[string]interface{}{
		"user_points":    userPoints,
		"balance_info":   userPoints.GetPointsBalance(),
		"recent_records": recentRecords,
	}

	result.Ok(response, "").Json(c)
}

// getPointConfig 获取积分配置
// @Summary 获取积分配置
// @Description 获取当前的积分规则配置
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=model.PointConfig}
// @Failure 500 {object} result.Result
// @Router /api/admin/points/config [get]
func getPointConfig(c *gin.Context) {
	config, err := adminPointConfigService.GetPointConfig()
	if err != nil {
		result.Err("获取积分配置失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(config, "").Json(c)
}

// updatePointConfig 更新积分配置
// @Summary 更新积分配置
// @Description 更新积分规则说明和邀请奖励积分配置
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Param request body UpdatePointConfigRequest true "更新请求"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/points/config [put]
func updatePointConfig(c *gin.Context) {
	var req UpdatePointConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	err := adminPointConfigService.UpdatePointConfig(req.RulesDescription, req.InviteRewardPoints)
	if err != nil {
		result.Err("更新积分配置失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(nil, "积分配置更新成功").Json(c)
}

// manualGrantPoints 手动发放积分
// @Summary 手动发放积分
// @Description 管理员手动给指定用户发放积分
// @Tags 管理员-积分系统
// @Accept json
// @Produce json
// @Param request body ManualGrantPointsRequest true "发放请求"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/points/manual-grant [post]
func manualGrantPoints(c *gin.Context) {
	var req ManualGrantPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	err := adminPointsService.ManualGrantPoints(req.UserID, req.Points, req.Description)
	if err != nil {
		result.Err("发放积分失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(nil, "积分发放成功").Json(c)
}
