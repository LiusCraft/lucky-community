package frontend

import (
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	pointsService      services.PointsService
	pointConfigService services.PointConfigService
)

// 用户积分查询相关的请求结构体
type PointRecordsQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Type     string `form:"type" json:"type"` // earn 或 spend
}

// InitPointsRouters 初始化积分相关路由
func InitPointsRouters(r *gin.Engine) {
	pointsGroup := r.Group("/community/user/points")
	{
		pointsGroup.GET("", getUserPoints)            // 获取用户积分概览
		pointsGroup.GET("/records", getPointRecords)  // 获取用户积分记录
		pointsGroup.GET("/ranking", getPointsRanking) // 获取积分排行榜
	}

	// 积分规则（无需登录）
	r.GET("/community/points/rules", getPointRules) // 获取积分规则说明
}

// getUserPoints 获取用户积分概览
// @Summary 获取用户积分概览
// @Description 获取当前用户的积分余额、累计获得、累计消费等信息
// @Tags 积分系统
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=model.UserPoints}
// @Failure 500 {object} result.Result
// @Router /api/user/points [get]
func getUserPoints(c *gin.Context) {

	userId := middleware.GetUserId(c)

	// 获取用户积分信息
	userPoints, err := pointsService.GetUserPoints(userId)
	if err != nil {
		result.Err("获取用户积分失败: " + err.Error()).Json(c)
		return
	}

	// 构造返回数据，包含积分余额统计
	response := map[string]interface{}{
		"user_points":  userPoints,
		"balance_info": userPoints.GetPointsBalance(),
	}

	result.Ok(response, "").Json(c)
}

// getPointRecords 获取用户积分记录
// @Summary 获取用户积分记录
// @Description 分页获取当前用户的积分变动记录，支持按类型筛选
// @Tags 积分系统
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param type query string false "记录类型 earn|spend"
// @Success 200 {object} result.Result{data=page.PageResult}
// @Failure 500 {object} result.Result
// @Router /api/user/points/records [get]
func getPointRecords(c *gin.Context) {
	userId := middleware.GetUserId(c)
	// 绑定查询参数
	var query PointRecordsQuery
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
	records, total, err := pointsService.GetUserPointRecords(userId, query.Page, query.PageSize, query.Type)
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

// getPointsRanking 获取积分排行榜
// @Summary 获取积分排行榜
// @Description 获取积分系统的用户排行榜
// @Tags 积分系统
// @Accept json
// @Produce json
// @Param limit query int false "排行榜数量，默认10"
// @Success 200 {object} result.Result{data=[]model.UserPoints}
// @Failure 500 {object} result.Result
// @Router /api/user/points/ranking [get]
func getPointsRanking(c *gin.Context) {
	// 获取排行榜数量参数
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// 限制最大数量
	if limit > 100 {
		limit = 100
	}

	// 获取积分排行榜
	ranking, err := pointsService.GetUsersPointsRanking(limit)
	if err != nil {
		result.Err("获取积分排行榜失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(ranking, "").Json(c)
}

// getPointRules 获取积分规则说明
// @Summary 获取积分规则说明
// @Description 获取积分系统的规则说明，包括如何获得积分、使用规则等
// @Tags 积分系统
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/points/rules [get]
func getPointRules(c *gin.Context) {
	config, err := pointConfigService.GetPointConfig()
	if err != nil {
		result.Err("获取积分规则失败: " + err.Error()).Json(c)
		return
	}

	// 构造返回数据
	response := map[string]interface{}{
		"rules_description":    config.RulesDescription,
		"invite_reward_points": config.InviteRewardPoints,
	}

	result.Ok(response, "").Json(c)
}
