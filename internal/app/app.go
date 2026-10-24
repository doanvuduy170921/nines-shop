package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"nineshop-be/internal/config"
	"nineshop-be/internal/db/sqlc"
	"nineshop-be/internal/repository"
	"nineshop-be/internal/routes"
	"nineshop-be/internal/service"
	"nineshop-be/internal/validation"
	"nineshop-be/pkg/auth"
	"nineshop-be/pkg/cache"
	"nineshop-be/pkg/email"
	"sync"
	"time"
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

func NewApplication(ctx context.Context, cfg *config.Config, dbQueries *sqlc.Queries) *Application {
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
		NewPendingOrderModule(ctx, redisCacheService),
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

func (a *Application) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:              a.config.ServerAddress,
		Handler:           a.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("shutting down server...")
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func getRoutes(modules []Module) []routes.Route {
	routerList := make([]routes.Route, len(modules))
	for i, m := range modules {
		routerList[i] = m.Route()
	}
	return routerList
}
func RunMigrations(dbURL string) {

	log.Println("🔄 Running database migrations...")

	m, err := migrate.New("file://internal/db/migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Database migrations applied successfully!")
}
