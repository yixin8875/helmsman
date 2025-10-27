package routers

import (
	"github.com/gin-gonic/gin"

	"helmsman/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		tradeTagsRouter(group, handler.NewTradeTagsHandler())
	})
}

func tradeTagsRouter(group *gin.RouterGroup, h handler.TradeTagsHandler) {
	g := group.Group("/tradeTags")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)                    // [post] /api/v1/tradeTags
	g.DELETE("/:tradeID", h.DeleteByTradeID) // [delete] /api/v1/tradeTags/:tradeID
	g.PUT("/:tradeID", h.UpdateByTradeID)    // [put] /api/v1/tradeTags/:tradeID
	g.GET("/:tradeID", h.GetByTradeID)       // [get] /api/v1/tradeTags/:tradeID
	g.POST("/list", h.List)                  // [post] /api/v1/tradeTags/list
	g.POST("/all", h.ListAll)                // [get] /api/v1/tradeTags/all
}
