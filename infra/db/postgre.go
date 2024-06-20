package db

import (
	"database/sql"
	"fmt"
	"log"

	"fashion-api/entity"
	"fashion-api/infra/config"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func handleDatabaseConnection() {

	appConfig := config.NewAppConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBName,
	)

	db, err = sql.Open(appConfig.DBDialect, dsn)

	if err != nil {
		log.Fatal("error occured while trying to validate database arguments: ", err.Error())
		return
	}

	if err := db.Ping(); err != nil {
		log.Fatal("error occured while trying to connect to database: ", err.Error())
		return
	}
}

func handleRequiredTables() {
	const (
		createTableUserQuery = `create table if not exists "user" (
			id serial primary key,
			full_name varchar(60) not null,
			email varchar(60) not null unique,
			password text not null,
			role varchar not null,
			created_at timestamptz default now(),
			updated_at timestamptz default now(),
			deleted_at timestamptz
		)`

		createTableCategoryQuery = `create table if not exists "category" (
			id serial primary key,
			type varchar(60) not null unique,
			created_at timestamptz default now(),
			updated_at timestamptz default now(),
			deleted_at timestamptz
		)
		`

		createTableProductQuery = `create table if not exists "product" (
			id serial primary key,
			name varchar(60) not null,
			description text not null,
			category_id int not null,
			price int not null,
			stock int not null,
			sold int default 0,
			created_at timestamptz default now(),
			updated_at timestamptz default now(),
			deleted_at timestamptz,
			constraint fk_category_id foreign key (category_id) references category(id)
		)
		`

		createTableOrderQuery = `create table if not exists "order" (
			id serial primary key,
			user_id int not null,
			product_id int not null,
			qty int not null,
			total_price int not null,
			created_at timestamptz default now(),
			updated_at timestamptz default now(),
			deleted_at timestamptz,
			constraint fk_user_id foreign key (user_id) references "user"(id),
			constraint fk_product_id foreign key (product_id) references product(id)
		)
		`

		createTableTransactionQuery = `create table if not exists "transaction" (
			id serial primary key,
			user_id int not null,
			order_id int not null,
			created_at timestamptz default now(),
			updated_at timestamptz default now(),
			deleted_at timestamptz,
			constraint fk_user_id foreign key (user_id) references "user"(id),
			constraint fk_order_id foreign key (order_id) references "order"(id)
		)`

		createAdminQuery = `insert into "user" (full_name, email, password, role) values($1, $2, $3, $4) on conflict(email) do nothing`
	)

	if _, err := db.Exec(createTableUserQuery); err != nil {
		log.Fatal("error occured while create table user : ", err.Error())
		return
	}

	if _, err := db.Exec(createTableCategoryQuery); err != nil {
		log.Fatal("error occured while create table category : ", err.Error())
		return
	}

	if _, err := db.Exec(createTableProductQuery); err != nil {
		log.Fatal("error occured while create table product : ", err.Error())
		return
	}

	if _, err := db.Exec(createTableOrderQuery); err != nil {
		log.Fatal("error occured while create table order : ", err.Error())
		return
	}

	if _, err := db.Exec(createTableTransactionQuery); err != nil {
		log.Fatal("error occured while create table transaction : ", err.Error())
		return
	}

	u := &entity.User{
		FullName: config.NewAppConfig().AdminFullName,
		Email:    config.NewAppConfig().AdminEmail,
		Password: config.NewAppConfig().AdminPassword,
	}

	u.GenerateHashPassword()

	if _, err := db.Exec(createAdminQuery, u.FullName, u.Email, u.Password, "admin"); err != nil {
		log.Fatal("error occured while create admin : ", err.Error())
		return
	}
}

func InitializeDatabase() {
	handleDatabaseConnection()
	handleRequiredTables()
}

func NewPostgres() *sql.DB {
	return db
}
