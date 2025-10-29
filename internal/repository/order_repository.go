package repository

import (
	"ordersvc/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) FindByProductID(productID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.DB.Where("product_id = ?", productID).Find(&orders).Error
	return orders, err
}
