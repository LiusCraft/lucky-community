package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	updateLogService = services.UpdateLogService{}
)

// InitUpdateLogRouters 初始化更新日志路由
func InitUpdateLogRouters(r *gin.Engine) {
	group := r.Group("/community/update-log")
	{
		group.GET("/list", getUpdateLogList)
		group.GET("/recent", getRecentUpdateLogs)
		group.GET("/:id", getUpdateLogDetail)
		group.GET("/type/:type", getUpdateLogListByType)
	}
}

// getUpdateLogList 获取更新日志列表
func getUpdateLogList(c *gin.Context) {
	var req model.UpdateLogListRequest

	// 解析查询参数
	if c.ShouldBindQuery(&req) != nil {
		req.Page = 1
		req.PageSize = 10
	}

	// 获取类型筛选
	req.Type = c.Query("type")

	// 获取活跃更新日志列表
	items, err := updateLogService.GetUpdateLogList(req)

	// 添加调试日志
	fmt.Printf("更新日志列表请求 - type: %s, page: %d, pageSize: %d, 返回条数: %d, 错误: %v\n",
		req.Type, req.Page, req.PageSize, len(items.Items), err)

	result.Auto(items, err).Json(c)
}

// getRecentUpdateLogs 获取最近的更新日志（用于首页展示）
func getRecentUpdateLogs(c *gin.Context) {
	// 获取限制数量参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20 // 限制最大数量
	}

	items, err := updateLogService.GetRecentUpdateLogs(limit)

	// 添加调试日志
	fmt.Printf("最近更新日志请求 - limit: %d, 返回条数: %d, 错误: %v\n", limit, len(items), err)

	result.Auto(items, err).Json(c)
}

// getUpdateLogDetail 获取更新日志详情
func getUpdateLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的更新日志ID").Json(c)
		return
	}

	item, err := updateLogService.GetUpdateLogByID(id)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	// 只返回活跃状态的更新日志
	if item.Status != "active" {
		result.Err("更新日志不存在").Json(c)
		return
	}

	result.Ok(item, "成功").Json(c)
}

// getUpdateLogListByType 根据类型获取更新日志列表
func getUpdateLogListByType(c *gin.Context) {
	logType := c.Param("type")

	// 验证类型参数长度
	if len(logType) > 50 {
		result.Err("更新日志类型长度不能超过50个字符").Json(c)
		return
	}

	items, err := updateLogService.GetUpdateLogListByType(logType)

	// 添加调试日志
	fmt.Printf("按类型获取更新日志请求 - type: %s, 返回条数: %d, 错误: %v\n", logType, len(items), err)

	result.Auto(items, err).Json(c)
}
