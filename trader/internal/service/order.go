package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"trader/internal/model"
	"trader/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo          *repository.OrderRepository
	orderbookURL       string
	orderbookAPIKeyID  string
	orderbookAPIKey    string
	orderbookAPISecret string
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	orderbookURL, apiKeyID, apiKey, apiSecret string,
) *OrderService {
	return &OrderService{
		orderRepo:          orderRepo,
		orderbookURL:       orderbookURL,
		orderbookAPIKeyID:  apiKeyID,
		orderbookAPIKey:    apiKey,
		orderbookAPISecret: apiSecret,
	}
}

type PlaceOrderInput struct {
	Symbol   string `json:"symbol"   binding:"required"`
	Side     string `json:"side"     binding:"required,oneof=bid ask"`
	Type     string `json:"type"     binding:"required,oneof=limit market"`
	Price    string `json:"price"`
	Quantity string `json:"quantity" binding:"required"`
}

func (s *OrderService) PlaceOrder(ctx context.Context, userID string, in PlaceOrderInput) (map[string]interface{}, error) {
	fmt.Println("URL>>>", s.orderbookURL)
	if s.orderbookAPIKey == "" || s.orderbookAPISecret == "" {
		return nil, fmt.Errorf("orderbook API credentials not configured")
	}

	fmt.Println("URL>>>", s.orderbookURL)
	reqBody := PlaceOrderRequest{
		Symbol:   in.Symbol,
		Side:     in.Side,
		Type:     in.Type,
		Price:    in.Price,
		Quantity: in.Quantity,
	}

	respBody, statusCode, err := callOrderbook(
		s.orderbookURL,
		s.orderbookAPIKeyID,
		s.orderbookAPIKey,
		s.orderbookAPISecret,
		http.MethodPost,
		"/api/v1/orders/",
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("orderbook call failed: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("invalid orderbook response: %w", err)
	}

	fmt.Println("result::", result)

	// Ekstrak order_code dari response orderbook
	orderCode := extractOrderCode(result)

	status := "submitted"
	if statusCode != http.StatusCreated && statusCode != http.StatusOK {
		status = "failed"
	}

	// Persist order record locally
	order := &model.Order{
		ID:          uuid.New().String(),
		UserID:      userID,
		OrderCode:   orderCode,
		Symbol:      in.Symbol,
		Side:        in.Side,
		Type:        in.Type,
		Price:       in.Price,
		Quantity:    in.Quantity,
		Status:      status,
		RawResponse: respBody,
	}
	_ = s.orderRepo.Create(ctx, order)

	if statusCode != http.StatusCreated && statusCode != http.StatusOK {
		return result, fmt.Errorf("orderbook returned status %d", statusCode)
	}
	return result, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]model.Order, error) {
	return s.orderRepo.FindByUserID(ctx, userID)
}

func (s *OrderService) GetOrderByCode(ctx context.Context, userID, orderCode string) (*model.Order, error) {
	return s.orderRepo.FindByOrderCode(ctx, userID, orderCode)
}

func extractOrderCode(result map[string]interface{}) string {
	keys := []string{"order_code", "code", "orderCode", "order_id", "id"}
	for _, k := range keys {
		if v, ok := result[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	// Coba satu level nested: result["data"]["order_code"]
	if data, ok := result["data"].(map[string]interface{}); ok {
		for _, k := range keys {
			if v, ok := data[k]; ok {
				if s, ok := v.(string); ok && s != "" {
					return s
				}
			}
		}
	}
	return ""
}
