package main

import (
	"fmt"
	"log"
	"net/http"
	"ordersvc/config"
	"ordersvc/internal/cache"
	"ordersvc/internal/db"
	"ordersvc/internal/handlers"
	"ordersvc/internal/rabbitmq"
	"ordersvc/internal/repository"
	"ordersvc/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg)
	rmq := rabbitmq.New(cfg.RabbitURL)
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	cache := cache.NewRedisCache(addr)
	defer rmq.Close()

	repo := repository.NewOrderRepository(database)
	svc := service.NewOrderService(repo, rmq, cfg.ProductAPI, cache)
	handler := handlers.NewOrderHandler(svc)

	go rmq.Consume("order.exchange", "order_log_queue", "order.created", func(payload map[string]interface{}) {
		log.Printf("Received order.created event: orderId=%v, productId=%v, quantity=%v, totalPrice=%v",
			payload["orderId"],
			payload["productId"],
			payload["quantity"],
			payload["totalPrice"],
		)
	})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.GetOrdersByProduct)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.ServicePort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Order-Service running on port %s", cfg.ServicePort)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}

}
