package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	adminExchangeService services.ExchangeService
)

// 管理员兑换管理相关的请求结构体
type AdminExchangeRequestsQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Status   string `form:"status" json:"status"` // pending|processing|completed|cancelled
}

type ProcessExchangeRequest struct {
	Status    string `json:"status" binding:"required"` // processing|completed|cancelled
	AdminNote string `json:"admin_note"`                // 管理员备注
}

type CreatePointProductRequest struct {
	Name        string `json:"name" binding:"required"`        // 商品名称
	Description string `json:"description"`                    // 商品描述
	PointsCost  int    `json:"points_cost" binding:"required"` // 积分成本
	RewardType  string `json:"reward_type" binding:"required"` // 奖励类型
	Stock       int    `json:"stock"`                          // 库存数量(-1表示无限)
	SortOrder   int    `json:"sort_order"`                     // 排序
	IsEnabled   bool   `json:"is_enabled"`                     // 是否启用
}

type UpdatePointProductRequest struct {
	Name        string `json:"name"`        // 商品名称
	Description string `json:"description"` // 商品描述
	PointsCost  int    `json:"points_cost"` // 积分成本
	RewardType  string `json:"reward_type"` // 奖励类型
	Stock       int    `json:"stock"`       // 库存数量
	SortOrder   int    `json:"sort_order"`  // 排序
	IsEnabled   *bool  `json:"is_enabled"`  // 是否启用(使用指针以区分false和nil)
}

// InitAdminExchangeRouters 初始化管理员兑换管理路由
func InitAdminExchangeRouters(r *gin.Engine) {
	adminExchangeGroup := r.Group("/community/admin/exchange")
	{
		// 兑换申请管理
		adminExchangeGroup.GET("/requests", getAdminExchangeRequests)    // 获取兑换申请列表
		adminExchangeGroup.GET("/requests/:id", getAdminExchangeRequest) // 获取单个兑换申请详情
		adminExchangeGroup.PUT("/requests/:id", processExchangeRequest)  // 处理兑换申请

		// 积分商品管理
		adminExchangeGroup.GET("/products", getAdminExchangeProducts)    // 获取商品列表
		adminExchangeGroup.GET("/products/:id", getAdminExchangeProduct) // 获取单个商品详情
		adminExchangeGroup.POST("/products", createPointProduct)         // 创建积分商品
		adminExchangeGroup.PUT("/products/:id", updatePointProduct)      // 更新积分商品
		adminExchangeGroup.DELETE("/products/:id", deletePointProduct)   // 删除积分商品
	}
}

// getAdminExchangeRequests 获取兑换申请列表
// @Summary 获取兑换申请列表
// @Description 分页获取所有兑换申请，支持状态和用户筛选
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param status query string false "申请状态筛选"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/requests [get]
func getAdminExchangeRequests(c *gin.Context) {
	// 绑定查询参数
	var query AdminExchangeRequestsQuery
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

	// 获取兑换申请列表
	requests, total, err := adminExchangeService.GetAllExchangeRequests(
		query.Page, query.PageSize, query.Status)
	if err != nil {
		result.Err("获取兑换申请列表失败: " + err.Error()).Json(c)
		return
	}

	// 构造分页结果
	pageResult := map[string]interface{}{
		"list":     requests,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	}

	result.Ok(pageResult, "").Json(c)
}

// getAdminExchangeRequest 获取单个兑换申请详情
// @Summary 获取单个兑换申请详情
// @Description 根据申请ID获取兑换申请详细信息
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} result.Result{data=model.ExchangeRequest}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/requests/{id} [get]
func getAdminExchangeRequest(c *gin.Context) {
	// 获取申请ID参数
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		result.Err("申请ID格式错误").Json(c)
		return
	}

	// 获取兑换申请
	exchangeRequest, err := adminExchangeService.GetExchangeRequestByID(requestID)
	if err != nil {
		result.Err("获取兑换申请失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(exchangeRequest, "").Json(c)
}

// processExchangeRequest 处理兑换申请
// @Summary 处理兑换申请
// @Description 管理员处理兑换申请，更新申请状态
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body ProcessExchangeRequest true "处理申请"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/requests/{id} [put]
func processExchangeRequest(c *gin.Context) {
	// 获取申请ID参数
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		result.Err("申请ID格式错误").Json(c)
		return
	}

	// 绑定请求参数
	var req ProcessExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 处理兑换申请
	if err := adminExchangeService.ProcessExchangeRequestStatus(requestID, req.Status); err != nil {
		result.Err("处理兑换申请失败: " + err.Error()).Json(c)
		return
	}

	result.Ok("处理成功", "").Json(c)
}

// getAdminExchangeProducts 获取商品列表
// @Summary 获取商品列表
// @Description 获取所有积分商品列表（包括已禁用的）
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Success 200 {object} result.Result{data=[]model.PointProduct}
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/products [get]
func getAdminExchangeProducts(c *gin.Context) {
	// 获取所有商品（包括已禁用的）
	products, err := adminExchangeService.GetAllProducts()
	if err != nil {
		result.Err("获取商品列表失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(products, "").Json(c)
}

// getAdminExchangeProduct 获取单个商品详情
// @Summary 获取单个商品详情
// @Description 根据商品ID获取商品详细信息
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} result.Result{data=model.PointProduct}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/products/{id} [get]
func getAdminExchangeProduct(c *gin.Context) {
	// 获取商品ID参数
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		result.Err("商品ID格式错误").Json(c)
		return
	}

	// 获取商品信息
	product, err := adminExchangeService.GetProductByID(productID)
	if err != nil {
		result.Err("获取商品信息失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(product, "").Json(c)
}

// createPointProduct 创建积分商品
// @Summary 创建积分商品
// @Description 创建新的积分商品
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param request body CreatePointProductRequest true "商品信息"
// @Success 200 {object} result.Result{data=model.PointProduct}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/products [post]
func createPointProduct(c *gin.Context) {
	// 绑定请求参数
	var req CreatePointProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 创建商品
	product, err := adminExchangeService.CreateProduct(
		req.Name, req.Description, req.PointsCost, req.RewardType,
		"", req.Stock, req.SortOrder, req.IsEnabled)
	if err != nil {
		result.Err("创建商品失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(product, "").Json(c)
}

// updatePointProduct 更新积分商品
// @Summary 更新积分商品
// @Description 更新积分商品信息
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param request body UpdatePointProductRequest true "商品信息"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/products/{id} [put]
func updatePointProduct(c *gin.Context) {
	// 获取商品ID参数
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		result.Err("商品ID格式错误").Json(c)
		return
	}

	// 绑定请求参数
	var req UpdatePointProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 更新商品
	if err := adminExchangeService.UpdateProduct(productID, req.Name, req.Description,
		req.PointsCost, req.RewardType, "", req.Stock, req.SortOrder, req.IsEnabled); err != nil {
		result.Err("更新商品失败: " + err.Error()).Json(c)
		return
	}

	result.Ok("更新成功", "").Json(c)
}

// deletePointProduct 删除积分商品
// @Summary 删除积分商品
// @Description 软删除积分商品（将is_enabled设为false）
// @Tags 管理员-兑换管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} result.Result
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/admin/exchange/products/{id} [delete]
func deletePointProduct(c *gin.Context) {
	// 获取商品ID参数
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		result.Err("商品ID格式错误").Json(c)
		return
	}

	// 软删除商品（设置为禁用状态）
	if err := adminExchangeService.DeleteProduct(productID); err != nil {
		result.Err("删除商品失败: " + err.Error()).Json(c)
		return
	}

	result.Ok("删除成功", "").Json(c)
}
