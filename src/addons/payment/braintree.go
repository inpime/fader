package payment

import (
	"addons/standard"
	braintree "github.com/lionelbarrow/braintree-go"
	"strconv"
	"time"
	"utils/sdata"
)

/*
	OrderID

	Amount

	Number
	ExpDate
	CVV

	FirstName
	LastName
	Email

	BillingAddress
		StreetAddress
		Locality
		Region
		PostalCode
*/

func OrderInfoFromM(orderId string, amount int64, opt *sdata.StringMap) OrderInfo {
	return OrderInfo{
		OrderID: orderId,
		Amount:  amount,

		CardholderName: opt.String("CardholderName"),
		CardNumber:     opt.String("CardNumber"),
		ExpDate:        opt.String("ExpDate"),
		CVV:            opt.String("CVV"),

		Customer: OrderCustomer{
			FirstName: opt.String("FirstName"),
			LastName:  opt.String("LastName"),
			Email:     opt.String("Email"),
			Company:   opt.String("Company"),
			Phone:     opt.String("Phone"),
		},

		BillingAddress: TxAddress{
			StreetAddress: opt.String("StreetAddress"),
			Locality:      opt.String("Locality"),
			Region:        opt.String("Region"),
			PostalCode:    opt.String("PostalCode"),
		},
	}
}

type OrderInfo struct {
	OrderID string
	Amount  int64

	CardholderName string
	CardNumber     string
	ExpDate        string
	CVV            string

	Customer OrderCustomer

	BillingAddress  TxAddress
	ShippingAddress TxAddress
}

type OrderCustomer struct {
	FirstName string
	LastName  string
	Email     string
	Company   string
	Phone     string
}

type TxAddress struct {
	StreetAddress string
	Locality      string
	Region        string
	PostalCode    string
}

func newBraintree() *braintree.Braintree {
	var env = braintree.Sandbox

	paymentOption := standard.MainSettings().Config.M("braintree")

	if paymentOption.Bool("production") {
		env = braintree.Production
	}

	return braintree.New(
		env,
		paymentOption.String("merchantId"),
		paymentOption.String("publicKey"),
		paymentOption.String("privateKey"),
	)
}

func PayViaBraintreeGateway(opt OrderInfo) (string, error) {
	bt := newBraintree()

	// tx, err := bt.Transaction().Find()

	tx, err := bt.Transaction().Create(&braintree.Transaction{
		Type:    "sale",
		Amount:  braintree.NewDecimal(opt.Amount, 2), // 100 cents = 1 USD
		OrderId: opt.OrderID,
		CreditCard: &braintree.CreditCard{
			CardholderName: opt.CardholderName,
			Number:         opt.CardNumber,
			ExpirationDate: opt.ExpDate,
			CVV:            opt.CVV,
		},
		Customer: &braintree.Customer{
			FirstName: opt.Customer.FirstName,
			LastName:  opt.Customer.LastName,
			Email:     opt.Customer.Email,
			Company:   opt.Customer.Company,
		},
		BillingAddress: &braintree.Address{
			StreetAddress: opt.BillingAddress.StreetAddress,
			Locality:      opt.BillingAddress.Locality,
			Region:        opt.BillingAddress.Region,
			PostalCode:    opt.BillingAddress.PostalCode,
		},
		// ShippingAddress: &braintree.Address{
		// 	StreetAddress: opt.ShippingAddress.StreetAddress,
		// 	Locality:      opt.ShippingAddress.Locality,
		// 	Region:        opt.ShippingAddress.Region,
		// 	PostalCode:    opt.ShippingAddress.PostalCode,
		// },
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
			// 	// StoreInVault:                     true,
			AddBillingAddressToPaymentMethod: true,
			StoreShippingAddressInVault:      true,
		},
	})

	if err != nil {
		return "", err
	}

	return tx.Id, nil
}

func ClientTokenBraintree() (string, error) {
	bt := newBraintree()
	return bt.ClientToken().Generate()
}

func PayFromNonceViaBraintreeGateway(payment_method_nonce string) (string, error) {
	bt := newBraintree()

	orderId := strconv.FormatInt(time.Now().Unix(), 10)

	tx, err := bt.Transaction().Create(&braintree.Transaction{
		Type:               "sale",
		Amount:             braintree.NewDecimal(1234, 2), // 100 cents = 1 USD
		OrderId:            orderId,
		PaymentMethodNonce: payment_method_nonce,
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
		},
	})

	return tx.Id, err
}
