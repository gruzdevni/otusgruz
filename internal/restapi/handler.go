package restapi

import (
	"otusgruz/internal/models"
	"otusgruz/internal/restapi/operations/other"
	"otusgruz/internal/restapi/operations/user_c_r_u_d"
	"otusgruz/internal/service/api/user"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	userSrv user.Service
}

func NewHandler(userSrv user.Service) *Handler {
	return &Handler{
		userSrv: userSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) GetUser(params user_c_r_u_d.GetUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewGetUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.GetUser(ctx, userGUID)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewGetUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewGetUserGUIDOK().WithPayload(res)
}

func (h *Handler) CreateUser(params user_c_r_u_d.PostUserParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	res, err := h.userSrv.CreateUser(ctx, params.Request)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPostUserInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewPostUserOK().WithPayload(res)
}

func (h *Handler) UpdateUser(params user_c_r_u_d.PatchUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPatchUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.UpdateUser(ctx, userGUID, params.Request)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPatchUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewPatchUserGUIDOK().WithPayload(res)
}

func (h *Handler) DeleteUser(params user_c_r_u_d.DeleteUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewDeleteUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.DeleteUser(ctx, userGUID)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewDeleteUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewDeleteUserGUIDOK().WithPayload(res)
}
