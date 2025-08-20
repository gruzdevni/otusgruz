package authhttp

import "github.com/google/uuid"

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
	UserGUID uuid.UUID `json:"user_guid"`
}
