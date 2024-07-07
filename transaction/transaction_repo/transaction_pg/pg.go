package transaction_pg

import (
	"context"
	"database/sql"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/transaction/transaction_repo"
	"log"

	"github.com/redis/go-redis/v9"
)

type transactionPg struct {
	db  *sql.DB
	rdb *redis.Client
}

const (
	addTransactionQuery = `insert into transaction (user_id, order_id) values($1, $2)`

	fetchUserIdQuery = `select id, user_id from transaction where id = $1`

	fetchAllCustomerTransactionQuery = `select t.id, o.product_id, p.name, o.qty, o.total_price, t.user_id, u.full_name, t.created_at, t.updated_at from transaction as t left join "user" as u on t.user_id = u.id left join "order" as o on t.order_id = o.id left join product as p on o.product_id = p.id where t.user_id = $1 and t.deleted_at is null order by created_at desc`

	fetchAllTransactionQuery = `select t.id, o.product_id, p.name, o.qty, o.total_price, t.user_id, u.full_name, t.created_at, t.updated_at from transaction as t left join "user" as u on t.user_id = u.id left join "order" as o on t.order_id = o.id left join product as p on o.product_id = p.id where t.deleted_at is null order by created_at desc`

	fetchTransactionByIdQuery = `select t.id, o.product_id, p.name, o.qty, o.total_price, t.user_id, u.full_name, t.created_at, t.updated_at from transaction as t left join "user" as u on t.user_id = u.id left join "order" as o on t.order_id = o.id left join product as p on o.product_id = p.id where t.id = $1 and t.deleted_at is null order by created_at desc`

	updateStockAndSoldQuery = `update product set stock = stock - (select qty from "order" where id = $1), sold = sold + (select qty from "order" where id = $1), updated_at = now() where id = (select product_id from "order" where id = $1)`
)

func NewTransactionPg(db *sql.DB, rdb *redis.Client) transaction_repo.TransactionRepo {
	return &transactionPg{
		db:  db,
		rdb: rdb,
	}
}

// Add implements transaction_repo.TransactionRepo.
func (pg *transactionPg) Add(transactionn *entity.Transaction) exception.Exception {

	tx, err := pg.db.Begin()

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(addTransactionQuery)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(transactionn.UserId, transactionn.OrderId); err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	updateStockAndSold, err := tx.Prepare(updateStockAndSoldQuery)

	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := updateStockAndSold.Exec(transactionn.OrderId); err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	if err := tx.Commit(); err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return exception.NewInternalServerError("something went wrong")
	}

	if err := pg.rdb.FlushDB(context.Background()).Err(); err != nil {
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	return nil
}

// CustomerTransaction implements transaction_repo.TransactionRepo.
func (pg *transactionPg) CustomerTransaction(id int) ([]*transaction_repo.TransactionWithProductsAndUserMapped, exception.Exception) {

	data := []*transaction_repo.TransactionWithProductsAndUserMapped{}

	stmt, err := pg.db.Prepare(fetchAllCustomerTransactionQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	rows, err := stmt.Query(id)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	for rows.Next() {

		transactionWithUserAndProduct := transaction_repo.TransactionWithProductsAndUserMapped{}

		if err := rows.Scan(
			&transactionWithUserAndProduct.Id,
			&transactionWithUserAndProduct.Product.Id,
			&transactionWithUserAndProduct.Product.Name,
			&transactionWithUserAndProduct.Qty,
			&transactionWithUserAndProduct.TotalPrice,
			&transactionWithUserAndProduct.User.Id,
			&transactionWithUserAndProduct.User.FullName,
			&transactionWithUserAndProduct.CreatedAt,
			&transactionWithUserAndProduct.UpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		data = append(data, &transactionWithUserAndProduct)
	}

	return data, nil
}

// FetchUserId implements transaction_repo.TransactionRepo.
func (pg *transactionPg) FetchUserId(id int) (*transaction_repo.TransactionWithProductsAndUserMapped, exception.Exception) {

	transactionWithUserAndProduct := transaction_repo.TransactionWithProductsAndUserMapped{}

	stmt, err := pg.db.Prepare(fetchUserIdQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	if err := stmt.QueryRow(id).Scan(
		&transactionWithUserAndProduct.Id,
		&transactionWithUserAndProduct.User.Id,
	); err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &transactionWithUserAndProduct, nil
}

// FetchTransactionById implements transaction_repo.TransactionRepo.
func (pg *transactionPg) FetchTransactionById(id int) (*transaction_repo.TransactionWithProductsAndUserMapped, exception.Exception) {

	transactionWithUserAndProduct := transaction_repo.TransactionWithProductsAndUserMapped{}

	stmt, err := pg.db.Prepare(fetchTransactionByIdQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	if err := stmt.QueryRow(id).Scan(
		&transactionWithUserAndProduct.Id,
		&transactionWithUserAndProduct.Product.Id,
		&transactionWithUserAndProduct.Product.Name,
		&transactionWithUserAndProduct.Qty,
		&transactionWithUserAndProduct.TotalPrice,
		&transactionWithUserAndProduct.User.Id,
		&transactionWithUserAndProduct.User.FullName,
		&transactionWithUserAndProduct.CreatedAt,
		&transactionWithUserAndProduct.UpdatedAt,
	); err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &transactionWithUserAndProduct, nil
}

// FetchTransactions implements transaction_repo.TransactionRepo.
func (pg *transactionPg) FetchTransactions() ([]*transaction_repo.TransactionWithProductsAndUserMapped, exception.Exception) {

	data := []*transaction_repo.TransactionWithProductsAndUserMapped{}

	rows, err := pg.db.Query(fetchAllTransactionQuery)

	if err != nil {
		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	for rows.Next() {

		transactionWithUserAndProduct := transaction_repo.TransactionWithProductsAndUserMapped{}

		if err := rows.Scan(
			&transactionWithUserAndProduct.Id,
			&transactionWithUserAndProduct.Product.Id,
			&transactionWithUserAndProduct.Product.Name,
			&transactionWithUserAndProduct.Qty,
			&transactionWithUserAndProduct.TotalPrice,
			&transactionWithUserAndProduct.User.Id,
			&transactionWithUserAndProduct.User.FullName,
			&transactionWithUserAndProduct.CreatedAt,
			&transactionWithUserAndProduct.UpdatedAt,
		); err != nil {
			log.Println(err.Error())
			return nil, exception.NewInternalServerError("something went wrong")
		}

		data = append(data, &transactionWithUserAndProduct)
	}

	return data, nil
}
