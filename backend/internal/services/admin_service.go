package services

// AdminService يقدم خدمات إدارة النظام المتقدمة
type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// يمكن إضافة الدوال اللازمة هنا عند الحاجة