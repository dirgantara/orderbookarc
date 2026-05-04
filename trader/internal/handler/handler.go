package handler

import (
	"fmt"
	"log"
	"net/http"
	"trader/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth  *service.AuthService
	order *service.OrderService
}

func New(auth *service.AuthService, order *service.OrderService) *Handler {
	return &Handler{auth: auth, order: order}
}

func respond(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

func respondErr(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var in service.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respondErr(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.auth.Register(c.Request.Context(), in)
	if err != nil {
		respondErr(c, http.StatusConflict, err.Error())
		return
	}
	respond(c, http.StatusCreated, user)
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var in service.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respondErr(c, http.StatusBadRequest, err.Error())
		return
	}
	token, user, err := h.auth.Login(c.Request.Context(), in)
	if err != nil {
		respondErr(c, http.StatusUnauthorized, err.Error())
		return
	}
	respond(c, http.StatusOK, gin.H{"token": token, "user": user})
}

// POST /orders
func (h *Handler) PlaceOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	var in service.PlaceOrderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respondErr(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("user", userID)
	log.Printf("trader service listening on :%s", userID)
	result, err := h.order.PlaceOrder(c.Request.Context(), userID.(string), in)
	if err != nil {
		respondErr(c, http.StatusBadGateway, err.Error())
		return
	}
	respond(c, http.StatusCreated, result)
}

// GET /orders
func (h *Handler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	orders, err := h.order.ListOrders(c.Request.Context(), userID.(string))
	if err != nil {
		respondErr(c, http.StatusInternalServerError, err.Error())
		return
	}
	respond(c, http.StatusOK, orders)
}

// GET /orders/:order_code
func (h *Handler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	orderCode := c.Param("order_code")
	order, err := h.order.GetOrderByCode(c.Request.Context(), userID.(string), orderCode)
	if err != nil {
		respondErr(c, http.StatusNotFound, "order not found")
		return
	}
	respond(c, http.StatusOK, order)
}
