package router

import (
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/handlers"
    "github.com/nawthtech/nawthtech/backend/internal/middleware"
)

func SetupRouter(aiHandler *handlers.AIHandler) *gin.Engine {
    r := gin.Default()
    
    // Middleware
    r.Use(middleware.CORS())
    r.Use(middleware.Logger())
    r.Use(middleware.Auth()) // JWT authentication
    
    // API routes
    api := r.Group("/api")
    {
        // AI Routes
        ai := api.Group("/ai")
        {
            // توليد محتوى نصي
            ai.POST("/generate", aiHandler.GenerateContentHandler)
            
            // تحليل الصور
            ai.POST("/analyze-image", aiHandler.AnalyzeImageHandler)
            
            // تحليل الاتجاهات
            ai.POST("/analyze-trends", aiHandler.AnalyzeTrendsHandler)
            
            // إنشاء استراتيجية
            ai.POST("/strategy", aiHandler.GenerateStrategyHandler)
            
            // الحصول على providers المتاحة
            ai.GET("/providers", aiHandler.GetProvidersHandler)
        }
        
        // User routes
        users := api.Group("/users")
        {
            users.POST("/register", handlers.RegisterHandler)
            users.POST("/login", handlers.LoginHandler)
            users.GET("/profile", handlers.GetProfileHandler)
        }
        
        // Dashboard routes
        dashboard := api.Group("/dashboard")
        {
            dashboard.GET("/metrics", handlers.GetMetricsHandler)
            dashboard.GET("/activities", handlers.GetActivitiesHandler)
        }
    }
    
    return r
}