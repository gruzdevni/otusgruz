package restapi

import (
	"errors"

	"otusgruz/internal/apperr"
	"otusgruz/internal/models"
	"otusgruz/internal/restapi/operations/other"
	"otusgruz/internal/restapi/operations/user_c_r_u_d"
	"otusgruz/internal/service/api/auth"
	"otusgruz/internal/service/api/user"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type Handler struct {
	userSrv user.Service
	authSrv auth.Service
}

func NewHandler(userSrv user.Service, authSrv auth.Service) *Handler {
	return &Handler{
		userSrv: userSrv,
		authSrv: authSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetPublicHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) Login(params other.PostPublicLoginParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	loginParams := params.Request

	authorizedGUID, err := h.authSrv.Login(ctx, loginParams.Email.String(), loginParams.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrNotCorrectData) {
			return other.NewPostPublicLoginUnauthorized().WithPayload(&models.DefaultStatusResponse{Message: apperr.ErrNotCorrectData.Error()})
		}

		if errors.Is(err, apperr.ErrNoSuchUser) {
			return other.NewPostPublicLoginNotFound().WithPayload(&models.DefaultStatusResponse{Message: apperr.ErrNoSuchUser.Error()})
		}

		return other.NewPostPublicLoginInternalServerError()
	}

	return other.NewPostPublicLoginOK().WithXUser(strfmt.UUID(authorizedGUID.String()))
}

func (h *Handler) Signup(params other.PostPublicSignupParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	regParams := params.Request

	res, err := h.authSrv.Singup(ctx, *regParams)
	if err != nil {
		if errors.Is(err, apperr.ErrEmailAlreadyUsed) {
			return other.NewPostPublicSignupOK().WithPayload(&models.DefaultStatusResponse{Message: apperr.ErrEmailAlreadyUsed.Error()})
		}

		return other.NewPostPublicSignupInternalServerError().WithPayload(&models.DefaultStatusResponse{Message: err.Error()})
	}

	return other.NewPostPublicSignupOK().WithPayload(res)
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
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewGetUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

		errText = err.Error()
		return user_c_r_u_d.NewGetUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewGetUserGUIDOK().WithPayload(res)
}

func (h *Handler) UpdateUser(params user_c_r_u_d.PatchUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewPatchUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

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
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewDeleteUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

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
