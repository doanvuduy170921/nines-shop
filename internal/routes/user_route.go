package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
	"nineshop-be/internal/middleware"
)

type UserRoute struct {
	handler *handler.UserHandler
}

func NewUserRoute(handler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		handler: handler,
	}
}

func (ur *UserRoute) Register(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.GET("/get-all", middleware.RoleMiddleware("admin"), ur.handler.GetAllUser)
		user.GET("/get-all-v2", middleware.RoleMiddleware("admin"), ur.handler.GetAllByFilter)
		user.PUT("/soft-delete/:uuid", middleware.RoleMiddleware("admin"), ur.handler.SoftDelete)
		user.PUT("/update/:uuid", ur.handler.UpdateUser)
		user.GET("/get/:uuid", ur.handler.GetByUuid)
	}

}
