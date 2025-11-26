package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterUploadsRoutes(router *gin.RouterGroup, uploadsHandler *UploadsHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc, sellerMiddleware gin.HandlerFunc) {
	uploadsRoutes := router.Group("/uploads")
	uploadsRoutes.Use(authMiddleware)
	
	// ==================== رفع الملفات العامة ====================
	uploadsRoutes.POST("", uploadsHandler.UploadFile)
	uploadsRoutes.POST("/multiple", uploadsHandler.UploadMultipleFiles)
	
	// ==================== الذكاء الاصطناعي والتحليل ====================
	uploadsRoutes.POST("/ai-analysis", uploadsHandler.UploadAndAnalyze)
	uploadsRoutes.POST("/optimize/:fileId", uploadsHandler.OptimizeImage)
	
	// ==================== إدارة الملفات ====================
	uploadsRoutes.GET("/info/:fileId", uploadsHandler.GetFileInfo)
	uploadsRoutes.GET("/user", uploadsHandler.GetUserFiles)
	uploadsRoutes.PUT("/:fileId", uploadsHandler.UpdateFileInfo)
	uploadsRoutes.DELETE("/:fileId", uploadsHandler.DeleteFile)
	
	// ==================== استخدام التخزين ====================
	uploadsRoutes.GET("/storage/usage", uploadsHandler.GetStorageUsage)
	
	// ==================== رفع ملفات المتجر ====================
	sellerUploadsRoutes := uploadsRoutes.Group("")
	sellerUploadsRoutes.Use(sellerMiddleware)
	sellerUploadsRoutes.POST("/bulk", uploadsHandler.BulkUploadFiles)
	sellerUploadsRoutes.POST("/services/images", uploadsHandler.UploadServiceImage)
	sellerUploadsRoutes.POST("/services/gallery", uploadsHandler.UploadServiceGallery)
	sellerUploadsRoutes.POST("/store/assets", uploadsHandler.UploadStoreAsset)
	
	// ==================== الإحصائيات والإدارة ====================
	adminUploadsRoutes := uploadsRoutes.Group("")
	adminUploadsRoutes.Use(adminMiddleware)
	adminUploadsRoutes.GET("/stats", uploadsHandler.GetUploadStats)
	adminUploadsRoutes.POST("/cleanup", uploadsHandler.CleanupFiles)
}