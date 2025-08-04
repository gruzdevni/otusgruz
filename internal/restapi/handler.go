package restapi

import (
	"otusgruz/internal/models"
	"otusgruz/internal/restapi/operations/other"

	"github.com/go-openapi/runtime/middleware"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}
