package shop

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/inpime/fader/api/context"
	"github.com/inpime/fader/store"
)

func NewProductsProvider(bucketName string) *ProductsProvider {
	return &ProductsProvider{
		bucketProducts: bucketName,
	}
}

type ProductsProvider struct {
	bucketProducts string
}

func (g *ProductsProvider) GetInfo(productId string) (*ProductInfo, error) {
	file, err := store.LoadOrNewFileID(
		strings.ToLower(g.bucketProducts),
		productId)

	if err != nil {
		return nil, err
	}

	return &ProductInfo{
		ID:        file.ID(),
		Title:     file.MMapData().String("Title"),
		Available: file.MMapData().Bool("AvailableForPurchase"),
		UnitPrice: file.MMapData().Float("UnitPrice"),
	}, nil
}

func initTplContext() {

	pongo2.DefaultSet.Globals["ShopBasketClear"] = func(ctx *pongo2.Value) *pongo2.Value {
		_ctx := ctx.Interface().(*context.Context)
		_ctx.Session().Set(shopBasketSessionKey, nil)

		return pongo2.AsValue("")
	}

	pongo2.DefaultSet.Globals["ShopBasket"] = func(ctx, bucketName *pongo2.Value) *pongo2.Value {
		_ctx := ctx.Interface().(*context.Context)

		basket := _ctx.Session().Get(shopBasketSessionKey)
		isUpdateSession := basket == nil

		basket = OrderFrom(basket)

		if basket, ok := basket.(*Order); ok {
			basket.SetProductProvider(NewProductsProvider(bucketName.String()))
		} else {
			logrus.WithField("ref", NAME).Errorf("Not expected basket type %T", basket)
		}

		if isUpdateSession {
			_ctx.Session().Set(shopBasketSessionKey, basket)
		}

		return pongo2.AsValue(basket)
	}

	pongo2.DefaultSet.Globals["ShopBasketSave"] = func(ctx, basket *pongo2.Value) *pongo2.Value {
		_ctx := ctx.Interface().(*context.Context)
		_ctx.Session().Set(shopBasketSessionKey, OrderFrom(basket.Interface()))

		return basket
	}
}
