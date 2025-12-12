package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/db"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/router"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/slack"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	if err := Run(); err != nil {
		logger.Error(context.Background(), "server failed", "error", err)
		os.Exit(1)
	}
}

func Run() error {
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة الـ logger مع معلومات التطبيق
	initLogger(cfg)

	logger.Info(context.Background(), "🚀 Starting NawthTech Backend Server",
		logger.ComponentAttr("server"),
		logger.StatusAttr("starting"),
		logger.MetricAttr("port", float64(parsePort(cfg.Port)), ""),
		"environment", cfg.Environment,
		"version", cfg.Version,
		"debug", cfg.Debug)

	// تسجيل معلومات النظام
	logSystemInfo()

	// تهيئة Slack client
	initSlack(cfg)

	// تهيئة قاعدة البيانات
	database, err := initDatabase(cfg)
	if err != nil {
		logger.Error(context.Background(), "❌ Failed to initialize database", 
			logger.ErrAttr(err),
			logger.ComponentAttr("database"))
		return err
	}
	defer closeDatabase(database)

	// تسجيل نجاح الاتصال بقاعدة البيانات
	logger.Health(context.Background(), "database", "healthy", 0,
		logger.MetricAttr("connection_time_ms", 0, "ms"))

	// إنشاء service container
	serviceContainer, zapLogger := initServices(database, cfg)
	defer closeServices(serviceContainer)

	// تشغيل فحص الصحة الأولي
	runInitialHealthCheck(serviceContainer)

	// تكوين تطبيق Gin
	app := setupGinApp(cfg, database, serviceContainer)

	// بدء الخادم مع graceful shutdown
	return startServer(app, cfg, serviceContainer)
}

func parsePort(port string) int {
	var portNum int
	fmt.Sscanf(port, "%d", &portNum)
	if portNum == 0 {
		portNum = 8080
	}
	return portNum
}

func initLogger(cfg *config.Config) {
	logger.Init(cfg.Environment, cfg.AppName, cfg.Version)
	
	logger.Info(context.Background(), "📝 Logger initialized",
		logger.ComponentAttr("logger"),
		"environment", cfg.Environment,
		"format", getLogFormat(cfg.Environment))
}

func getLogFormat(env string) string {
	if env == "development" {
		return "text"
	}
	return "json"
}

func logSystemInfo() {
	logger.Info(context.Background(), "💻 System Information",
		logger.SystemAttr(),
		logger.ComponentAttr("system"))
}

func initSlack(cfg *config.Config) {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	appName := cfg.AppName
	environment := cfg.Environment

	if slackToken != "" && slackChannel != "" {
		client, err := slack.New(
			slack.WithToken(slackToken),
			slack.WithChannelURL(slackChannel),
			slack.WithAppName(appName),
			slack.WithEnvironment(environment),
		)
		if err != nil {
			logger.Warn(context.Background(), "⚠️ Failed to initialize Slack client", 
				logger.ErrAttr(err),
				logger.ComponentAttr("slack"))
		} else {
			logger.Info(context.Background(), "✅ Slack client initialized successfully",
				logger.ComponentAttr("slack"))
			
			// إرسال إشعار بدء التشغيل
			go func() {
				startTime := time.Now()
				_, _, err := client.SendAlert("info", "🚀 Backend Server Started", 
					fmt.Sprintf("%s backend server v%s has started successfully in %s environment", 
						appName, cfg.Version, environment))
				
				logger.Metric(context.Background(), "slack_notification_ms", 
					float64(time.Since(startTime).Milliseconds()),
					logger.ComponentAttr("slack"),
					logger.StatusAttr(func() string {
						if err != nil { return "failed" }; return "success"
					}()))
					
				if err != nil {
					logger.Warn(context.Background(), "⚠️ Failed to send Slack notification", 
						logger.ErrAttr(err))
				}
			}()
		}
	} else {
		logger.Info(context.Background(), "ℹ️ Slack not configured, running without notifications",
			logger.ComponentAttr("slack"))
	}
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	logger.Info(context.Background(), "🔌 Initializing database connection",
		logger.ComponentAttr("database"),
		"driver", cfg.Database.Driver,
		"name", cfg.Database.Name)

	startTime := time.Now()
	
	// تهيئة اتصال قاعدة البيانات
	database, err := db.InitializeFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	// التحقق من الاتصال
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		return nil, err
	}

	connectionTime := time.Since(startTime)
	
	logger.Health(context.Background(), "database", "healthy", connectionTime,
		logger.MetricAttr("connection_time_ms", float64(connectionTime.Milliseconds()), "ms"),
		logger.ComponentAttr("database"))

	// تشغيل عمليات الترحيل
	logger.Info(context.Background(), "🔄 Running database migrations",
		logger.ComponentAttr("database"))
	
	if err := db.RunMigrations(ctx, database); err != nil {
		logger.Warn(context.Background(), "⚠️ Migrations failed or already applied", 
			logger.ErrAttr(err),
			logger.ComponentAttr("database"))
	} else {
		logger.Info(context.Background(), "✅ Migrations completed successfully",
			logger.ComponentAttr("database"))
	}

	// التحقق من الجداول الأساسية
	checkTables(ctx, database)

	return database, nil
}

