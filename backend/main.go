package main

import (
	"cmp"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nawthtech/backend/internal/config"
	"github.com/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/backend/internal/logger"
	"github.com/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/backend/internal/services"
	"github.com/nawthtech/backend/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()

	// تهيئة النظام
	if err := initializeSystem(cfg); err != nil {
		logger.Stdout.Error("فشل في تهيئة النظام", logger.ErrAttr(err))
		os.Exit(1)
	}

	// إنشاء الموجه
	r := mux.NewRouter()

	// تسجيل الوسائط
	middleware.Register(r)

	// تهيئة الخدمات
	services := initializeServices()

	// تسجيل المسارات
	handlers.Register(r, services)

	// إعداد الخادم
	server := &http.Server{
		Addr:              ":" + cmp.Or(os.Getenv("PORT"), "3000"),
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// بدء الخادم
	logger.Stdout.Info("بدء تشغيل الخادم", 
		slog.String("port", server.Addr),
		slog.String("environment", cfg.Environment),
		slog.String("version", cfg.Version),
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Stdout.Error("فشل في بدء الخادم", logger.ErrAttr(err))
		os.Exit(1)
	}
}

// initializeSystem تهيئة مكونات النظام
func initializeSystem(cfg *config.Config) error {
	// تهيئة قاعدة البيانات
	if err := utils.InitDatabase(cfg.DatabaseURL); err != nil {
		return err
	}

	// تهيئة المدقق
	if err := utils.InitValidator(); err != nil {
		return err
	}

	// تهيئة نظام التشفير
	if err := utils.InitEncryption(cfg.EncryptionKey); err != nil {
		return err
	}

	// التحقق من المساحات التخزينية
	if err := checkStorage(); err != nil {
		return err
	}

	logger.Stdout.Info("تم تهيئة النظام بنجاح")
	return nil
}

// initializeServices تهيئة جميع الخدمات
func initializeServices() *handlers.Services {
	// خدمات الإدارة
	adminService := services.NewAdminService()

	// خدمات المستخدمين
	userService := services.NewUserService()

	// خدمات المصادقة
	authService := services.NewAuthService()

	// خدمات المتجر
	storeService := services.NewStoreService()

	// خدمات السلة
	cartService := services.NewCartService()

	// خدمات المدفوعات
	paymentService := services.NewPaymentService()

	// خدمات الذكاء الاصطناعي
	aiService := services.NewAIService()

	// خدمات البريد الإلكتروني
	emailService := services.NewEmailService()

	// خدمات الرفع
	uploadService := services.NewUploadService()

	return &handlers.Services{
		Admin:   adminService,
		User:    userService,
		Auth:    authService,
		Store:   storeService,
		Cart:    cartService,
		Payment: paymentService,
		AI:      aiService,
		Email:   emailService,
		Upload:  uploadService,
	}
}

// checkStorage التحقق من المساحات التخزينية
func checkStorage() error {
	// التحقق من مساحة القرص
	// (يمكن إضافة منطق فعلي للتحقق من المساحة)
	
	// إنشاء المجلدات الضرورية إذا لم تكن موجودة
	dirs := []string{
		"./uploads",
		"./logs", 
		"./backups",
		"./temp",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// gracefulShutdown إغلاق النظام بشكل آمن
func gracefulShutdown(server *http.Server) {
	// إنشاء قناة للإشارات
	sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// انتظار الإشارة
	<-sigChan

	logger.Stdout.Info("استلام إشارة إغلاق، بدء الإغلاق الآمن")

	// إعطاء وقت للإغلاق الآمن
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// إيقاف الخادم
	if err := server.Shutdown(ctx); err != nil {
		logger.Stdout.Error("فشل في إيقاف الخادم بشكل آمن", logger.ErrAttr(err))
	} else {
		logger.Stdout.Info("تم إيقاف الخادم بشكل آمن")
	}

	// إغلاق اتصالات قاعدة البيانات
	utils.CloseDatabase()

	// إغلاق الموارد الأخرى
	// (يمكن إضافة إغلاق لموارد أخرى مثل اتصالات Redis، إلخ)
}
