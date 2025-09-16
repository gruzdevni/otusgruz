package build

import (
	internalclient "otusgruz/internal/client"
	authHTTP "otusgruz/internal/client/authhttp"
	"otusgruz/internal/client/billhttp"
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
