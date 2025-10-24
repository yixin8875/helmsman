package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"helmsman/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		usersRouter(group, handler.NewUsersHandler())
	})
}

func usersRouter(group *gin.RouterGroup, h handler.UsersHandler) {
	g := group.Group("/users")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)                   // [post] /api/v1/users/register
	g.POST("/", h.Create, middleware.Auth())          // [post] /api/v1/users
	g.DELETE("/:id", h.DeleteByID, middleware.Auth()) // [delete] /api/v1/users/:id
	g.PUT("/:id", h.UpdateByID, middleware.Auth())    // [put] /api/v1/users/:id
	g.GET("/:id", h.GetByID, middleware.Auth())       // [get] /api/v1/users/:id

	g.POST("/condition", h.GetByCondition, middleware.Auth()) // [post] /api/v1/users/condition
}
