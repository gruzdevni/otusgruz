package authhttp

import "github.com/gofrs/uuid"

type User struct {
	GUID       string `json:"guid"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
}

type UserInfoResponse struct {
	Status int  `json:"status"`
	Result User `json:"result"`
	Errors any  `json:"errors"`
}

type DefaultResponse struct {
	Status int `json:"status"`
	Errors any `json:"errors"`
}

type LoginResponse struct {
	Status   int       `json:"status"`
	UserGUID uuid.UUID `json:"user_guid"`
	Errors   any       `json:"errors"`
}
