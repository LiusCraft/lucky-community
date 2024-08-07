package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

type Dashboard struct {
	UserCount    int64 `json:"userCount"`
	Profit       int64 `json:"profit"`
	ArticleCount int64 `json:"articleCount"`
}

func InitDashboardRouters(r *gin.Engine) {
	group := r.Group("/community/admin/dashboard")
	group.GET("", dashboard)
}
func dashboard(ctx *gin.Context) {
	// 查出用户数量

	var codes []model.InviteCodes
	model.InviteCode().Where("state = 1").Select("id", "member_id").Find(&codes)
	var userCount int64

	model.User().Count(&userCount)

	// 查出文章数量
	var articleCount int64
	model.Article().Count(&articleCount)
	// 查出盈利
	var orderService = services.OrderServices{}

	d := Dashboard{
		UserCount:    userCount,
		ArticleCount: articleCount,
		Profit:       orderService.CalculateProfit(),
	}
	result.Ok(d, "").Json(ctx)
}
