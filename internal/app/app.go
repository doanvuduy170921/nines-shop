package app

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/internal/validation"
	"nineshop-be/pkg/cache"
	"nineshop-be/pkg/email"
)

type Module interface {
	Route() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	redis  *redis.Client
	email  email.EmailService
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	redisClient := config.NewRedisConfig()
	redisCacheService := cache.NewRedisCacheService(redisClient)

	emailConfig := config.NewEmailConfig()
	emailService := email.NewGmailService(*emailConfig)
	userService := service.NewUserService(repository.NewUserRepository(db.DB), redisClient, emailService)

	module := []Module{
		NewUserModule(redisClient, emailService),
		NewAuthModule(redisCacheService, userService),
		NewProductModule(),
		NewCategoryModule(),
		NewBrandModule(),
		NewProductImagesModule(),
		NewCartModule(),
		NewPendingOrderModule(),
		NewPaymentModule(),
		NewOrderItemModule(),
	}

	// init validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterValidation(v) // <-- Đăng ký strong_pass, strong_user
	}

	routes.RegisterRoute(redisCacheService, r, getRoutes(module)...)
	return &Application{
		config: cfg,
		router: r,
		redis:  redisClient,
		email:  emailService,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getRoutes(modules []Module) []routes.Route {
	routerList := make([]routes.Route, len(modules))
	for i, m := range modules {
		routerList[i] = m.Route()
	}
	return routerList
}
