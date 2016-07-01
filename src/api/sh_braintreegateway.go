package api

import (
	"api/config"
	braintree "github.com/lionelbarrow/braintree-go"
	"utils"
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

func OrderInfoFromM(orderId string, amount int64, opt utils.M) OrderInfo {
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

func PayViaBraintreegateway(opt OrderInfo) (string, error) {
	var env = braintree.Sandbox

	if config.AppSettings().M("braintree").Bool("prod") {
		env = braintree.Production
	}

	bt := braintree.New(
		env,
		config.AppSettings().M("braintree").String("merchantId"),
		config.AppSettings().M("braintree").String("publicKey"),
		config.AppSettings().M("braintree").String("privateKey"),
	)

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
