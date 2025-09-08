package build

import (
	internalclient "otusgruz/internal/client"
	authHTTP "otusgruz/internal/client/authhttp"
	"otusgruz/internal/client/billhttp"
)

func (b *Builder) NewAuthClient(doer internalclient.Doer) authHTTP.Client {
	client := authHTTP.NewClient(b.config.AuthInternal.HTTPAddress, doer)

	return client
}

func (b *Builder) NewBillClient(doer internalclient.Doer) billhttp.Client {
	client := billhttp.NewClient(b.config.BillInternal.HTTPAddress, doer)

	return client
}
