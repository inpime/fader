package payment

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/inpime/fader/utils/sdata"
)

func initTplContext() {
	pongo2.DefaultSet.Globals["PayViaBraintree"] = func(orderId, amount, opt *pongo2.Value) *pongo2.Value {
		orderOpt := OrderInfoFromM(
			orderId.String(),
			int64(amount.Float()*100),
			opt.Interface().(*sdata.StringMap),
		)

		txId, err := PayViaBraintreeGateway(orderOpt)

		logrus.WithFields(logrus.Fields{
			"order":  orderId.String(),
			"amount": amount.Integer(),
			"err":    err,
			"txid":   txId,
			"_api":   NAME,
		}).Infof("pay via braintree")

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(txId)
	}

	pongo2.DefaultSet.Globals["Payment"] = func() *pongo2.Value {
		return pongo2.AsValue(&Payment{})
	}

	pongo2.DefaultSet.Globals["BraintreeClientToken"] = func() *pongo2.Value {
		token, err := ClientTokenBraintree()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":  err,
				"_api": NAME,
			}).Infof("generate braintree client token")
			return pongo2.AsValue("")
		}

		return pongo2.AsValue(token)
	}

	pongo2.DefaultSet.Globals["PayFromNonceViaBraintreeGateway"] = func(payment_method_nonce, orderId, amount, opt *pongo2.Value) *pongo2.Value {
		orderOpt := OrderInfoFromM(
			orderId.String(),
			int64(amount.Float()*100),
			opt.Interface().(*sdata.StringMap),
		)

		txId, err := PayFromNonceViaBraintreeGateway(payment_method_nonce.String(), orderOpt)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"order":  orderId.String(),
				"amount": amount.Integer(),
				"err":    err,
				"txid":   txId,
				"_api":   NAME,
			}).Infof("pay via nonce")

			return pongo2.AsValue("")
		}

		return pongo2.AsValue(txId)
	}
}
