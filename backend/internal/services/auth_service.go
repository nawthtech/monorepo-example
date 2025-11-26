package services

// AuthService يقدم خدمات المصادقة والأمان
type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// يمكن إضافة الدوال اللازمة هنا عند الحاجة