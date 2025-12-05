package ai

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// CostManager مدير التكاليف والحصص
type CostManager struct {
    mu           sync.RWMutex
    dataPath     string
    monthlyLimit float64
    dailyLimit   float64
    
    // إحصائيات الاستخدام
    Usage struct {
        TotalCost      float64                     `json:"total_cost"`
        MonthlyCost    map[string]float64          `json:"monthly_cost"`    // الشهر -> التكلفة
        DailyCost      map[string]float64          `json:"daily_cost"`      // اليوم -> التكلفة
        UserUsage      map[string]*UserUsageStats  `json:"user_usage"`      // المستخدم -> الإحصائيات
        ProviderUsage  map[string]*ProviderStats   `json:"provider_usage"`  // المزود -> الإحصائيات
        LastReset      time.Time                   `json:"last_reset"`
    }
}

// UserUsageStats إحصائيات استخدام المستخدم
type UserUsageStats struct {
    UserID        string                       `json:"user_id"`
    Tier          string                       `json:"tier"`
    TotalCost     float64                      `json:"total_cost"`
    MonthlyCost   map[string]float64           `json:"monthly_cost"`
    DailyCost     map[string]float64           `json:"daily_cost"`
    Quotas        map[string]*Quota            `json:"quotas"`
    LastActive    time.Time                    `json:"last_active"`
}

// ProviderStats إحصائيات المزود
type ProviderStats struct {
    ProviderName  string                       `json:"provider_name"`
    TotalRequests int64                        `json:"total_requests"`
    TotalCost     float64                      `json:"total_cost"`
    SuccessRate   float64                      `json:"success_rate"`
    AvgLatency    float64                      `json:"avg_latency"`
    LastUsed      time.Time                    `json:"last_used"`
}

// Quota حصة المستخدم
type Quota struct {
    Type          string    `json:"type"`           // text, image, video, audio
    Used          int64     `json:"used"`           // الكمية المستخدمة
    Limit         int64     `json:"limit"`          // الحد الأقصى
    ResetPeriod   string    `json:"reset_period"`   // daily, weekly, monthly
    LastReset     time.Time `json:"last_reset"`
}

// NewCostManager إنشاء CostManager جديد
func NewCostManager() (*CostManager, error) {
    dataPath := os.Getenv("AI_DATA_PATH")
    if dataPath == "" {
        dataPath = "./data/ai"
    }
    
    cm := &CostManager{
        dataPath:     dataPath,
        monthlyLimit: 0.0,  // 0 = لا يوجد حد (مجاني)
        dailyLimit:   0.0,
    }
    
    // تهيئة البيانات
    cm.Usage.MonthlyCost = make(map[string]float64)
    cm.Usage.DailyCost = make(map[string]float64)
    cm.Usage.UserUsage = make(map[string]*UserUsageStats)
    cm.Usage.ProviderUsage = make(map[string]*ProviderStats)
    cm.Usage.LastReset = time.Now()
    
    // تحميل البيانات المحفوظة
    if err := cm.load(); err != nil {
        fmt.Printf("Warning: Could not load cost data: %v\n", err)
    }
    
    // جدولة إعادة التعيين التلقائي
    go cm.startAutoReset()
    
    return cm, nil
}

// startAutoReset بدء إعادة التعيين التلقائي
func (cm *CostManager) startAutoReset() {
    // إعادة تعيين يومية في منتصف الليل
    go func() {
        for {
            now := time.Now()
            nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
            durationUntilMidnight := nextMidnight.Sub(now)
            
            time.Sleep(durationUntilMidnight)
            cm.resetDailyQuotas()
        }
    }()
    
    // إعادة تعيين شهرية في أول يوم من الشهر
    go func() {
        for {
            now := time.Now()
            nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
            durationUntilNextMonth := nextMonth.Sub(now)
            
            time.Sleep(durationUntilNextMonth)
            cm.resetMonthlyQuotas()
        }
    }()
}

