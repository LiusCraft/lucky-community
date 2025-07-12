package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	welfareService = services.WelfareService{}
)

// InitWelfareRouters 初始化福利路由
func InitWelfareRouters(r *gin.Engine) {
	group := r.Group("/community/welfare")
	{
		group.GET("/list", getWelfareList)
		group.GET("/:id", getWelfareDetail)
		group.GET("/tag/:tag", getWelfareListByTag)
	}
}

// getWelfareList 获取福利列表
func getWelfareList(c *gin.Context) {
	tag := c.Query("tag")
	
	// 获取活跃福利列表
	items, err := welfareService.GetWelfareListByTag(tag)
	
	// 添加调试日志
	fmt.Printf("福利列表请求 - tag: %s, 返回条数: %d, 错误: %v\n", tag, len(items), err)
	
	result.Auto(items, err).Json(c)
}

// getWelfareDetail 获取福利详情
func getWelfareDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		result.Err("无效的福利ID").Json(c)
		return
	}
	
	item, err := welfareService.GetWelfareByID(id)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	
	// 只返回活跃状态的福利
	if item.Status != "active" {
		result.Err("福利不存在").Json(c)
		return
	}
	
	result.Ok(item, "成功").Json(c)
}

// getWelfareListByTag 根据标签获取福利列表
func getWelfareListByTag(c *gin.Context) {
	tag := c.Param("tag")
	
	if tag == "" {
		result.Err("标签不能为空").Json(c)
		return
	}
	
	items, err := welfareService.GetWelfareListByTag(tag)
	result.Auto(items, err).Json(c)
}