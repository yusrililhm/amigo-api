package order_pg

import (
	"database/sql"
	"fashion-api/entity"
	"fashion-api/order/order_repo"
	"fashion-api/pkg/exception"
	"log"
)

type orderPg struct {
	db *sql.DB
}

func NewOrderPg(db *sql.DB) order_repo.OrderRepo {
	return &orderPg{
		db: db,
	}
}

const (
	addOrderQuery = `insert into "order" (user_id, product_id, qty, total_price) values ($1, $2, $3, ((select p.price from product as p where id = $2)) * $3);`

	fetchOrderQuery = `select o.id, o.user_id, o.product_id, p.name, p.price, o.qty, o.total_price, o.created_at, o.updated_at from "order" as o left join product as p on o.product_id = p.id where o.user_id = $1 and o.deleted_at is null;`

	fetchUserIdQuery = `select user_id, product_id from "order" where id = $1`

	modifyOrderQuery = `update "order" set qty = $2, total_price = ((select p.price from product as p where id = (select o.product_id from "order" as o where id = $1)) * $2), updated_at = now() where id = $1;`

	deleteOrderQuery = `update "order" set deleted_at = now(), updated_at = now() where id = $1;`
)

// Add implements order_repo.OrderRepo.
func (pg *orderPg) Add(order *entity.Order) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(addOrderQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(order.UserId, order.ProductId, order.Qty); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt.Close()

	return nil
}

// Fetch implements order_repo.OrderRepo.
func (pg *orderPg) Fetch(userId int) ([]*entity.OrderWithProduct, exception.Exception) {
	orders := []*entity.OrderWithProduct{}

	stmt, err := pg.db.Prepare(fetchOrderQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	rows, err := stmt.Query(userId)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	for rows.Next() {
		order := entity.OrderWithProduct{}

		if err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.ProductId,
			&order.ProductName,
			&order.ProductPrice,
			&order.Qty,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		orders = append(orders, &order)
	}

	stmt.Close()

	return orders, nil
}

// FetchOrderById implements order_repo.OrderRepo.
func (pg *orderPg) FetchOrderById(id int) (*entity.Order, exception.Exception) {

	order := entity.Order{}

	stmt, err := pg.db.Prepare(fetchUserIdQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	if err := stmt.QueryRow(id).Scan(
		&order.UserId,
		&order.ProductId,
	); err != nil {

		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, exception.NewNotFoundError("order not found")
		}

		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &order, nil
}

// Modify implements order_repo.OrderRepo.
func (pg *orderPg) Modify(order *entity.Order) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(modifyOrderQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(order.Id, order.Qty); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt.Close()

	return nil
}

// Remove implements order_repo.OrderRepo.
func (pg *orderPg) Remove(id int) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(deleteOrderQuery)

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

	stmt.Close()

	return nil
}
