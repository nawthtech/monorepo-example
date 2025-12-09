package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// Implementations use services interfaces from services package.
// We'll provide simplified implementations that return example data.

type authHandler struct {
	auth services.AuthService
}

type userHandler struct {
	user services.UserService
}

type serviceHandler struct {
	service services.ServiceService
}

func NewAuthHandler(a services.AuthService) *authHandler {
	return &authHandler{auth: a}
}
func NewUserHandler(u services.UserService) *userHandler {
	return &userHandler{user: u}
}
func NewServiceHandler(s services.ServiceService) *serviceHandler {
	return &serviceHandler{service: s}
}

// Auth handlers (simplified)
func (h *authHandler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "register (stub)"})
}
func (h *authHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "login (stub)"})
}
func (h *authHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "refresh (stub)"})
}

// User handlers
func (h *userHandler) GetProfile(c *gin.Context) {
	// in real impl: get user id from claims & call h.user.GetProfile
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":        "user_123",
			"name":      "نواف",
			"email":     "user@example.com",
			"createdAt": time.Now(),
		},
	})
}
func (h *userHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "profile updated (stub)"})
}
func (h *userHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "password changed (stub)"})
}
func (h *userHandler) GetUserStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{"total_orders": 5, "total_spent": 1200.5},
	})
}

// Service handlers (stubbed)
func (h *serviceHandler) GetServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": []gin.H{
			{"id": "svc_1", "title": "خدمة 1", "price": 100.0},
			{"id": "svc_2", "title": "خدمة 2", "price": 250.0},
		},
	})
}
func (h *serviceHandler) GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"id": id, "title": "مثال خدمة", "price": 199.99}})
}
func (h *serviceHandler) GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []string{"تطوير الويب", "تصميم"}})
}