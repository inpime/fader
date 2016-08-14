package shop

import (
	"testing"
)

func newProduct(id string, unitPrice float64, available bool) *ProductInfo {
	return &ProductInfo{
		ID:        id,
		UnitPrice: unitPrice,
		Title:     id,
		Available: available,
	}
}

var products = map[string]*ProductInfo{
	"1": newProduct("1", float64(1.2), true),
	"2": newProduct("2", float64(3.2), true),
	"3": newProduct("3", float64(3.2), false),
	"4": newProduct("4", float64(2.2), true),
}

type manager struct {
	products map[string]*ProductInfo
}

func (m *manager) GetInfo(productId string) (*ProductInfo, error) {
	return m.products[productId], nil
}

func TestShop_OrderBasket(t *testing.T) {
	order := NewOrder()
	order.SetProductProvider(&manager{products})
	order.SetProduct("1", 3)
	order.SetProduct("2", 2)
	order.SetProduct("3", 4)
	order.SetProduct("4", 3)
	order.RemoveProduct("4")
	order.IncrProducts("2") // 2 -> 3
	order.DecrProducts("1") // 3 -> 2

	order.Recalculate()

	want := 1.2*2.0 + 3.2*3.0
	if order.TotalAmount() != want {
		t.Errorf("not correct the total amount, %v, %v", order.TotalAmount(), want)
	}
}
