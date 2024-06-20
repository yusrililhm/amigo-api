package category_handler

import (
	"fashion-api/category/category_service"
	"fashion-api/dto"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"

	"encoding/json"
	"strconv"
	"strings"

	"net/http"

	"github.com/go-chi/chi/v5"
)

type categoryHandler struct {
	cs category_service.CategoryService
}

type CategoryHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	FetchById(w http.ResponseWriter, r *http.Request)
	Modify(w http.ResponseWriter, r *http.Request)
}

func NewCategoryHandler(cs category_service.CategoryService) CategoryHandler {
	return &categoryHandler{
		cs: cs,
	}
}

// Add implements CategoryHandler.
func (ch *categoryHandler) Add(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload := &dto.CategoryPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidJSON := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidJSON.Status())
		w.Write(helper.ResponseJSON(invalidJSON))
		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := ch.cs.Add(payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Delete implements CategoryHandler.
func (ch *categoryHandler) Delete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(path[2])

	res, err := ch.cs.Delete(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Fetch implements CategoryHandler.
func (ch *categoryHandler) Fetch(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res, err := ch.cs.Fetch()

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// FetchById implements CategoryHandler.
func (ch *categoryHandler) FetchById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(path[2])

	res, err := ch.cs.FetchById(id)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}

// Modify implements CategoryHandler.
func (ch *categoryHandler) Modify(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	chi.URLParam(r, "id")

	path := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(path[2])

	payload := &dto.CategoryPayload{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		invalidJSON := exception.NewUnprocessableEntityError("invalid JSON body request")

		w.WriteHeader(invalidJSON.Status())
		w.Write(helper.ResponseJSON(invalidJSON))
		return
	}

	if err := helper.ValidateStruct(payload); err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	res, err := ch.cs.Modify(id, payload)

	if err != nil {
		w.WriteHeader(err.Status())
		w.Write(helper.ResponseJSON(err))
		return
	}

	w.WriteHeader(res.Status)
	w.Write(helper.ResponseJSON(res))
}
