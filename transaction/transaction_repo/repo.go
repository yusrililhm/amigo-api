package transaction_repo

import (
	"fashion-api/entity"
	"fashion-api/pkg/exception"
)

type TransactionRepo interface {
	Add(transactionn *entity.Transaction) exception.Exception
	CustomerTransaction(id int) ([]*TransactionWithProductsAndUserMapped, exception.Exception)
	FetchUserId(id int) (*TransactionWithProductsAndUserMapped, exception.Exception)
	FetchTransactions() ([]*TransactionWithProductsAndUserMapped, exception.Exception)
	FetchTransactionById(id int) (*TransactionWithProductsAndUserMapped, exception.Exception)
}
