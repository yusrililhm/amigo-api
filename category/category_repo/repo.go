package category_repo

import (
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type CategoryRepo interface {
	Add(category *entity.Category) exception.Exception
	Fetch() ([]*entity.Category, exception.Exception)
	FetchById(id int) (*CategoryWithProduct, exception.Exception)
	FetchId(id int) (*entity.Category, exception.Exception)
	Modify(id int, category *entity.Category) exception.Exception
	Delete(id int) exception.Exception
}
