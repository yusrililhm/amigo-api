package user_service

import (
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/user/user_repo"
	"sync"

	"context"
	"net/http"
)

type userService struct {
	ur user_repo.UserRepository
	wg *sync.WaitGroup
}

type UserService interface {
	SignUp(payload *dto.UserSignUpPayload) (*helper.ResponseBody, exception.Exception)
	SignIn(payload *dto.UserSignInPayload) (*helper.ResponseBody, exception.Exception)
	Modify(id int, payload *dto.UserModifyPayload) (*helper.ResponseBody, exception.Exception)
	ChangePassword(id int, payload *dto.UserChangePasswordPayload) (*helper.ResponseBody, exception.Exception)
	Profile(id int) (*helper.ResponseBody, exception.Exception)
	Authentication(next http.Handler) http.Handler
	Authorization(next http.Handler) http.Handler
}

func NewUserService(ur user_repo.UserRepository, wg *sync.WaitGroup) UserService {
	return &userService{
		ur: ur,
		wg: wg,
	}
}

// Authorization implements UserService.
func (us *userService) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearerToken := r.Header.Get("Authorization")
		user := &entity.User{}

		if err := user.ValidateToken(bearerToken); err != nil {
			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))
			return
		}

		if user.Role != "admin" {
			err := exception.NewUnauthorizedError("you're not authorized to access this endpoint")

			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// Authentication implements UserService.
func (us *userService) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearerToken := r.Header.Get("Authorization")
		user := &entity.User{}

		if err := user.ValidateToken(bearerToken); err != nil {
			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))
			return
		}

		userData, err := us.ur.FetchByEmail(user.Email)

		if err != nil {
			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))
			return
		}

		ctx := context.WithValue(context.Background(), "userData", userData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// ChangePassword implements UserService.
func (us *userService) ChangePassword(id int, payload *dto.UserChangePasswordPayload) (*helper.ResponseBody, exception.Exception) {

	errCh := make(chan exception.Exception, 1)

	us.wg.Add(2)

	go func() {
		defer us.wg.Done()

		user, err := us.ur.FetchById(id)

		if err != nil {
			errCh <- err
			return
		}

		isValidPassword := user.CompareHashPassword(payload.OldPassword)

		if !isValidPassword {
			errCh <- exception.NewBadRequestError("invalid user")
			return
		}

		if payload.NewPassword != payload.ConfirmNewPassword {
			errCh <- exception.NewBadRequestError("password didn't match")
			return
		}
	}()

	u := &entity.User{
		Password: payload.NewPassword,
	}

	u.GenerateHashPassword()

	go func() {
		defer us.wg.Done()

		if err := us.ur.ChangePassword(id, u); err != nil {
			errCh <- err
			return
		}
	}()

	us.wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusOK,
			Message: "password successfully changed",
			Data:    nil,
		}, nil
	}
}

// Modify implements UserService.
func (us *userService) Modify(id int, payload *dto.UserModifyPayload) (*helper.ResponseBody, exception.Exception) {

	errCh := make(chan exception.Exception, 1)

	user := &entity.User{
		FullName: payload.FullName,
		Email:    payload.Email,
	}

	us.wg.Add(1)

	go func() {

		defer us.wg.Done()

		if err := us.ur.Modify(id, user); err != nil {
			errCh <- err
			return
		}
	}()

	us.wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusOK,
			Message: "user successfully modified",
			Data:    nil,
		}, nil
	}
}

// Profile implements UserService.
func (us *userService) Profile(id int) (*helper.ResponseBody, exception.Exception) {

	userCh, errCh := make(chan *entity.User, 1), make(chan exception.Exception, 1)

	us.wg.Add(1)

	go func() {
		defer us.wg.Done()

		user, err := us.ur.FetchById(id)

		if err != nil {
			errCh <- err
			return
		}

		userCh <- user
	}()

	user := <-userCh

	us.wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusOK,
			Message: "user successfully fetched",
			Data: &dto.UserData{
				Id:        user.Id,
				FullName:  user.FullName,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		}, nil
	}
}

// SignIn implements UserService.
func (us *userService) SignIn(payload *dto.UserSignInPayload) (*helper.ResponseBody, exception.Exception) {

	userCh, errCh := make(chan *entity.User, 1), make(chan exception.Exception, 1)

	us.wg.Add(1)

	go func() {
		defer us.wg.Done()

		user, err := us.ur.FetchByEmail(payload.Email)

		if err != nil {
			errCh <- err
			return
		}

		userCh <- user
	}()

	user := <-userCh

	isValidPassword := user.CompareHashPassword(payload.Password)

	if !isValidPassword {
		return nil, exception.NewBadRequestError("invalid email/password")
	}

	us.wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusOK,
			Message: "user successfully sign in",
			Data: &dto.TokenString{
				Token: user.GenereateTokenString(),
			},
		}, nil
	}
}

// SignUp implements UserService.
func (us *userService) SignUp(payload *dto.UserSignUpPayload) (*helper.ResponseBody, exception.Exception) {

	user := &entity.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: payload.Password,
	}

	user.GenerateHashPassword()

	chErr := make(chan exception.Exception, 1)
	defer close(chErr)

	us.wg.Add(1)

	go func() {
		defer us.wg.Done()

		if err := us.ur.Add(user); err != nil {
			chErr <- err
			return
		}
	}()

	us.wg.Wait()

	select {
	case err := <-chErr:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusCreated,
			Message: "user successfully sign up",
			Data:    nil,
		}, nil
	}
}
