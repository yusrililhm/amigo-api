package transaction_handler

import (
	"encoding/json"
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/transaction/transaction_service"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type transactionHandler struct {
	ts transaction_service.TransactionService
}

type TransactionHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	FetchTransactionById(w http.ResponseWriter, r *http.Request)
	CustomersTransaction(w http.ResponseWriter, r *http.Request)
	FetchAllTransaction(w http.ResponseWriter, r *http.Request)
}

func NewTransactionHandler(ts transaction_service.TransactionService) TransactionHandler {
	return &transactionHandler{
		ts: ts,
	}
}

// Add implements TransactionHandler.
func (th *transactionHandler) Add(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload := &dto.AddTransactionPayload{}
	user := r.Context().Value("userData").(*entity.User)

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidJSONRequest := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidJSONRequest.Status())
		w.Write(helper.ResponseJSON(invalidJSONRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := th.ts.Add(user.Id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// CustomersTransaction implements TransactionHandler.
func (th *transactionHandler) CustomersTransaction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value("userData").(*entity.User)

	res, err := th.ts.CustomersTransaction(user.Id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// FetchTransactionById implements TransactionHandler.
func (th *transactionHandler) FetchTransactionById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := th.ts.FetchTransactionById(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// FetchAllTransaction implements TransactionHandler.
func (th *transactionHandler) FetchAllTransaction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res, err := th.ts.FetchAllTransaction()

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}
