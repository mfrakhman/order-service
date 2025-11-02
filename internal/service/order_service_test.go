package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder_Success(t *testing.T) {
	mockProduct := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
		Qty   int    `json:"qty"`
	}{
		ID: "mock123", Name: "Test Product", Price: 100000, Qty: 10,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockProduct)
	}))
	defer mockServer.Close()

	svc := &OrderService{
		productAPI: mockServer.URL,
		orderQueue: make(chan *orderQueueItem, 10),
	}

	order, err := svc.CreateOrder("mock123", 2)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "mock123", order.ProductID)
	assert.Equal(t, 200000, order.TotalPrice)
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "product not found", http.StatusNotFound)
	}))
	defer mockServer.Close()

	svc := &OrderService{
		productAPI: mockServer.URL,
		orderQueue: make(chan *orderQueueItem, 10),
	}

	order, err := svc.CreateOrder("unknown", 2)

	assert.Error(t, err)
	assert.Nil(t, order)
}
