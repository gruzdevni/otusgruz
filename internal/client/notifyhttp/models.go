package notifyhttp

import (
	"github.com/google/uuid"
)

type User struct {
	GUID string `json:"guid"`
}

type CreateNotificationRequest struct {
	UserGUID     uuid.UUID `json:"user_guid"`
	OperationRef string    `json:"operation_ref"`
	NotifyType   string    `json:"notify_type"`
}

type DefaultResponse struct {
	Status int `json:"status"`
	Errors any `json:"errors"`
}
