package category_pg

import (
	"database/sql"
	"log"

	"fashion-api/category/category_repo"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type categoryPg struct {
	db *sql.DB
}

const (
	addCategoryQuery = `insert into "category" (type) values ($1)`

	deleteCategoryQuery = `update "category" set updated_at = now(), deleted_at = now() where id = $1`

	fetchCategoriesQuery = `select id, type, created_at, updated_at from "category" where deleted_at is null`

	fetchCategoryByIdQuery = `select c.id, c.type, p.id, p.name, p.price, p.stock, p.sold, c.created_at, c.updated_at from "category" as c left join product as p on c.id = p.category_id where c.id = $1 and c.deleted_at is null`

	fetchIdCategoryQuery = `select id from category where id = $1`

	modifyCategoryQuery = `update "category" set type = $2, updated_at = now() where id = $1`
)

func NewCategoryPg(db *sql.DB) category_repo.CategoryRepo {
	return &categoryPg{
		db: db,
	}
}

// Add implements category_repo.CategoryRepo.
func (pg *categoryPg) Add(category *entity.Category) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(addCategoryQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(category.Type); err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "category_type_key"` {
			tx.Rollback()
			log.Println(err.Error())
			return exception.NewConflictError("type has been created")
		}

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

// Delete implements category_repo.CategoryRepo.
func (pg *categoryPg) Delete(id int) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(deleteCategoryQuery)

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

// Fetch implements category_repo.CategoryRepo.
func (pg *categoryPg) Fetch() ([]*entity.Category, exception.Exception) {

	categories := []*entity.Category{}

	rows, err := pg.db.Query(fetchCategoriesQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	for rows.Next() {
		category := entity.Category{}

		if err := rows.Scan(
			&category.Id,
			&category.Type,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// FetchById implements category_repo.CategoryRepo.
func (pg *categoryPg) FetchById(id int) (*category_repo.CategoryWithProduct, exception.Exception) {

	categoriesProducts := []*category_repo.CategoryProduct{}
	categoryWithProducts := &category_repo.CategoryWithProduct{}

	stmt, err := pg.db.Prepare(fetchCategoryByIdQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	rows, err := stmt.Query(id)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	for rows.Next() {

		categoryProduct := categoryWithProduct{}

		if err := rows.Scan(
			&categoryProduct.categoryId,
			&categoryProduct.categoryType,
			&categoryProduct.productId,
			&categoryProduct.productName,
			&categoryProduct.productPrice,
			&categoryProduct.productStock,
			&categoryProduct.productSold,
			&categoryProduct.categoryCreatedAt,
			&categoryProduct.categoryUpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		categoriesProducts = append(categoriesProducts, categoryProduct.categoryWithProductToAggregate())
	}

	return categoryWithProducts.HandleCategoryWithProduct(categoriesProducts), nil
}

// FetchId implements category_repo.CategoryRepo.
func (pg *categoryPg) FetchId(id int) (*entity.Category, exception.Exception) {
	category := entity.Category{}

	stmt, err := pg.db.Prepare(fetchIdCategoryQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	if err := stmt.QueryRow(id).Scan(
		&category.Id,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, exception.NewNotFoundError("category not found")
		}

		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &category, nil
}

// Modify implements category_repo.CategoryRepo.
func (pg *categoryPg) Modify(id int, category *entity.Category) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(modifyCategoryQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(id, category.Type); err != nil {
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
