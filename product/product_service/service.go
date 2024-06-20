package product_service

import (
	"fashion-api/category/category_repo"
	"fashion-api/dto"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/pkg/helper"
	"fashion-api/product/product_repo"
	"sync"

	"net/http"
)

type productService struct {
	pr product_repo.ProductRepo
	cr category_repo.CategoryRepo
	wg *sync.WaitGroup
}

type ProductService interface {
	Fetch() (*helper.ResponseBody, exception.Exception)
	FetchById(id int) (*helper.ResponseBody, exception.Exception)
	Add(payload *dto.ProductPayload) (*helper.ResponseBody, exception.Exception)
	Modify(id int, payload *dto.ProductPayload) (*helper.ResponseBody, exception.Exception)
	Delete(id int) (*helper.ResponseBody, exception.Exception)
}

func NewProductService(pr product_repo.ProductRepo, cr category_repo.CategoryRepo, wg *sync.WaitGroup) ProductService {
	return &productService{
		pr: pr,
		cr: cr,
		wg: wg,
	}
}

// Add implements ProductService.
func (ps *productService) Add(payload *dto.ProductPayload) (*helper.ResponseBody, exception.Exception) {

	errCh := make(chan exception.Exception, 1)

	_, err := ps.cr.FetchId(payload.CategoryId)

	if err != nil {
		return nil, err
	}
	
	ps.wg.Add(1)

	go func() {
		defer ps.wg.Done()

		if err := ps.pr.Add(&entity.Product{
			Name:        payload.Name,
			Description: payload.Description,
			CategoryId:  payload.CategoryId,
			Price:       payload.Price,
			Stock:       payload.Stock,
		}); err != nil {
			errCh <- err
			return
		}
	}()

	ps.wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return &helper.ResponseBody{
			Status:  http.StatusCreated,
			Message: "product successfully added",
			Data:    nil,
		}, nil
	}
}

// Delete implements ProductService.
func (ps *productService) Delete(id int) (*helper.ResponseBody, exception.Exception) {

	if _, err := ps.FetchById(id); err != nil {
		return nil, err
	}

	if err := ps.pr.Delete(id); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "product successfully deleted",
		Data:    nil,
	}, nil
}

// Fetch implements ProductService.
func (ps *productService) Fetch() (*helper.ResponseBody, exception.Exception) {

	products, err := ps.pr.Fetch()

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "products successfully fetched",
		Data:    products,
	}, nil
}

// FetchById implements ProductService.
func (ps *productService) FetchById(id int) (*helper.ResponseBody, exception.Exception) {

	product, err := ps.pr.FetchById(id)

	if err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "product with id successfully fetched",
		Data: &dto.ProductData{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			CategoryId:  product.CategoryId,
			Price:       product.Price,
			Stock:       product.Stock,
			Sold:        product.Sold,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	}, nil
}

// Modify implements ProductService.
func (ps *productService) Modify(id int, payload *dto.ProductPayload) (*helper.ResponseBody, exception.Exception) {

	_, err := ps.cr.FetchById(payload.CategoryId)

	if err != nil {
		return nil, err
	}

	if _, err := ps.FetchById(id); err != nil {
		return nil, err
	}

	if err := ps.pr.Modify(id, &entity.Product{
		Name:        payload.Name,
		Description: payload.Description,
		CategoryId:  payload.CategoryId,
		Price:       payload.Price,
		Stock:       payload.Stock,
	}); err != nil {
		return nil, err
	}

	return &helper.ResponseBody{
		Status:  http.StatusOK,
		Message: "product successfully modified",
		Data:    nil,
	}, nil
}
