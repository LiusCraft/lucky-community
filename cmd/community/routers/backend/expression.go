package backend

import (
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/utils/page"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitExpressionRouters(r *gin.Engine) {
	group := r.Group("/community/admin/expression")
	group.GET("", listExpressions)
	group.POST("", saveExpression, middleware.OperLogger())
	group.PUT("/:id", updateExpression, middleware.OperLogger())
	group.DELETE("/:id", deleteExpression, middleware.OperLogger())
	group.PUT("/:id/toggle", toggleExpressionStatus, middleware.OperLogger())
}

// listExpressions 获取表情类型列表
func listExpressions(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)

	reactionService := services.NewReactionService(ctx)
	expressions, count, err := reactionService.PageExpressionTypes(p, limit)
	if err != nil {
		log.Errorf("获取表情类型列表失败: %v", err)
		result.Err("获取表情类型列表失败").Json(ctx)
		return
	}
	
	result.Page(expressions, count, nil).Json(ctx)
}

// saveExpression 创建表情类型
func saveExpression(ctx *gin.Context) {
	var expression model.ExpressionType
	if err := ctx.ShouldBindJSON(&expression); err != nil {
		log.Warnf("用户id: %d 添加表情类型参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(expression, err)).Json(ctx)
		return
	}

	reactionService := services.NewReactionService(ctx)
	created, err := reactionService.CreateExpressionType(&expression)
	if err != nil {
		log.Warnf("用户id: %d 添加表情类型失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	
	result.OkWithMsg(created, "表情类型创建成功").Json(ctx)
}

// updateExpression 更新表情类型
func updateExpression(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("用户id: %d 更新表情类型ID解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("ID参数错误").Json(ctx)
		return
	}

	var expression model.ExpressionType
	if err := ctx.ShouldBindJSON(&expression); err != nil {
		log.Warnf("用户id: %d 更新表情类型参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(expression, err)).Json(ctx)
		return
	}

	expression.ID = id
	reactionService := services.NewReactionService(ctx)
	err = reactionService.UpdateExpressionType(&expression)
	if err != nil {
		log.Warnf("用户id: %d 更新表情类型失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.OkWithMsg(nil, "表情类型更新成功").Json(ctx)
}

// deleteExpression 删除表情类型
func deleteExpression(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("用户id: %d 删除表情类型ID解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("ID参数错误").Json(ctx)
		return
	}

	reactionService := services.NewReactionService(ctx)
	
	// 检查是否有使用该表情的回复
	hasReactions, err := reactionService.CheckExpressionInUse(id)
	if err != nil {
		log.Warnf("用户id: %d 检查表情使用状态失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("检查表情使用状态失败").Json(ctx)
		return
	}
	
	if hasReactions {
		log.Warnf("用户id: %d 删除表情类型失败,该表情正在使用中", middleware.GetUserId(ctx))
		result.Err("删除失败,该表情正在使用中").Json(ctx)
		return
	}

	err = reactionService.DeleteExpressionType(id)
	if err != nil {
		log.Warnf("用户id: %d 删除表情类型失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.OkWithMsg(nil, "表情类型删除成功").Json(ctx)
}

// toggleExpressionStatus 切换表情启用状态
func toggleExpressionStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("用户id: %d 切换表情状态ID解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("ID参数错误").Json(ctx)
		return
	}

	reactionService := services.NewReactionService(ctx)
	newStatus, err := reactionService.ToggleExpressionStatus(id)
	if err != nil {
		log.Warnf("用户id: %d 切换表情状态失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}

	statusText := "禁用"
	if newStatus {
		statusText = "启用"
	}

	result.OkWithMsg(map[string]interface{}{
		"isActive": newStatus,
	}, "表情已"+statusText).Json(ctx)
}