package order_handler

import (
	"encoding/json"
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/order/order_service"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"net/http"
	"strconv"
	"strings"
)

type orderHandler struct {
	os order_service.OrderService
}

type OrderHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	Modify(w http.ResponseWriter, r *http.Request)
	Remove(w http.ResponseWriter, r *http.Request)
}

func NewOrderHandler(os order_service.OrderService) OrderHandler {
	return &orderHandler{
		os: os,
	}
}

// Add implements OrderHandler.
func (oh *orderHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value("userData").(*entity.User)
	payload := &dto.AddOrderPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidJsonBodyRequest := exception.NewUnprocessableEntityError("invalid json body request")

		w.WriteHeader(invalidJsonBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidJsonBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := oh.os.Add(user.Id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Fetch implements OrderHandler.
func (oh *orderHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, _ := r.Context().Value("userData").(*entity.User)

	res, err := oh.os.Fetch(user.Id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Modify implements OrderHandler.
func (oh *orderHandler) Modify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(path[2])

	payload := &dto.ModifyOrderPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidJsonBodyRequest := exception.NewUnprocessableEntityError("invalid json body request")

		w.WriteHeader(invalidJsonBodyRequest.Status())
		w.Write(helper.ResponseJSON(invalidJsonBodyRequest))

		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := oh.os.Modify(id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Remove implements OrderHandler.
func (oh *orderHandler) Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(path[2])

	res, err := oh.os.Remove(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}
