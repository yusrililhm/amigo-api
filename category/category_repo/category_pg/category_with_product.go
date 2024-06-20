package category_pg

import (
	"database/sql"
	"fashion-api/category/category_repo"
	"fashion-api/entity"
	"time"
)

type categoryWithProduct struct {
	categoryId        int
	categoryType      string
	productId         sql.NullInt64
	productName       sql.NullString
	productPrice      sql.NullInt64
	productStock      sql.NullInt64
	productSold       sql.NullInt64
	categoryCreatedAt time.Time
	categoryUpdatedAt time.Time
}

func (c *categoryWithProduct) categoryWithProductToAggregate() *category_repo.CategoryProduct {
	return &category_repo.CategoryProduct{
		Category: &entity.Category{
			Id:        c.categoryId,
			Type:      c.categoryType,
			CreatedAt: c.categoryCreatedAt,
			UpdatedAt: c.categoryUpdatedAt,
		},
		Product: &entity.Product{
			Id:    int(c.productId.Int64),
			Name:  c.productName.String,
			Price: int(c.productPrice.Int64),
			Stock: int(c.productStock.Int64),
			Sold:  int(c.productSold.Int64),
		},
	}
}
