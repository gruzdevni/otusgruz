package build

import (
	internalclient "otusgruz/internal/client"
	authHTTP "otusgruz/internal/client/authhttp"
	"otusgruz/internal/client/billhttp"
	"otusgruz/internal/client/deliveryhttp"
	"otusgruz/internal/client/goodshttp"
	"otusgruz/internal/client/notifyhttp"
)

func (b *Builder) NewAuthClient(doer internalclient.Doer) authHTTP.Client {
	client := authHTTP.NewClient(b.config.AuthInternal.HTTPAddress, doer)

	return client
}

func (b *Builder) NewBillClient(doer internalclient.Doer) billhttp.Client {
	client := billhttp.NewClient(b.config.BillInternal.HTTPAddress, doer)

	return client
}

func (b *Builder) NewNotifyClient(doer internalclient.Doer) notifyhttp.Client {
	client := notifyhttp.NewClient(b.config.NotifyInternal.HTTPAddress, doer)

	return client
}

func (b *Builder) NewDeliveryClient(doer internalclient.Doer) deliveryhttp.Client {
	client := deliveryhttp.NewClient(b.config.DeliveryInternal.HTTPAddress, doer)

	return client
}

func (b *Builder) NewGoodsClient(doer internalclient.Doer) goodshttp.Client {
	client := goodshttp.NewClient(b.config.GoodsInternal.HTTPAddress, doer)

	return client
}
