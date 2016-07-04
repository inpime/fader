package api

import (
	// "api/config"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	// "github.com/inpime/dbox"
	// "net/url"
	// "store"
	// "strings"
	"utils"
	// "github.com/levigross/grequests"
)

func pongo2InitGlobalCustoms() {

	pongo2.DefaultSet.Globals["PayViaBraintreegateway"] = func(orderId, amount, opt *pongo2.Value) *pongo2.Value {
		orderOpt := OrderInfoFromM(
			orderId.String(),
			int64(amount.Integer()),
			opt.Interface().(utils.M),
		)

		txId, err := PayViaBraintreegateway(orderOpt)

		logrus.WithFields(logrus.Fields{
			"order":    orderId.String(),
			"amount":   amount.Integer(),
			"err":      err,
			"txid":     txId,
			"_service": "payment",
		}).Infof("pay via braintree")

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(txId)
	}
}
