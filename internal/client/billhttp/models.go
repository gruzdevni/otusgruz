package billhttp

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	GUID string `json:"guid"`
}

type UserBalanceResponse struct {
	GUID   int             `json:"guid"`
	Amount decimal.Decimal `json:"amount"`
}

type ChangeBalanceRequest struct {
	UserGUID     uuid.UUID       `json:"user_guid"`
	OperationRef string          `json:"operation_ref"`
	Amount       decimal.Decimal `json:"amount"`
}

type DefaultResponse struct {
	Status int `json:"status"`
	Errors any `json:"errors"`
}
