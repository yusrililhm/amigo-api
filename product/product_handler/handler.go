package product_handler

import (
	"fashion-api/dto"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/product/product_service"
	"strings"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type productHandler struct {
	ps product_service.ProductService
}

type ProductHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	FetchById(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Modify(w http.ResponseWriter, r *http.Request)
}

func NewProductHandler(ps product_service.ProductService) ProductHandler {
	return &productHandler{
		ps: ps,
	}
}

// Add implements ProductHandler.
func (ph *productHandler) Add(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload := &dto.ProductPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		err := exception.NewUnprocessableEntityError("invalid JSON body request")
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := ph.ps.Add(payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Delete implements ProductHandler.
func (ph *productHandler) Delete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(path[2])

	res, err := ph.ps.Delete(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Fetch implements ProductHandler.
func (ph *productHandler) Fetch(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res, err := ph.ps.Fetch()

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// FetchById implements ProductHandler.
func (ph *productHandler) FetchById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	res, err := ph.ps.FetchById(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Modify implements ProductHandler.
func (ph *productHandler) Modify(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(path[2])

	payload := &dto.ProductPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		err := exception.NewUnprocessableEntityError("invalid JSON body request")
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := ph.ps.Modify(id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}
