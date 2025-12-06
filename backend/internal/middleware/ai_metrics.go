package middleware

import (
    "fmt"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    aiRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_requests_total",
            Help: "Total number of AI requests",
        },
        []string{"endpoint", "status"},
    )
    
    aiRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "ai_request_duration_seconds",
            Help:    "Duration of AI requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"endpoint"},
    )
)

func init() {
    prometheus.MustRegister(aiRequestsTotal, aiRequestDuration)
}

// AIMetrics middleware لتتبع مقاييس AI
func AIMetrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := fmt.Sprintf("%d", c.Writer.Status())
        
        aiRequestsTotal.WithLabelValues(path, status).Inc()
        aiRequestDuration.WithLabelValues(path).Observe(duration)
    }
}
