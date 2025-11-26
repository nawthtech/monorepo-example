package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("جلب سلة المستخدم", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب السلة بنجاح",
		"data": map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":       "cart_item_1",
					"serviceId": "service_1",
					"name":      "خدمة متابعين إنستغرام",
					"price":     150.00,
					"quantity":  2,
					"subtotal":  300.00,
				},
			},
			"total":     300.00,
			"itemCount": 1,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var cartData struct {
		ServiceID string                 `json:"serviceId"`
		Quantity  int                    `json:"quantity"`
		Options   map[string]interface{} `json:"options"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&cartData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إضافة عنصر إلى السلة", 
		"userID", userID, 
		"serviceID", cartData.ServiceID, 
		"quantity", cartData.Quantity, 
		"options", cartData.Options)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إضافة العنصر إلى السلة بنجاح",
		"data": map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":       "cart_item_new",
					"serviceId": cartData.ServiceID,
					"name":      "خدمة مضافة حديثاً",
					"price":     150.00,
					"quantity":  cartData.Quantity,
					"subtotal":  float64(cartData.Quantity) * 150.00,
				},
			},
			"total":     float64(cartData.Quantity) * 150.00,
			"itemCount": 1,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	itemID := chi.URLParam(r, "itemId")
	
	var updateData struct {
		Quantity int `json:"quantity"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث عنصر في السلة", 
		"userID", userID, 
		"itemID", itemID, 
		"quantity", updateData.Quantity)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث العنصر بنجاح",
		"data": map[string]interface{}{
			"updated": true,
			"itemId":  itemID,
			"newQuantity": updateData.Quantity,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	itemID := chi.URLParam(r, "itemId")

	logger.Stdout.Info("إزالة عنصر من السلة", "userID", userID, "itemID", itemID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إزالة العنصر من السلة بنجاح",
		"data": map[string]interface{}{
			"removed": true,
			"itemId":  itemID,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("تفريغ السلة", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تفريغ السلة بنجاح",
		"data": map[string]interface{}{
			"cleared":   true,
			"itemCount": 0,
			"total":     0.00,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) GetCartSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("جلب ملخص السلة", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب ملخص السلة بنجاح",
		"data": map[string]interface{}{
			"itemCount": 2,
			"subtotal":  450.00,
			"shipping":  0.00,
			"tax":       45.00,
			"total":     495.00,
			"discount":  0.00,
		},
	}

	respondJSON(w, response)
}

func (h *CartHandler) ValidateCartItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	logger.Stdout.Info("التحقق من صحة عناصر السلة", "userID", userID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم التحقق من صحة السلة بنجاح",
		"data": map[string]interface{}{
			"valid": true,
			"issues": []string{},
			"warnings": []string{
				"بعض الخدمات قد تتطلب وقتاً إضافياً للتجهيز",
			},
		},
	}

	respondJSON(w, response)
}