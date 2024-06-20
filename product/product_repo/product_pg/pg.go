package product_pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/product/product_repo"

	"github.com/redis/go-redis/v9"
)

type productPg struct {
	db  *sql.DB
	rdb *redis.Client
}

const (
	addProductQuery = `insert into "product" (name, description, category_id, price, stock) values ($1, $2, $3, $4, $5)`

	fetchProductQuery = `select id, name, description, category_id, price, stock, sold, created_at, updated_at from "product" where deleted_at is null`

	fetchByIdProductQuery = `select id, name, description, category_id, price, stock, sold, created_at, updated_at from "product" where id = $1 and deleted_at is null`

	deleteProductQuery = `update "product" set updated_at = now(), deleted_at = now() where id = $1`

	modifyProductQuery = `update "product" set name = $2, description = $3, category_id = $4, price = $5, stock = $5, updated_at = now() where id = $1`
)

func NewProductPg(db *sql.DB, rdb *redis.Client) product_repo.ProductRepo {
	return &productPg{
		db:  db,
		rdb: rdb,
	}
}

// Add implements product_repo.ProductRepo.
func (pg *productPg) Add(product *entity.Product) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(addProductQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(
		product.Name,
		product.Description,
		product.CategoryId,
		product.Price,
		product.Stock,
	); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	return nil
}

// Delete implements product_repo.ProductRepo.
func (pg *productPg) Delete(id int) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(deleteProductQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(id); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	return nil
}

// Fetch implements product_repo.ProductRepo.
func (pg *productPg) Fetch() ([]*entity.Product, exception.Exception) {

	cmd := pg.rdb.Get(context.Background(), "products")

	products := []*entity.Product{}

	if cmd.Err() != redis.Nil {

		err := json.Unmarshal([]byte(cmd.Val()), &products)

		if err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		return products, nil
	}

	rows, err := pg.db.Query(fetchProductQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	for rows.Next() {

		product := entity.Product{}

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.CategoryId,
			&product.Price,
			&product.Stock,
			&product.Sold,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		products = append(products, &product)
	}

	if err := pg.setDataToRedis(products); err != nil {
		return nil, err
	}

	return products, nil
}

// FetchById implements product_repo.ProductRepo.
func (pg *productPg) FetchById(id int) (*entity.Product, exception.Exception) {

	product := entity.Product{}

	stmt, err := pg.db.Prepare(fetchByIdProductQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	if err := stmt.QueryRow(id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.CategoryId,
		&product.Price,
		&product.Stock,
		&product.Sold,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {

		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, exception.NewNotFoundError("product not found")
		}

		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &product, nil
}

// Modify implements product_repo.ProductRepo.
func (pg *productPg) Modify(id int, product *entity.Product) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(modifyProductQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(
		id,
		product.Name,
		product.Description,
		product.CategoryId,
		product.Price,
		product.Stock,
	); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	return nil
}

func (pg *productPg) setDataToRedis(products []*entity.Product) exception.Exception {

	data, _ := json.Marshal(products)

	if err := pg.rdb.Set(context.Background(), "products", data, 1*time.Hour).Err(); err != nil {
		log.Println(err.Error())
		return exception.NewInternalServerError("somethong went wrong")
	}

	return nil
}
