package main

import (
	"context"
	"log"
	"net/http"

	"trader/internal/config"
	"trader/internal/handler"
	"trader/internal/middleware"
	"trader/internal/repository"
	"trader/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────────────────────
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}
	log.Println("connected to postgres")

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	orderSvc := service.NewOrderService(
		orderRepo,
		cfg.OrderbookURL,
		cfg.OrderbookAPIKeyID,
		cfg.OrderbookAPIKey,
		cfg.OrderbookAPISecret,
	)

	// ── Handlers ──────────────────────────────────────────────────────────────
	h := handler.New(authSvc, orderSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "trader"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	orders := r.Group("/orders")
	orders.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		orders.POST("", h.PlaceOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/:order_code", h.GetOrder) // ← lookup by order_code
	}

	log.Printf("trader service listening on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
