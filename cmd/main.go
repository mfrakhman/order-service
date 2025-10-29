package main

import (
	"log"
	"ordersvc/config"
	"ordersvc/internal/db"
	"ordersvc/internal/handlers"
	"ordersvc/internal/rabbitmq"
	"ordersvc/internal/repository"
	"ordersvc/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg)
	rmq := rabbitmq.New(cfg.RabbitURL)
	defer rmq.Close()

	repo := repository.NewOrderRepository(database)
	svc := service.NewOrderService(repo, rmq, cfg.ProductAPI)
	handler := handlers.NewOrderHandler(svc)

	r := gin.Default()
	r.POST("/orders", handler.CreateOrder)
	r.GET("/orders", handler.GetOrdersByProduct)

	log.Printf("🚀 Order-Service running on port %s", cfg.ServicePort)
	r.Run(":" + cfg.ServicePort)
}