// RecordUsage تسجيل استخدام
func (cm *CostManager) RecordUsage(record *UsageRecord) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    
    now := time.Now()
    monthKey := now.Format("2006-01")
    dayKey := now.Format("2006-01-02")
    
    // تحديث التكلفة الإجمالية
    cm.Usage.TotalCost += record.Cost
    cm.Usage.MonthlyCost[monthKey] += record.Cost
    cm.Usage.DailyCost[dayKey] += record.Cost
    
    // تحديث إحصائيات المستخدم
    if record.UserID != "" {
        if _, exists := cm.Usage.UserUsage[record.UserID]; !exists {
            cm.Usage.UserUsage[record.UserID] = &UserUsageStats{
                UserID:      record.UserID,
                Tier:        record.UserTier,
                MonthlyCost: make(map[string]float64),
                DailyCost:   make(map[string]float64),
                Quotas:      cm.getDefaultQuotas(record.UserTier),
            }
        }
        
        userStats := cm.Usage.UserUsage[record.UserID]
        userStats.TotalCost += record.Cost
        userStats.MonthlyCost[monthKey] += record.Cost
        userStats.DailyCost[dayKey] += record.Cost
        userStats.LastActive = now
        
        // تحديث الحصص
        if quota, exists := userStats.Quotas[record.Type]; exists {
            quota.Used += record.Quantity
        }
    }
    
    // تحديث إحصائيات المزود
    if _, exists := cm.Usage.ProviderUsage[record.Provider]; !exists {
        cm.Usage.ProviderUsage[record.Provider] = &ProviderStats{
            ProviderName: record.Provider,
        }
    }
    
    providerStats := cm.Usage.ProviderUsage[record.Provider]
    providerStats.TotalRequests++
    providerStats.TotalCost += record.Cost
    providerStats.LastUsed = now
    
    // تحديث متوسط زمن الاستجابة
    if record.Latency > 0 {
        if providerStats.AvgLatency == 0 {
            providerStats.AvgLatency = record.Latency
        } else {
            providerStats.AvgLatency = (providerStats.AvgLatency*float64(providerStats.TotalRequests-1) + record.Latency) / float64(providerStats.TotalRequests)
        }
    }
    
    // تحديث نسبة النجاح
    if record.Success {
        providerStats.SuccessRate = (providerStats.SuccessRate*float64(providerStats.TotalRequests-1) + 1.0) / float64(providerStats.TotalRequests)
    } else {
        providerStats.SuccessRate = (providerStats.SuccessRate * float64(providerStats.TotalRequests-1)) / float64(providerStats.TotalRequests)
    }
    
    // التحقق من تجاوز الحدود
    if cm.monthlyLimit > 0 && cm.Usage.MonthlyCost[monthKey] > cm.monthlyLimit {
        return fmt.Errorf("monthly cost limit exceeded: %.2f/%.2f", 
            cm.Usage.MonthlyCost[monthKey], cm.monthlyLimit)
    }
    
    if cm.dailyLimit > 0 && cm.Usage.DailyCost[dayKey] > cm.dailyLimit {
        return fmt.Errorf("daily cost limit exceeded: %.2f/%.2f", 
            cm.Usage.DailyCost[dayKey], cm.dailyLimit)
    }
    
    // الحفظ التلقائي
    go cm.save()
    
    return nil
}

// CanUseAI التحقق من إمكانية استخدام AI
func (cm *CostManager) CanUseAI(userID, requestType string) (bool, string) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    // إذا لم يكن هناك حدود، السماح دائماً
    if cm.monthlyLimit == 0 && cm.dailyLimit == 0 {
        return true, ""
    }
    
    now := time.Now()
    monthKey := now.Format("2006-01")
    dayKey := now.Format("2006-01-02")
    
    // التحقق من الحدود العامة
    if cm.monthlyLimit > 0 && cm.Usage.MonthlyCost[monthKey] >= cm.monthlyLimit {
        return false, "تم تجاوز الحد الشهري للت