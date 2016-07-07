package payment

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"utils"
)

func initTplContext() {
	pongo2.DefaultSet.Globals["PayViaBraintree"] = func(orderId, amount, opt *pongo2.Value) *pongo2.Value {
		orderOpt := OrderInfoFromM(
			orderId.String(),
			int64(amount.Float()*100),
			opt.Interface().(utils.M),
		)

		txId, err := PayViaBraintreeGateway(orderOpt)

		logrus.WithFields(logrus.Fields{
			"order":    orderId.String(),
			"amount":   amount.Integer(),
			"err":      err,
			"txid":     txId,
			"_service": NAME,
		}).Infof("pay via braintree")

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(txId)
	}
}
