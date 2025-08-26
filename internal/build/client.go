package build

import (
	internalclient "otusgruz/internal/client"
	authHTTP "otusgruz/internal/client/authhttp"
)

func (b *Builder) NewAuthClient(doer internalclient.Doer) authHTTP.Client {
	client := authHTTP.NewClient(b.config.AuthInternal.HTTPAddress, doer)

	return client
}
