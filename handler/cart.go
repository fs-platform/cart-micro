package handler

import (
	"cart/domain/model"
	"cart/domain/service"
	cart "cart/proto/cart"
	"context"
	common "github.com/fs-platform/go-tool"
	log "github.com/micro/go-micro/v2/logger"
)

type Cart struct {
	CartDataService service.ICartDataService
}

func (c Cart) AddCart(ctx context.Context, in *cart.CartInfo, out *cart.ResponseAdd) error {
	cartModel := &model.Cart{}
	err := common.SwapTo(in, cartModel)
	if err != nil {
		return err
	}
	cartId, err := c.CartDataService.AddCart(cartModel)
	if err != nil {
		return err
	}
	out.CartId = cartId
	out.Msg = "新增成功"
	return nil
}

func (c Cart) ClearCart(ctx context.Context, in *cart.Clean, out *cart.Response) error {
	err := c.CartDataService.DeleteCart(in.UserId)
	if err != nil {
		return err
	}
	out.Msg = "清除成功"
	return nil
}

func (c Cart) Incr(ctx context.Context, in *cart.Item, out *cart.Response) error {
	err := c.CartDataService.IncrNum(in.Id, in.ChangeNum)
	if err != nil {
		return err
	}
	out.Msg = "新增成功"
	return nil
}

func (c Cart) Decr(ctx context.Context, in *cart.Item, out *cart.Response) error {
	err := c.CartDataService.DecrNum(in.Id, in.ChangeNum)
	if err != nil {
		return err
	}
	out.Msg = "减少成功"
	return nil
}

func (c Cart) DeleteItemById(ctx context.Context, in *cart.CartId, out *cart.Response) error {
	err := c.CartDataService.DeleteCart(in.Id)
	if err != nil {
		return err
	}
	out.Msg = "删除成功"
	return nil
}

func (c Cart) GetAll(ctx context.Context, in *cart.CartFindAll, out *cart.CartAll) error {
	data, err := c.CartDataService.FindAllCart(in.UserId)
	if err != nil {
		return err
	}
	for _, v := range data {
		cartInfo := &cart.CartInfo{}
		err := common.SwapTo(v, cartInfo)
		if err != nil {
			log.Error(err)
			break
		}
		out.CartInfo = append(out.CartInfo, cartInfo)
	}
	return nil
}
