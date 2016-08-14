package shop

import (
	"math"
)

// ExecuteCommand
func (o *Order) ExecuteCommand(cmd, productId string, count float64) *Order {
	switch cmd {
	case "incr":
		o.IncrProducts(productId)
	case "decr":
		o.DecrProducts(productId)
	case "set":
		o.SetProduct(productId, count)
	case "remove":
		o.RemoveProduct(productId)
	case "empty":
		o.Items = []*OrderItem{}
	}

	return o
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
