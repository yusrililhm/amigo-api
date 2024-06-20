package user_service

import (
	"fashion-api/dto"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"net/http"
)

type serviceMock struct {
}

var (
	Authentication func(next http.Handler) http.Handler
	Authorization  func(next http.Handler) http.Handler
	ChangePassword func(id int, payload *dto.UserChangePasswordPayload) (*helper.ResponseBody, exception.Exception)
	Modify         func(id int, payload *dto.UserModifyPayload) (*helper.ResponseBody, exception.Exception)
	Profile        func(id int) (*helper.ResponseBody, exception.Exception)
	SignIn         func(payload *dto.UserSignInPayload) (*helper.ResponseBody, exception.Exception)
	SignUp         func(payload *dto.UserSignUpPayload) (*helper.ResponseBody, exception.Exception)
)

// Authentication implements UserService.
func (s *serviceMock) Authentication(next http.Handler) http.Handler {
	return Authentication(next)
}

// Authorization implements UserService.
func (s *serviceMock) Authorization(next http.Handler) http.Handler {
	return Authorization(next)
}

// ChangePassword implements UserService.
func (s *serviceMock) ChangePassword(id int, payload *dto.UserChangePasswordPayload) (*helper.ResponseBody, exception.Exception) {
	return ChangePassword(id, payload)
}

// Modify implements UserService.
func (s *serviceMock) Modify(id int, payload *dto.UserModifyPayload) (*helper.ResponseBody, exception.Exception) {
	return Modify(id, payload)
}

// Profile implements UserService.
func (s *serviceMock) Profile(id int) (*helper.ResponseBody, exception.Exception) {
	return Profile(id)
}

// SignIn implements UserService.
func (s *serviceMock) SignIn(payload *dto.UserSignInPayload) (*helper.ResponseBody, exception.Exception) {
	return SignIn(payload)
}

// SignUp implements UserService.
func (s *serviceMock) SignUp(payload *dto.UserSignUpPayload) (*helper.ResponseBody, exception.Exception) {
	return SignUp(payload)
}

func NewServiceMock() UserService {
	return &serviceMock{}
}
