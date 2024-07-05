package user_pg

import (
	"database/sql"
	"fashion-api/entity"
	"fashion-api/pkg/exception"
	"fashion-api/user/user_repo"
	"log"
)

type userPg struct {
	db *sql.DB
}

const (
	addUserQuery = `insert into "user" (full_name, email, password, role) values($1, $2, $3, 'customers')`

	modifyUserQuery = `update "user" set full_name = $2, email = $3, address = $4, updated_at = now() where id = $1`

	changePasswordQuery = `update "user" set password = $2, updated_at = now() where id = $1`

	fetchUserByEmailQuery = `select id, full_name, email, password, role, address, created_at, updated_at from "user" where email = $1`

	fetchUserByIdQuery = `select id, full_name, email, password, role, address, created_at, updated_at from "user" where id = $1`
)

func NewUserPg(db *sql.DB) user_repo.UserRepository {
	return &userPg{
		db: db,
	}
}

// Add implements user_repo.UserRepository.
func (pg *userPg) Add(user *entity.User) exception.Exception {
	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(addUserQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(
		user.FullName,
		user.Email,
		user.Password,
	); err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"` {
			tx.Rollback()
			log.Println(err.Error())
			return exception.NewConflictError("email has been used")
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

// ChangePassword implements user_repo.UserRepository.
func (pg *userPg) ChangePassword(id int, user *entity.User) exception.Exception {
	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(changePasswordQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(
		id,
		user.Password,
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

// FetchByEmail implements user_repo.UserRepository.
func (pg *userPg) FetchByEmail(email string) (*entity.User, exception.Exception) {

	user := userData{}

	stmt, _ := pg.db.Prepare(fetchUserByEmailQuery)

	if err := stmt.QueryRow(email).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, exception.NewNotFoundError("user not found")
		}

		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &entity.User{
		Id:        user.Id,
		FullName:  user.FullName,
		Email:     email,
		Password:  user.Password,
		Role:      user.Role,
		Address:   user.Address.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// FetchById implements user_repo.UserRepository.
func (pg *userPg) FetchById(id int) (*entity.User, exception.Exception) {

	user := userData{}

	stmt, _ := pg.db.Prepare(fetchUserByIdQuery)

	if err := stmt.QueryRow(id).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err.Error())
			return nil, exception.NewNotFoundError("user not found")
		}

		log.Println(err.Error())
		return nil, exception.NewInternalServerError("something went wrong")
	}

	return &entity.User{
		Id:        user.Id,
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Address:   user.Address.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Modify implements user_repo.UserRepository.
func (pg *userPg) Modify(id int, user *entity.User) exception.Exception {
	tx, err := pg.db.Begin()

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	stmt, err := tx.Prepare(modifyUserQuery)

	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return exception.NewInternalServerError("something went wrong")
	}

	if _, err := stmt.Exec(
		id,
		user.FullName,
		user.Email,
		user.Address,
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
