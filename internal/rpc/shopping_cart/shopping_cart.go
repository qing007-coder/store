package shopping_cart

import (
	"context"
	"store/internal/proto/shopping_cart"
	"store/internal/rpc/base"
)

type ShoppingCart struct {
	*base.Base
}

func NewShoppingCart(b *base.Base) *ShoppingCart {
	return &ShoppingCart{
		b,
	}
}

func (s *ShoppingCart) AddShoppingCart(ctx context.Context, req *shopping_cart.AddShoppingCartReq, resp *shopping_cart.ShoppingCartResp) error {

}

func (s *ShoppingCart) RemoveShoppingCart(ctx context.Context, req *shopping_cart.RemoveShoppingCartReq, resp *shopping_cart.ShoppingCartResp) error {

}

func (s *ShoppingCart) GetShoppingCartList(ctx context.Context, req *shopping_cart.GetShoppingCartListReq, resp *shopping_cart.ShoppingCartResp) error {

}
