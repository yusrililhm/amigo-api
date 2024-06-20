package category_repo

import (
	"fashion-api/entity"
	"time"
)

type CategoryProduct struct {
	Category *entity.Category
	Product  *entity.Product
}

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	Sold  int    `json:"sold"`
}

type CategoryWithProduct struct {
	Id        int        `json:"id"`
	Type      string     `json:"type"`
	Products  []*Product `json:"products"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (cwp *CategoryWithProduct) HandleCategoryWithProduct(categoriesProducts []*CategoryProduct) *CategoryWithProduct {

	if len(categoriesProducts) == 0 {
		return nil
	}

	categoryWithProducts := &CategoryWithProduct{}
	products := []*Product{}

	for _, eachCategoryProduct := range categoriesProducts {

		product := &Product{
			Id:    eachCategoryProduct.Product.Id,
			Name:  eachCategoryProduct.Product.Name,
			Price: eachCategoryProduct.Product.Price,
			Stock: eachCategoryProduct.Product.Stock,
			Sold:  eachCategoryProduct.Product.Sold,
		}

		products = append(products, product)

		categoryWithProducts = &CategoryWithProduct{
			Id:        eachCategoryProduct.Category.Id,
			Type:      eachCategoryProduct.Category.Type,
			Products:  products,
			CreatedAt: eachCategoryProduct.Category.CreatedAt,
			UpdatedAt: eachCategoryProduct.Category.UpdatedAt,
		}
	}

	if len(categoriesProducts) == 1 {
		categoryWithProducts.Products = []*Product{}
	}

	return categoryWithProducts
}
