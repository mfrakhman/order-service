package dto

type CreateOrderDTO struct {
	ProductID string `json:"productId" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}
