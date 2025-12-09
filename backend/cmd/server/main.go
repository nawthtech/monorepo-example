package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/db"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// initLogger تهيئة logger
func initLogger() {
	if slog.Default().Handler() == nil {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(handler))
	}
}

// Run تشغيل خادم API
func Run() error {
	// تهيئة logger
	initLogger()

	// تحميل الإعدادات
	cfg := config.Load()

	// تسجيل بدء التشغيل
	slog.Info("🚀 بدء تشغيل خادم NawthTech API",
		"environment", cfg.Environment,
		"version", cfg.Version,
		"port", cfg.Port,
	)

	// تهيئة قاعدة بيانات D1
	db.InitializeD1(cfg)

	// تهيئة الحاوية للخدمات
	serviceContainer := services.NewServiceContainer(db.DB)

	// إنشاء تطبيق Gin
	app := initGinApp(cfg)

	// تسجيل الوسائط
	registerMiddlewares(app)

	// تسجيل المسارات
	registerRoutes(app, serviceContainer, cfg)

	// بدء الخادم
	return startServer(app, cfg)
}

// initGinApp تهيئة تطبيق Gin
func initGinApp(cfg *config.Config) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.New()
	app.ForwardedByClientIP = true
	app.MaxMultipartMemory = 10 << 20 // 10 MB

	if cfg.IsProduction() {
		app.SetTrustedProxies([]string{
			"127.0.0.1",
			"::1",
			"10.0.0.0/8",
			"172.16.0.0/12",
			"192.168.0.0/16",
		})
	} else {
		app.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	return app
}

// registerMiddlewares تسجيل الوسائط
func registerMiddlewares(app *gin.Engine) {
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.SecurityHeaders())
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		slog.Info("طلب HTTP",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
		)
		return ""
	}))
	app.Use(gin.Recovery())
	app.Use(middleware.RateLimitMiddlewareFunc())
}

// registerRoutes تسجيل المسارات
func registerRoutes(app *gin.Engine, svcContainer *services.ServiceContainer, cfg *config.Config) {
	handlerContainer := &handlers.HandlerContainer{
		Auth:    handlers.NewAuthHandler(svcContainer.Auth),
		User:    handlers.NewUserHandler(svcContainer.User),
		Service: handlers.NewServiceHandler(svcContainer.Service),
	}

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	handlers.RegisterV1Routes(v1Group, handlerContainer)

	// مسارات الصحة
	app.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		dbStatus := "healthy"
		if err := db.DB.Ping(ctx); err != nil {
			dbStatus = "unhealthy"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      dbStatus,
			"timestamp":   time.Now().UTC(),
			"version":     cfg.Version,
			"environment": cfg.Environment,
		})
	})

	app.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "live",
			"timestamp": time.Now().UTC(),
		})
	})

	app.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		dbStatus := "healthy"
		if err := db.DB.Ping(ctx); err != nil {
			dbStatus = "unhealthy"
		}

		if dbStatus != "healthy" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not_ready",
				"timestamp": time.Now().UTC(),
				"error":     "Database not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
		})
	})
}

// startServer بدء الخادم
func startServer(app *gin.Engine, cfg *config.Config) error {
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		slog.Info("🌐 بدء الخادم",
			"port", cfg.Port,
			"environment", cfg.Environment,
			"version", cfg.Version,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("❌ فشل في بدء الخادم", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-sigChan
	slog.Info("🛑 استلام إشارة إغلاق", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("⏳ إيقاف الخادم بشكل أنيق...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("❌ فشل في إيقاف الخادم", "error", err)
		return err
	}

	slog.Info("✅ تم إيقاف الخادم بنجاح")
	return nil
}

// main الدالة الرئيسية
func main() {
	initLogger()
	if err := Run(); err != nil {
		slog.Error("❌ فشل في تشغيل الخادم", "error", err)
		os.Exit(1)
	}
}