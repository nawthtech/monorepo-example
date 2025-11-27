package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// GenerateMD5Hash إنشاء تجزئة MD5
func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// ParseJSON تحليل JSON مع معالجة الأخطاء
func ParseJSON(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

// ToJSON تحويل البيانات إلى JSON
func ToJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// ToJSONString تحويل البيانات إلى سلسلة JSON
func ToJSONString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}