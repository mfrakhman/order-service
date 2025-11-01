package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ordersvc/internal/cache"
	"ordersvc/internal/models"
	"ordersvc/internal/rabbitmq"
	"ordersvc/internal/repository"
	"time"
)

type OrderService struct {
	repo       *repository.OrderRepository
	rmq        *rabbitmq.RabbitMQ
	productAPI string
	cache      *cache.RedisCache
	orderQueue chan *orderQueueItem
	Quantity   int
}

type orderQueueItem struct {
	Order    *models.Order
	Quantity int
}

func NewOrderService(repo *repository.OrderRepository, rmq *rabbitmq.RabbitMQ, productAPI string, cache *cache.RedisCache) *OrderService {
	svc := &OrderService{
		repo:       repo,
		rmq:        rmq,
		productAPI: productAPI,
		cache:      cache,
		orderQueue: make(chan *orderQueueItem, 1000),
	}

	go svc.startOrderWorker()
	return svc
}

func (s *OrderService) CreateOrder(productId string, quantity int) (*models.Order, error) {
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

	select {
	case s.orderQueue <- &orderQueueItem{Order: order, Quantity: quantity}:
		log.Printf("Queued order for async processing: %+v", order)
	default:
		log.Println("Order queue full — dropping order!")
	}

	return order, nil
}

func (s *OrderService) startOrderWorker() {
	log.Println("Order processing worker started...")

	for item := range s.orderQueue {
		order := item.Order
		quantity := item.Quantity

		if err := s.repo.Create(order); err != nil {
			log.Printf("Failed to insert order: %v", err)
			continue
		}

		payload := map[string]interface{}{
			"orderId":    order.ID,
			"productId":  order.ProductID,
			"quantity":   quantity, // ✅ now available
			"totalPrice": order.TotalPrice,
			"createdAt":  order.CreatedAt,
		}

		body, _ := json.Marshal(payload)
		if err := s.rmq.Publish("order.exchange", "order.created", body); err != nil {
			log.Printf("Failed to publish order event: %v", err)
			continue
		}

		log.Printf("Processed order async: %+v", order)
	}
}

func (s *OrderService) GetOrdersByProductID(productID string) ([]models.Order, error) {
	cacheKey := fmt.Sprintf("orders:%s", productID)

	var orders []models.Order
	found, err := s.cache.Get(cacheKey, &orders)
	if err != nil {
		log.Printf("Redis get error: %v", err)
	}

	if found {
		log.Printf("Cache hit for product %s", productID)
		return orders, nil
	}
	log.Printf("Cache miss for product %s, querying DB", productID)
	orders, err = s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(cacheKey, orders, 30*time.Second)
	return orders, nil
}

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
