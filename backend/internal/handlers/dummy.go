package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func SetupRoutes(router *gin.Engine) {
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    
    // AI routes
    aiRoutes := router.Group("/api/ai")
    {
        aiRoutes.POST("/generate", generateAIHandler)
        aiRoutes.GET("/models", getModelsHandler)
    }
}

func generateAIHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "AI generated"})
}

func getModelsHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"models": []string{"gemini", "stability"}})
}
