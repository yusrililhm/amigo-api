package order_repo

import (
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type OrderRepo interface {
	Add(order *entity.Order) exception.Exception
	Fetch(userId int) ([]*entity.OrderWithProduct, exception.Exception)
	Modify(order *entity.Order) exception.Exception
	Remove(id int) exception.Exception
	FetchOrderById(id int) (*entity.Order, exception.Exception)
}
