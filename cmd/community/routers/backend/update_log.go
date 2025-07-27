package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	adminUpdateLogService = services.UpdateLogService{}
)

// InitAdminUpdateLogRouters 初始化管理员更新日志路由
func InitAdminUpdateLogRouters(r *gin.Engine) {
	group := r.Group("/community/admin/update-log")
	{
		group.GET("/list", adminGetUpdateLogList)
		group.GET("/:id", adminGetUpdateLogDetail)
		group.POST("/create", adminCreateUpdateLog)
		group.PUT("/:id", adminUpdateUpdateLog)
		group.DELETE("/:id", adminDeleteUpdateLog)
		group.PUT("/:id/status", adminUpdateUpdateLogStatus)
	}
}

// adminGetUpdateLogList 获取更新日志列表（管理员）
func adminGetUpdateLogList(c *gin.Context) {
	var req model.UpdateLogListRequest

	// 解析查询参数
	if c.ShouldBindQuery(&req) != nil {
		req.Page = 1
		req.PageSize = 10
	}

	// 获取类型和状态筛选
	req.Type = c.Query("type")
	req.Status = c.Query("status")

	// 获取更新日志列表（管理员可以查看所有状态）
	response, err := adminUpdateLogService.GetAdminUpdateLogList(req)

	result.Auto(response, err).Json(c)
}

// adminGetUpdateLogDetail 获取更新日志详情（管理员）
func adminGetUpdateLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的更新日志ID").Json(c)
		return
	}

	item, err := adminUpdateLogService.GetUpdateLogByID(id)
	result.Auto(item, err).Json(c)
}

// adminCreateUpdateLog 创建更新日志（管理员）
func adminCreateUpdateLog(c *gin.Context) {
	var req model.UpdateLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("请求参数格式错误: " + err.Error()).Json(c)
		return
	}

	item, err := adminUpdateLogService.CreateUpdateLog(req)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(item, "更新日志创建成功").Json(c)
}

// adminUpdateUpdateLog 更新更新日志（管理员）
func adminUpdateUpdateLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的更新日志ID").Json(c)
		return
	}

	var req model.UpdateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("请求参数格式错误: " + err.Error()).Json(c)
		return
	}

	item, err := adminUpdateLogService.UpdateUpdateLog(id, req)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(item, "更新日志更新成功").Json(c)
}

// adminDeleteUpdateLog 删除更新日志（管理员）
func adminDeleteUpdateLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的更新日志ID").Json(c)
		return
	}

	err = adminUpdateLogService.DeleteUpdateLog(id)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(nil, "更新日志删除成功").Json(c)
}

// adminUpdateUpdateLogStatus 更新更新日志状态（管理员）
func adminUpdateUpdateLogStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的更新日志ID").Json(c)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("请求参数格式错误: " + err.Error()).Json(c)
		return
	}

	err = adminUpdateLogService.UpdateUpdateLogStatus(id, req.Status)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(nil, "更新日志状态更新成功").Json(c)
}
