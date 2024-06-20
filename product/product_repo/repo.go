package product_repo

import (
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type ProductRepo interface {
	Fetch() ([]*entity.Product, exception.Exception)
	FetchById(id int) (*entity.Product, exception.Exception)
	Add(product *entity.Product) exception.Exception
	Modify(id int, product *entity.Product) exception.Exception
	Delete(id int) exception.Exception
}
