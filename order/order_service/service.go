package order_service

import (
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/order/order_repo"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/product/product_repo"
	"strconv"
	"strings"

	"net/http"
)

type orderService struct {
	or order_repo.OrderRepo
	pr product_repo.ProductRepo
}

type OrderService interface {
	Add(userId int, payload *dto.AddOrderPayload) (*helper.ResponseBody, exception.Exception)
	Fetch(userId int) (*helper.ResponseBody, exception.Exception)
	Modify(id int, payload *dto.ModifyOrderPayload) (*helper.ResponseBody, exception.Exception)
	Remove(id int) (*helper.ResponseBody, exception.Exception)
	Authorization(next http.Handler) http.Handler
}

func NewOrderService(or order_repo.OrderRepo, pr product_repo.ProductRepo) OrderService {
	return &orderService{
		or: or,
		pr: pr,
	}
}

// Add implements OrderService.
func (os *orderService) Add(userId int, payload *dto.AddOrderPayload) (*helper.ResponseBody, exception.Exception) {

	product, err := os.pr.FetchById(payload.ProductId)

	if err != nil {
		return nil, err
	}

	if payload.Qty > product.Stock {
		return nil, exception.NewBadRequestError("qty is greater than stock")
	}

	if err := os.or.Add(&entity.Order{
		UserId:    userId,
		ProductId: payload.ProductId,
		Qty:       payload.Qty,
	}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "order successfully added",
		Data:    nil,
	}, nil
}

// Authorization implements OrderService.
func (os *orderService) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value("userData").(*entity.User)

		path := strings.Split(r.URL.Path, "/")

		id, _ := strconv.Atoi(path[2])

		order, err := os.or.FetchOrderById(id)

		if err != nil {
			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))
			return
		}

		if user.Id != order.UserId {
			unauthorizationError := exception.NewUnauthorizedError("you're not authorized this order")

			w.WriteHeader(unauthorizationError.Status())
			w.Write(helper.ResponseJSON(unauthorizationError))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// Fetch implements OrderService.
func (os *orderService) Fetch(userId int) (*helper.ResponseBody, exception.Exception) {

	orders, err := os.or.Fetch(userId)

	if err != nil {
		return nil, err
	}

	data := []*entity.OrderWithProductMapped{}

	for _, eachOrder := range orders {
		orderWithProduct := &entity.OrderWithProductMapped{
			Id:     eachOrder.Id,
			UserId: eachOrder.UserId,
			Product: &entity.OrderProduct{
				Id:    eachOrder.ProductId,
				Name:  eachOrder.ProductName,
				Price: eachOrder.ProductPrice,
			},
			Qty:        eachOrder.Qty,
			TotalPrice: eachOrder.TotalPrice,
			CreatedAt:  eachOrder.CreatedAt,
			UpdatedAt:  eachOrder.UpdatedAt,
		}

		data = append(data, orderWithProduct)
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "order successfully fetched",
		Data:    data,
	}, nil
}

// Modify implements OrderService.
func (os *orderService) Modify(id int, payload *dto.ModifyOrderPayload) (*helper.ResponseBody, exception.Exception) {

	order, err := os.or.FetchOrderById(id)

	if err != nil {
		return nil, err
	}

	product, err := os.pr.FetchById(order.ProductId)

	if err != nil {
		return nil, err
	}

	if payload.Qty > product.Stock {
		return nil, exception.NewBadRequestError("qty is greater than stock")
	}

	if err := os.or.Modify(&entity.Order{
		Id:  id,
		Qty: payload.Qty,
	}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "order successfully modified",
		Data:    nil,
	}, nil
}

// Remove implements OrderService.
func (os *orderService) Remove(id int) (*helper.ResponseBody, exception.Exception) {

	if err := os.or.Remove(id); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "order successfully removed",
		Data:    nil,
	}, nil
}
