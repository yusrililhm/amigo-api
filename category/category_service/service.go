package category_service

import (
	"fashion-api/category/category_repo"
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"

	"net/http"
)

type categoryService struct {
	cr category_repo.CategoryRepo
}

type CategoryService interface {
	Add(payload *dto.CategoryPayload) (*helper.ResponseBody, exception.Exception)
	Fetch() (*helper.ResponseBody, exception.Exception)
	FetchById(id int) (*helper.ResponseBody, exception.Exception)
	Modify(id int, payload *dto.CategoryPayload) (*helper.ResponseBody, exception.Exception)
	Delete(id int) (*helper.ResponseBody, exception.Exception)
}

func NewCategoryService(cr category_repo.CategoryRepo) CategoryService {
	return &categoryService{
		cr: cr,
	}
}

// Add implements CategoryService.
func (cs *categoryService) Add(payload *dto.CategoryPayload) (*helper.ResponseBody, exception.Exception) {
	if err := cs.cr.Add(&entity.Category{
		Type: payload.Type,
	}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusCreated,
		Message: "category successfully added",
		Data:    nil,
	}, nil
}

// Delete implements CategoryService.
func (cs *categoryService) Delete(id int) (*helper.ResponseBody, exception.Exception) {

	_, err := cs.FetchById(id)

	if err != nil {
		return nil, err
	}

	if err := cs.cr.Delete(id); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "category successfully deleted",
		Data:    nil,
	}, nil
}

// Fetch implements CategoryService.
func (cs *categoryService) Fetch() (*helper.ResponseBody, exception.Exception) {

	categories, err := cs.cr.Fetch()

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "categories successfully fetched",
		Data:    categories,
	}, nil
}

// FetchById implements CategoryService.
func (cs *categoryService) FetchById(id int) (*helper.ResponseBody, exception.Exception) {

	catergory, err := cs.cr.FetchById(id)

	if err != nil {
		return nil, err
	}

	if catergory == nil {
		return nil, exception.NewNotFoundError("category not found")
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "category successfully fetched",
		Data:    catergory,
	}, nil
}

// Modify implements CategoryService.
func (cs *categoryService) Modify(id int, payload *dto.CategoryPayload) (*helper.ResponseBody, exception.Exception) {

	if err := cs.cr.Modify(id, &entity.Category{
		Type: payload.Type,
		}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "category successfully modified",
		Data:    nil,
	}, nil
}
