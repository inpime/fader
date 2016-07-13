package payment

import (
	"github.com/Sirupsen/logrus"
	"github.com/inpime/fader/utils/sdata"
)

type Payment struct {
	Braintree BraintreePayment
}

type BraintreePayment struct {
	Error error
}

func (b *BraintreePayment) ClientToken() string {
	token, err := ClientTokenBraintree()
	b.Error = err

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"_api": NAME,
		}).Infof("generate braintree client token")
		return ""
	}

	return token
}

func (b *BraintreePayment) Pay(orderId string, amount float64, opt *sdata.StringMap) string {
	orderOpt := OrderInfoFromM(
		orderId,
		int64(amount*100),
		opt,
	)

	txId, err := PayViaBraintreeGateway(orderOpt)
	b.Error = err

	logrus.WithFields(logrus.Fields{
		"order":  orderId,
		"amount": amount,
		"err":    err,
		"txid":   txId,
		"_api":   NAME,
	}).Infof("pay via braintree")

	return txId
}

func (b BraintreePayment) IsError() bool {
	return b.Error != nil
}
