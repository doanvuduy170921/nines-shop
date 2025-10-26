package routes

import (
	"github.com/gin-gonic/gin"
	"nineshop-be/internal/handler"
)

type AuthRoute struct {
	handler *handler.AuthHandler
}

func NewAuthRoute(handler *handler.AuthHandler) *AuthRoute {
	return &AuthRoute{
		handler: handler,
	}
}

func (ur *AuthRoute) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ur.handler.Login)
		auth.POST("refresh", ur.handler.RefreshToken)
		auth.POST("/logout", ur.handler.Logout)
		auth.POST("/create", ur.handler.CreateUser)
	}
}
