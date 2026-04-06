package app

import (
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/internal/validation"
	"nineshop-be/pkg/auth"
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

var registerValidatorOnce sync.Once

func NewApplication(cfg *config.Config, dbQueries *sqlc.Queries) *Application {
	r := gin.Default()
	r.Static("/uploads", "./uploads")
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	redisClient := config.NewRedisConfig()
	redisCacheService := cache.NewRedisCacheService(redisClient)

	emailConfig := config.NewEmailConfig()
	emailService := email.NewGmailService(*emailConfig)
	userRepo := repository.NewUserRepository(dbQueries)
	userService := service.NewUserService(userRepo, redisClient, emailService)

	tokenService := auth.NewJwtService(redisCacheService)
	authService := service.NewAuthService(userRepo, tokenService, redisCacheService)

	module := []Module{
		NewUserModule(userService),
		NewAuthModule(authService, userService),
		NewProductModule(),
		NewCategoryModule(),
		NewBrandModule(),
		NewProductImagesModule(),
		NewCartModule(),
		NewPendingOrderModule(),
		NewPaymentModule(),
		NewOrderItemModule(),
		NewOrderModule(),
		NewMediaModule(),
	}

	registerValidatorOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			validation.RegisterValidation(v)
		}
	})

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
