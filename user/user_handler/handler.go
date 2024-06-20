package user_handler

import (
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/user/user_service"

	"encoding/json"
	"net/http"
)

type userHandler struct {
	us user_service.UserService
}

type UserHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	Modify(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(us user_service.UserService) UserHandler {
	return &userHandler{
		us: us,
	}
}

// ChangePassword implements UserHandler.
func (uh *userHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value("userData").(*entity.User)
	payload := &dto.UserChangePasswordPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidBodyRequest := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := uh.us.ChangePassword(u.Id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Modify implements UserHandler.
func (uh *userHandler) Modify(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value("userData").(*entity.User)
	payload := &dto.UserModifyPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidBodyRequest := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := uh.us.Modify(u.Id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Profile implements UserHandler.
func (uh *userHandler) Profile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	u := r.Context().Value("userData").(*entity.User)

	res, err := uh.us.Profile(u.Id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// SignIn implements UserHandler.
func (uh *userHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload := &dto.UserSignInPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidBodyRequest := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := uh.us.SignIn(payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// SignUp implements UserHandler.
func (uh *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload := &dto.UserSignUpPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidBodyRequest := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := uh.us.SignUp(payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}