func checkTables(ctx context.Context, db *sql.DB) {
	query := `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.Warn(context.Background(), "⚠️ Failed to check tables", 
			logger.ErrAttr(err),
			logger.ComponentAttr("database"))
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	logger.Metric(context.Background(), "database_tables", float64(len(tables)),
		logger.ComponentAttr("database"),
		"tables", tables)
}

func closeDatabase(db *sql.DB) {
	if db != nil {
		startTime := time.Now()
		if err := db.Close(); err != nil {
			logger.Error(context.Background(), "❌ Failed to close database", 
				logger.ErrAttr(err),
				logger.ComponentAttr("database"))
		} else {
			closeTime := time.Since(startTime)
			logger.Info(context.Background(), "✅ Database connection closed",
				logger.MetricAttr("close_time_ms", float64(closeTime.Milliseconds()), "ms"),
				logger.ComponentAttr("database"))
		}
	}
}

func initServices(db *sql.DB, cfg *config.Config) (*services.ServiceContainer, *zap.Logger) {
	logger.Info(context.Background(), "🛠️ Initializing services",
		logger.ComponentAttr("services"))

	// إنشاء zap logger للتكامل مع services
	zapLogger, _ := zap.NewProduction()
	if cfg.Environment == "development" {
		zapLogger, _ = zap.NewDevelopment()
	}

	serviceContainer := services.NewServiceContainerWithConfig(db, cfg, zapLogger)

	// اختبار الخدمات الأساسية
	testBasicServices(serviceContainer)

	return serviceContainer, zapLogger
}

func testBasicServices(sc *services.ServiceContainer) {
	ctx := context.Background()
	startTime := time.Now()
	
	// اختبار خدمة التخزين المؤقت
	testKey := "server_start_test"
	testValue := "nawthtech_backend_" + time.Now().Format(time.RFC3339)

	if err := sc.Cache.Set(testKey, testValue, 1*time.Minute); err != nil {
		logger.Warn(context.Background(), "⚠️ Cache service test failed", 
			logger.ErrAttr(err),
			logger.ComponentAttr("cache"))
	} else {
		cacheTime := time.Since(startTime)
		logger.Health(context.Background(), "cache", "healthy", cacheTime,
			logger.MetricAttr("response_time_ms", float64(cacheTime.Milliseconds()), "ms"),
			logger.ComponentAttr("cache"))
	}
}

func runInitialHealthCheck(sc *services.ServiceContainer) {
	ctx := context.Background()
	logger.Info(context.Background(), "🏥 Running initial health check",
		logger.ComponentAttr("health"))

	// فحص صحة الخدمات الأساسية
	healthReq := &services.HealthRequest{
		CheckDatabase: true,
		CheckCache:    true,
		CheckStorage:  false, // اختياري في البداية
		CheckServices: true,
	}

	result, err := sc.Health.CheckHealth(ctx, healthReq)
	if err != nil {
		logger.Warn(context.Background(), "⚠️ Initial health check failed",
			logger.ErrAttr(err),
			logger.ComponentAttr("health"))
		return
	}

	logger.Health(context.Background(), "system", result.Status, 0,
		logger.MetricAttr("checks_count", float64(len(result.Checks)), ""),
		"checks", result.Checks,
		logger.ComponentAttr("health"))
}

func closeServices(sc *services.ServiceContainer) {
	if sc != nil {
		logger.Info(context.Background(), "🧹 Closing services",
			logger.ComponentAttr("services"))
		
		if err := sc.Close(); err != nil {
			logger.Warn(context.Background(), "⚠️ Error closing services",
				logger.ErrAttr(err),
				logger.ComponentAttr("services"))
		}
	}
}

func setupGinApp(cfg *config.Config, db *sql.DB, serviceContainer *services.ServiceContainer) *gin.Engine {
	// تكوين Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()

	// Middleware الأساسية
	app.Use(gin.Recovery())
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RateLimitMiddleware(cfg))
	app.Use(middleware.MetricsMiddleware())

	// مسارات الصحة
	setupHealthRoutes(app, db, serviceContainer)

	// API Routes
	setupAPIRoutes(app, serviceContainer)

	// مسارات غير موجودة
	app.NoRoute(func(c *gin.Context) {
		logger.Warn(context.Background(), "🔍 Endpoint not found",
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusNotFound, 0),
			logger.ComponentAttr("router"))
		
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "endpoint not found",
			"path":    c.Request.URL.Path,
		})
	})

	return app
}

func setupHealthRoutes(app *gin.Engine, db *sql.DB, serviceContainer *services.ServiceContainer) {
	// مسارات الصحة الأساسية
	app.GET("/health", func(c *gin.Context) {
		startTime := time.Now()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			logger.Health(context.Background(), "database", "unhealthy", time.Since(startTime),
				logger.ErrAttr(err),
				logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusServiceUnavailable, time.Since(startTime)))
			
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"message": "database connection failed",
				"error":   err.Error(),
			})
			return
		}

		logger.Health(context.Background(), "api", "healthy", time.Since(startTime),
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusOK, time.Since(startTime)))
		
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "nawthtech-backend",
			"database":  "connected",
		})
	})

	app.GET("/health/ready", func(c *gin.Context) {
		startTime := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
			"message":   "service is ready to accept requests",
		})
		
		logger.Health(context.Background(), "readiness", "ready", time.Since(startTime),
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusOK, time.Since(startTime)))
	})

	app.GET("/health/live", func(c *gin.Context) {
		startTime := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"status":    "live",
			"timestamp": time.Now().UTC(),
			"message":   "service is alive",
		})
		
		logger.Health(context.Background(), "liveness", "live", time.Since(startTime),
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusOK, time.Since(startTime)))
	})

	// مسار صحة متقدم (يحتاج مصادقة)
	app.GET("/health/detailed", middleware.AuthMiddleware(), func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()

		healthResult, err := serviceContainer.Health.GetDetailedHealth(ctx)
		if err != nil {
			logger.Error(context.Background(), "❌ Failed to get detailed health",
				logger.ErrAttr(err),
				logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusInternalServerError, time.Since(startTime)))
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get health information",
			})
			return
		}

		logger.Health(context.Background(), "detailed_health", healthResult.Status, time.Since(startTime),
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusOK, time.Since(startTime)))
		
		c.JSON(http.StatusOK, healthResult)
	})

	// مسار المقاييس (للمراقبة)
	app.GET("/health/metrics", middleware.AuthMiddleware(), func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()

		metrics, err := serviceContainer.Health.GetMetrics(ctx)
		if err != nil {
			logger.Warn(context.Background(), "⚠️ Failed to get metrics",
				logger.ErrAttr(err),
				logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusInternalServerError, time.Since(startTime)))
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get metrics",
			})
			return
		}

		logger.Metric(context.Background(), "metrics_collection_ms", 
			float64(time.Since(startTime).Milliseconds()),
			logger.RequestAttr(c.Request.Method, c.Request.URL.Path, http.StatusOK, time.Since(startTime)))
		
		c.JSON(http.StatusOK, metrics)
	})
}

func setupAPIRoutes(app *gin.Engine, serviceContainer *services.ServiceContainer) {
	// Group لـ API v1
	apiV1 := app.Group("/api/v1")

	// إنشاء Handlers
	handlerContainer := handlers.NewHandlerContainer(serviceContainer)

	// Register routes
	routes.RegisterV1Routes(apiV1, handlerContainer, middleware.AuthMiddleware())
}

func startServer(app *gin.Engine, cfg *config.Config, serviceContainer *services.ServiceContainer) error {
	// إعداد الخادم
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// بدء الخادم في goroutine
	serverErr := make(chan error, 1)
	go func() {
		logger.Info(context.Background(), "🌐 Server starting",
			logger.ComponentAttr("server"),
			"address", srv.Addr,
			"environment", cfg.Environment,
			"port", cfg.Port)

		if cfg.IsProduction() && cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
			logger.Info(context.Background(), "🔒 Starting server with TLS",
				logger.ComponentAttr("server"))
			if err := srv.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		} else {
			logger.Info(context.Background(), "🌍 Starting server without TLS",
				logger.ComponentAttr("server"))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		}
	}()

	// انتظار إشارة الإغلاق
	select {
	case err := <-serverErr:
		logger.Error(context.Background(), "💥 Server error",
			logger.ErrAttr(err),
			logger.ComponentAttr("server"))
		return err
	case sig := <-quit:
		logger.Info(context.Background(), "🛑 Received shutdown signal",
			"signal", sig.String(),
			logger.ComponentAttr("server"))
		
		// تسجيل مقاييس قبل الإغلاق
		logFinalMetrics(serviceContainer)
		
		// إعداد context للإغلاق
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		
		// محاولة إرسال إشعار Slack
		sendShutdownNotification(cfg)
		
		// إغلاق الخادم
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(context.Background(), "❌ Server shutdown failed", 
				logger.ErrAttr(err),
				logger.ComponentAttr("server"))
			return err
		}
		
		logger.Info(context.Background(), "✅ Server shutdown completed",
			logger.ComponentAttr("server"))
		return nil
	}
}

func logFinalMetrics(serviceContainer *services.ServiceContainer) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// الحصول على إحصائيات الخدمات
	if stats, err := serviceContainer.Health.GetServiceStats(ctx); err == nil {
		logger.Metric(context.Background(), "final_user_count", float64(stats.TotalUsers),
			logger.ComponentAttr("metrics"),
			"type", "total_users")
		logger.Metric(context.Background(), "final_service_count", float64(stats.TotalServices),
			logger.ComponentAttr("metrics"),
			"type", "total_services")
		logger.Metric(context.Background(), "final_revenue", stats.TotalRevenue,
			logger.ComponentAttr("metrics"),
			"type", "total_revenue")
	}
}

func sendShutdownNotification(cfg *config.Config) {
	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	
	if slackToken != "" && slackChannel != "" {
		client, err := slack.New(
			slack.WithToken(slackToken),
			slack.WithChannelURL(slackChannel),
			slack.WithAppName(cfg.AppName),
			slack.WithEnvironment(cfg.Environment),
		)
		if err == nil {
			startTime := time.Now()
			_, _, err := client.SendAlert("warning", "🛑 Backend Server Shutdown",
				fmt.Sprintf("%s backend server v%s is shutting down gracefully", 
					cfg.AppName, cfg.Version))
			
			logger.Metric(context.Background(), "slack_shutdown_notification_ms", 
				float64(time.Since(startTime).Milliseconds()),
				logger.ComponentAttr("slack"),
				logger.StatusAttr(func() string {
					if err != nil { return "failed" }; return "success"
				}()))
			
			if err != nil {
				logger.Warn(context.Background(), "⚠️ Failed to send shutdown notification", 
					logger.ErrAttr(err),
					logger.ComponentAttr("slack"))
			}
		}
	}
}