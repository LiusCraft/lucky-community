package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	welfareAdminService = services.WelfareService{}
)

// InitWelfareRouters 初始化福利管理路由
func InitWelfareRouters(r *gin.Engine) {
	group := r.Group("/admin/welfare")
	{
		group.GET("/list", getWelfareList)
		group.GET("/:id", getWelfareDetail)
		group.POST("", createWelfare)
		group.PUT("/:id", updateWelfare)
		group.DELETE("/:id", deleteWelfare)
		group.PATCH("/:id/status", updateWelfareStatus)
	}
}

// getWelfareList 获取福利列表（管理员）
func getWelfareList(c *gin.Context) {
	var req model.WelfareListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	response, err := welfareAdminService.GetAdminWelfareList(req)
	result.Auto(response, err).Json(c)
}

// getWelfareDetail 获取福利详情（管理员）
func getWelfareDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的福利ID").Json(c)
		return
	}

	item, err := welfareAdminService.GetWelfareByID(id)
	result.Auto(item, err).Json(c)
}

// createWelfare 创建福利
func createWelfare(c *gin.Context) {
	var req model.WelfareItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	item, err := welfareAdminService.CreateWelfare(req)
	result.Auto(item, err).Json(c)
}

// updateWelfare 更新福利
func updateWelfare(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的福利ID").Json(c)
		return
	}

	var req model.WelfareItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	item, err := welfareAdminService.UpdateWelfare(id, req)
	result.Auto(item, err).Json(c)
}

// deleteWelfare 删除福利
func deleteWelfare(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的福利ID").Json(c)
		return
	}

	err = welfareAdminService.DeleteWelfare(id)
	result.Auto("删除成功", err).Json(c)
}

// updateWelfareStatus 更新福利状态
func updateWelfareStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的福利ID").Json(c)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	err = welfareAdminService.UpdateWelfareStatus(id, req.Status)
	result.Auto("状态更新成功", err).Json(c)
}
