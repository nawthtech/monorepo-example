// backend/internal/middleware/ai_metrics.go
package middleware

func AIMetrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        // تسجيل metrics
        metrics.RecordAIRequest(
            c.GetString("user_id"),
            c.Request.URL.Path,
            time.Since(start),
            c.GetFloat64("ai_cost"),
        )
    }
}