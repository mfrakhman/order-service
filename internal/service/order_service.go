package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ordersvc/internal/models"
	"ordersvc/internal/rabbitmq"
	"ordersvc/internal/repository"
	"time"
)

type OrderService struct {
	repo *repository.OrderRepository
	rmq  *rabbitmq.RabbitMQ
	productAPI string
}

func NewOrderService(repo *repository.OrderRepository, rmq *rabbitmq.RabbitMQ, productAPI string) *OrderService {
	return &OrderService{
		repo: repo,
		rmq:  rmq,
		productAPI: productAPI,
	}
}

func (s *OrderService) CreateOrder(productId string, quantity int) (*models.Order, error) {
	// ✅ Fetch product info from Product-Service
	product, err := s.fetchProduct(productId)
	if err != nil {
		return nil, fmt.Errorf("product not found or unavailable: %w", err)
	}

	totalPrice := product.Price * quantity
	order := &models.Order{
		ProductID:  productId,
		TotalPrice: totalPrice,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
	}

	err = s.repo.Create(order)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"orderId":   order.ID,
		"productId": order.ProductID,
		"quantity":  quantity,
		"totalPrice": totalPrice,
		"createdAt": order.CreatedAt,
	}

	s.rmq.Publish("order.exchange", "order.created", payload)
	log.Printf("✅ Order created: %+v", order)

	return order, nil
}

func (s *OrderService) GetOrdersByProductID(productID string) ([]models.Order, error) {
	return s.repo.FindByProductID(productID)
}

// Fetch product from Product-Service
func (s *OrderService) fetchProduct(productId string) (*struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", s.productAPI, productId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch product: %s", string(body))
	}

	var product struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
		Qty   int    `json:"qty"`
	}
	json.NewDecoder(resp.Body).Decode(&product)
	return &product, nil
}
