package frontend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	exchangeService services.ExchangeService
)

// 兑换相关的请求结构体
type ExchangeRequestsQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Status   string `form:"status" json:"status"` // pending|processing|completed|cancelled
}

type CreateExchangeRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UserInfo  string `json:"user_info"` // 用户填写的联系信息等
}

// InitExchangeRouters 初始化兑换相关路由
func InitExchangeRouters(r *gin.Engine) {
	exchangeGroup := r.Group("/community/exchange")
	{
		exchangeGroup.GET("/products", getExchangeProducts)    // 获取积分商品列表
		exchangeGroup.GET("/products/:id", getExchangeProduct) // 获取单个商品详情
		exchangeGroup.POST("/request", createExchangeRequest)  // 提交兑换申请
		exchangeGroup.GET("/requests", getExchangeRequests)    // 获取用户兑换申请记录
		exchangeGroup.GET("/requests/:id", getExchangeRequest) // 获取单个兑换申请详情
	}
}

// getExchangeProducts 获取积分商品列表
// @Summary 获取积分商品列表
// @Description 获取所有启用的积分商品
// @Tags 兑换系统
// @Accept json
// @Produce json
// @Param reward_type query string false "奖励类型筛选"
// @Success 200 {object} result.Result{data=[]model.PointProduct}
// @Failure 500 {object} result.Result
// @Router /api/exchange/products [get]
func getExchangeProducts(c *gin.Context) {
	// 获取奖励类型筛选参数
	var products interface{}
	var err error

	products, err = exchangeService.GetActiveProducts()

	if err != nil {
		result.Err("获取商品列表失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(products, "").Json(c)
}

// getExchangeProduct 获取单个商品详情
// @Summary 获取单个商品详情
// @Description 根据商品ID获取商品详细信息
// @Tags 兑换系统
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} result.Result{data=model.PointProduct}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/exchange/products/{id} [get]
func getExchangeProduct(c *gin.Context) {
	// 获取商品ID参数
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		result.Err("商品ID格式错误").Json(c)
		return
	}

	// 获取商品信息
	product, err := exchangeService.GetProductByID(productID)
	if err != nil {
		result.Err("获取商品信息失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(product, "").Json(c)
}

// createExchangeRequest 提交兑换申请
// @Summary 提交兑换申请
// @Description 用户提交积分兑换申请，立即扣除积分
// @Tags 兑换系统
// @Accept json
// @Produce json
// @Param request body CreateExchangeRequest true "兑换申请"
// @Success 200 {object} result.Result{data=model.ExchangeRequest}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/exchange/request [post]
func createExchangeRequest(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 绑定请求参数
	var req CreateExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 验证兑换操作
	if err := exchangeService.ValidateExchangeOperation(uid, req.ProductID, req.Quantity); err != nil {
		result.Err("兑换验证失败: " + err.Error()).Json(c)
		return
	}

	// 创建兑换申请
	exchangeRequest, err := exchangeService.CreateExchangeRequest(uid, req.ProductID, req.Quantity, req.UserInfo)
	if err != nil {
		result.Err("创建兑换申请失败: " + err.Error()).Json(c)
		return
	}

	result.Ok(exchangeRequest, "").Json(c)
}

// getExchangeRequests 获取用户兑换申请记录
// @Summary 获取用户兑换申请记录
// @Description 分页获取当前用户的兑换申请记录
// @Tags 兑换系统
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param status query string false "申请状态筛选"
// @Success 200 {object} result.Result{data=page.PageResult}
// @Failure 500 {object} result.Result
// @Router /api/exchange/requests [get]
func getExchangeRequests(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 绑定查询参数
	var query ExchangeRequestsQuery
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

	// 获取兑换申请记录
	requests, total, err := exchangeService.GetUserExchangeRequests(uid, query.Page, query.PageSize, query.Status)
	if err != nil {
		result.Err("获取兑换申请记录失败: " + err.Error()).Json(c)
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

// getExchangeRequest 获取单个兑换申请详情
// @Summary 获取单个兑换申请详情
// @Description 根据申请ID获取兑换申请详细信息
// @Tags 兑换系统
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Success 200 {object} result.Result{data=model.ExchangeRequest}
// @Failure 400 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /api/exchange/requests/{id} [get]
func getExchangeRequest(c *gin.Context) {
	// 通过中间件获取用户ID
	uid := middleware.GetUserId(c)

	// 获取申请ID参数
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		result.Err("申请ID格式错误").Json(c)
		return
	}

	// 获取兑换申请
	exchangeRequest, err := exchangeService.GetExchangeRequestByID(requestID)
	if err != nil {
		result.Err("获取兑换申请失败: " + err.Error()).Json(c)
		return
	}

	// 验证申请是否属于当前用户
	if exchangeRequest.UserID != uid {
		result.Err("无权访问该兑换申请").Json(c)
		return
	}

	result.Ok(exchangeRequest, "").Json(c)
}
