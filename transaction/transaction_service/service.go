package transaction_service

import (
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/order/order_repo"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/transaction/transaction_repo"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type transactionService struct {
	tr transaction_repo.TransactionRepo
	or order_repo.OrderRepo
}

type TransactionService interface {
	Authorization(next http.Handler) http.Handler
	Add(userId int, payload *dto.AddTransactionPayload) (*helper.ResponseBody, exception.Exception)
	FetchTransactionById(id int) (*helper.ResponseBody, exception.Exception)
	CustomersTransaction(userId int) (*helper.ResponseBody, exception.Exception)
	FetchAllTransaction() (*helper.ResponseBody, exception.Exception)
}

func NewTransactionService(tr transaction_repo.TransactionRepo, or order_repo.OrderRepo) TransactionService {
	return &transactionService{
		tr: tr,
		or: or,
	}
}

// Add implements TransactionService.
func (ts *transactionService) Add(userId int, payload *dto.AddTransactionPayload) (*helper.ResponseBody, exception.Exception) {

	order, err := ts.or.FetchOrderById(payload.OrderId)

	if err != nil {
		return nil, err
	}

	if order.UserId != userId {
		return nil, exception.NewUnauthorizedError("You're not authorized to access this order")
	}

	if err := ts.tr.Add(&entity.Transaction{
		UserId:  userId,
		OrderId: payload.OrderId,
	}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusCreated,
		Message: "transaction successfully added",
		Data:    nil,
	}, nil
}

// Authorization implements TransactionService.
func (ts *transactionService) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value("userData").(*entity.User)
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		transaction, err := ts.tr.FetchUserId(id)

		if err != nil {
			w.WriteHeader(err.Status())
			w.Write(helper.ResponseJSON(err))
			return
		}

		if user.Id != transaction.User.Id {
			unauthorized := exception.NewUnauthorizedError("you're not authorized to access this endpoint")

			w.WriteHeader(unauthorized.Status())
			w.Write(helper.ResponseJSON(unauthorized))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// CustomersTransaction implements TransactionService.
func (ts *transactionService) CustomersTransaction(userId int) (*helper.ResponseBody, exception.Exception) {

	data, err := ts.tr.CustomerTransaction(userId)

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "transaction successfully fetched",
		Data:    data,
	}, nil
}

// FetchAllTransaction implements TransactionService.
func (ts *transactionService) FetchAllTransaction() (*helper.ResponseBody, exception.Exception) {

	data, err := ts.tr.FetchTransactions()

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "transaction successfully fetched",
		Data:    data,
	}, nil
}

// FetchTransactionById implements TransactionService.
func (ts *transactionService) FetchTransactionById(id int) (*helper.ResponseBody, exception.Exception) {

	data, err := ts.tr.FetchTransactionById(id)

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "transaction successfully fetched",
		Data:    data,
	}, nil
}
