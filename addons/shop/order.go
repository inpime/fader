package shop

import (
	"encoding/gob"
	"log"
	"math"

	"github.com/inpime/sdata"
)

func init() {
	gob.Register(Order{})
}

func OrderFrom(order interface{}) *Order {
	if order, ok := order.(*Order); ok {
		return order
	}

	if order, ok := order.(Order); ok {
		return &order
	}

	return NewOrder()
}

func NewOrder() *Order {
	return &Order{
		Props: sdata.NewStringMap(),
	}
}

// Order order
type Order struct {
	Items []*OrderItem

	Props *sdata.StringMap

	manager ManagerProducts

	// DiscountCode string
	// Discount     float64

	InitialTotal float64 // sum(orders.InitialTotal)
	Total        float64 // TODO: (InitialTotal - Discount)
}

func (o *Order) SetProductProvider(manager ManagerProducts) *Order {
	o.manager = manager

	return o
}

func (o *Order) ValidBasket() bool {
	if len(o.Items) == 0 {
		return false
	}

	var countActiveItems = 0

	for _, item := range o.Items {
		if len(item.ProductID) == 0 ||
			item.Count == 0 ||
			item.Available == false ||
			item.Error != nil {

			return false
		}

		if !item.Removed {
			countActiveItems++
		}
	}

	if countActiveItems == 0 {
		return false
	}

	return true
}

// SetProduct add product
func (o *Order) SetProduct(productId string, count float64) *Order {
	count = math.Abs(count)

	if itemExists := o.GetProduct(productId); itemExists != nil {
		itemExists.Count = Round(count, 2)
		itemExists.Removed = false
	} else {
		o.Items = append(o.Items, NewOrderItem(productId, count))
	}

	return o
}

// IncrProducts increment count products
func (o *Order) IncrProducts(productId string) *Order {
	if itemExists := o.GetProduct(productId); itemExists != nil {
		itemExists.Count = Round(itemExists.Count+1, 2)
		itemExists.Removed = false
	} else {
		o.SetProduct(productId, 1)
	}

	return o
}

// IncrProducts decrement count products
func (o *Order) DecrProducts(productId string) *Order {
	if itemExists := o.GetProduct(productId); itemExists != nil {
		if (itemExists.Count - 1) > 0 {
			itemExists.Count = Round(itemExists.Count-1, 2)
			itemExists.Removed = false
		}
	}

	return o
}

// RemoveProduct remove product by productId
func (o *Order) RemoveProduct(productId string) *Order {
	if itemExists := o.GetProduct(productId); itemExists != nil {
		itemExists.Removed = true
		itemExists.Count = 0
	}

	return o
}

// GetProduct return product by productId or nil
func (o *Order) GetProduct(productId string) *OrderItem {
	for _, item := range o.Items {
		if item.ProductID == productId {
			return item
		}
	}

	return nil
}

// Recalculate check the availability of products,
// update the total amount of order
// update the total amount of order
func (o *Order) Recalculate() error {
	o.InitialTotal, o.Total = 0, 0

	for _, orderItem := range o.Items {
		if orderItem.Removed {
			continue
		}

		info, err := o.manager.GetInfo(orderItem.ProductID)

		if err != nil {
			log.Printf("refresh order: failed to get information about the product %q, %s", orderItem.ProductID, err)
			orderItem.Error = err
			continue
		}

		orderItem.Hydrate(info)
		orderItem.Recalculate()

		if orderItem.Available {
			o.InitialTotal = Round(o.InitialTotal+orderItem.InitialTotal, 2)
			o.Total = o.InitialTotal // TODO: discount
		}
	}

	return nil
}

func (o *Order) TotalAmount() float64 {

	return o.Total
}

func NewOrderItem(productId string, count float64) *OrderItem {
	return &OrderItem{
		ProductID: productId,
		Count:     count,
	}
}

// OrderItem product information
type OrderItem struct {
	ProductID string // reference to product
	Available bool   // is available to order product
	Error     error

	UnitPrice float64 // of product
	Count     float64

	Title string // of product
	// Description string

	// DiscountCode string
	// Discount     float64

	Removed bool // keep info of deleted products from your cart, for marketing

	InitialTotal float64 // UnitPrice*Count
	Total        float64 // TODO: (InitialTotal - Discount)
}

func (oi *OrderItem) Recalculate() *OrderItem {
	oi.InitialTotal = Round(oi.Count*oi.UnitPrice, 2)
	oi.Total = oi.InitialTotal // TODO: discount code

	return oi
}

func (oi *OrderItem) Hydrate(info *ProductInfo) *OrderItem {
	oi.Title = info.Title
	oi.UnitPrice = info.UnitPrice
	oi.Available = info.Available

	return oi
}

// Product info

type ProductInfo struct {
	ID        string
	UnitPrice float64
	Title     string
	Available bool
}
