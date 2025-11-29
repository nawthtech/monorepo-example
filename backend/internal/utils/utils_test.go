package utils

import (
	"testing"
)

func TestGetMemoryUsageMB(t *testing.T) {
	mem := GetMemoryUsageMB()
	
	if mem.UsedMB < 0 {
		t.Errorf("Expected used memory to be non-negative, got %f", mem.UsedMB)
	}
	
	if mem.TotalMB < 0 {
		t.Errorf("Expected total memory to be non-negative, got %f", mem.TotalMB)
	}
	
	if mem.UsagePercentage < 0 || mem.UsagePercentage > 100 {
		t.Errorf("Expected usage percentage between 0 and 100, got %f", mem.UsagePercentage)
	}
}

func TestGetGoroutineCount(t *testing.T) {
	count := GetGoroutineCount()
	
	if count <= 0 {
		t.Errorf("Expected goroutine count to be positive, got %d", count)
	}
}

// إذا كان لديك دالة GenerateID في utils، أضف هذا الاختبار
// إذا لم تكن موجودة، احذف هذا الاختبار أو أنشئ الدالة أولاً

func TestGenerateRandomString(t *testing.T) {
	// اختبار دالة مساعدة إذا كانت موجودة
	str1 := generateRandomString(10)
	str2 := generateRandomString(10)
	
	if len(str1) != 10 {
		t.Errorf("Expected string length 10, got %d", len(str1))
	}
	
	if str1 == str2 {
		t.Errorf("Expected unique strings, got duplicates: %s", str1)
	}
}

// إذا لم تكن generateRandomString موجودة، أضف هذه الدالة المساعدة إلى utils.go:
/*
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
*/