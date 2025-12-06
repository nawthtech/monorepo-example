package router

import (
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/handlers"
    "github.com/nawthtech/nawthtech/backend/internal/middleware"
)

func NewRouter() *gin.Engine {
    // Create router with default middleware
    router := gin.Default()
    
    // Add middleware
    router.Use(middleware.CORS())
    router.Use(middleware.Logger())
    
    // Setup routes
    handlers.SetupRoutes(router)
    
    return router
}
