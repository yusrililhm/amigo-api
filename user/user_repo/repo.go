package user_repo

import (
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type UserRepository interface {
	Add(user *entity.User) exception.Exception
	FetchById(id int) (*entity.User, exception.Exception)
	FetchByEmail(email string) (*entity.User, exception.Exception)
	Modify(id int, user *entity.User) exception.Exception
	ChangePassword(id int, user *entity.User) exception.Exception
}
