// backend/cmd/video-server/main.go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/gin-gonic/gin"
    "github.com/nawthtech/nawthtech/backend/internal/ai/video"
)

func main() {
    // إنشاء مزود فيديو هجين
    provider := video.NewHybridVideoProvider()
    
    r := gin.Default()
    
    r.POST("/api/video/generate", func(c *gin.Context) {
        var req struct {
            Prompt   string `json:"prompt"`
            ImageURL string `json:"image_url,omitempty"`
            UserID   string `json:"user_id"`
            Tier     string `json:"tier" default:"free"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        // تحقق من الحصة
        if !provider.CostManager().CanGenerateVideo(req.UserID, req.Tier) {
            c.JSON(403, gin.H{
                "error": "Monthly quota exceeded. Upgrade your plan.",
                "quota": provider.CostManager().GetRemainingQuota(req.UserID),
            })
            return
        }
        
        // توليد الفيديو
        videoData, err := provider.Generate(req.Prompt, req.ImageURL)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        // تسجيل الاستخدام
        provider.CostManager().RecordGeneration(req.UserID, req.Tier, 0)
        
        c.Header("Content-Type", "video/mp4")
        c.Data(200, "video/mp4", videoData)
    })
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    
    log.Printf("🎬 Free Video Generation Server running on port %s", port)
    log.Fatal(r.Run(":" + port))
}