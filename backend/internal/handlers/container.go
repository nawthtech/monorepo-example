package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// HandlerContainer حاوية لجميع الـ handlers
type HandlerContainer struct {
	Auth         *AuthHandler
	User         *UserHandler
	Service      *ServiceHandler
	Category     *CategoryHandler
	Order        *OrderHandler
	Payment      *PaymentHandler
	Upload       *UploadHandler
	Notification *NotificationHandler
	Admin        *AdminHandler
	Health       *HealthHandler
}

// NewHandlerContainer إنشاء حاوية handlers جديدة
func NewHandlerContainer(serviceContainer *services.ServiceContainer) *HandlerContainer {
	if serviceContainer == nil {
		return &HandlerContainer{}
	}

	container := &HandlerContainer{}

	// إنشاء الـ handlers مع التحقق من وجود الخدمات
	if serviceContainer.Auth != nil {
		container.Auth = NewAuthHandler(serviceContainer.Auth)
	}

	if serviceContainer.User != nil {
		container.User = NewUserHandler(serviceContainer.User)
	}

	if serviceContainer.Service != nil {
		container.Service = NewServiceHandler(serviceContainer.Service)
	}

	if serviceContainer.Category != nil {
		container.Category = NewCategoryHandler(serviceContainer.Category)
	}

	if serviceContainer.Order != nil {
		container.Order = NewOrderHandler(serviceContainer.Order)
	}

	if serviceContainer.Payment != nil {
		container.Payment = NewPaymentHandler(serviceContainer.Payment)
	}

	if serviceContainer.Upload != nil {
		container.Upload = NewUploadHandler(serviceContainer.Upload)
	}

	if serviceContainer.Notification != nil {
		container.Notification = NewNotificationHandler(serviceContainer.Notification)
	}

	if serviceContainer.Admin != nil {
		container.Admin = NewAdminHandler(serviceContainer.Admin)
	}

	if serviceContainer.Health != nil {
		container.Health = NewHealthHandler(serviceContainer.Health)
	}

	return container
}