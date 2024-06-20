package app

import (
	"fashion-api/category/category_handler"
	"fashion-api/category/category_repo/category_pg"
	"fashion-api/category/category_service"
	"fashion-api/infra/config"
	"fashion-api/infra/db"
	"fashion-api/order/order_handler"
	"fashion-api/order/order_repo/order_pg"
	"fashion-api/order/order_service"
	"fashion-api/product/product_handler"
	"fashion-api/product/product_repo/product_pg"
	"fashion-api/product/product_service"
	"fashion-api/transaction/transaction_handler"
	"fashion-api/transaction/transaction_repo/transaction_pg"
	"fashion-api/transaction/transaction_service"
	"fashion-api/user/user_handler"
	"fashion-api/user/user_repo/user_pg"
	"fashion-api/user/user_service"
	"sync"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func StartApplication() {

	config.LoadEnv()

	db.InitializeDatabase()

	pg := db.NewPostgres()
	rdb := db.NewRedisClient()

	wg := &sync.WaitGroup{}

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// dependency injection
	ur := user_pg.NewUserPg(pg)
	us := user_service.NewUserService(ur, wg)
	uh := user_handler.NewUserHandler(us)

	cr := category_pg.NewCategoryPg(pg, rdb)
	cs := category_service.NewCategoryService(cr)
	ch := category_handler.NewCategoryHandler(cs)

	pr := product_pg.NewProductPg(pg, rdb)
	ps := product_service.NewProductService(pr, cr, wg)
	ph := product_handler.NewProductHandler(ps)

	or := order_pg.NewOrderPg(pg)
	os := order_service.NewOrderService(or, pr)
	oh := order_handler.NewOrderHandler(os)

	tr := transaction_pg.NewTransactionPg(pg)
	ts := transaction_service.NewTransactionService(tr, or)
	th := transaction_handler.NewTransactionHandler(ts)

	// user routes
	r.Group(func(r chi.Router) {
		r.Post("/user/signup", uh.SignUp)
		r.Post("/user/signin", uh.SignIn)

		r.Group(func(r chi.Router) {
			r.Use(us.Authentication)
			r.Get("/user", uh.Profile)
			r.Patch("/user", uh.Modify)
			r.Patch("/user/change-password", uh.ChangePassword)
		})
	})

	// product routes
	r.Group(func(r chi.Router) {
		r.Get("/products", ph.Fetch)
		r.Get("/products/{id}", ph.FetchById)

		r.Group(func(r chi.Router) {
			r.Use(us.Authentication, us.Authorization)
			r.Post("/products", ph.Add)
			r.Delete("/products/{id}", ph.Delete)
			r.Patch("/products/{id}", ph.Modify)
		})
	})

	// category routes
	r.Group(func(r chi.Router) {
		r.Get("/category", ch.Fetch)
		r.Get("/category/{id}", ch.FetchById)

		r.Group(func(r chi.Router) {
			r.Use(us.Authentication, us.Authorization)
			r.Post("/category", ch.Add)

			r.Patch("/category/{id}", ch.Modify)
			r.Delete("/category/{id}", ch.Delete)
		})
	})

	// order routes
	r.Group(func(r chi.Router) {
		r.Use(us.Authentication)
		r.Post("/orders", oh.Add)
		r.Get("/orders", oh.Fetch)

		r.Group(func(r chi.Router) {
			r.Use(us.Authentication, os.Authorization)
			r.Patch("/orders/{id}", oh.Modify)
			r.Delete("/orders/{id}", oh.Remove)
		})
	})

	// transaction routes
	r.Group(func(r chi.Router) {
		r.Use(us.Authentication)
		r.Post("/transaction", th.Add)
		r.Get("/transaction", th.CustomersTransaction)

		r.Group(func(r chi.Router) {
			r.Use(ts.Authorization)
			r.Get("/transaction/{id}", th.FetchTransactionById)
		})

		r.Group(func(r chi.Router) {
			r.Use(us.Authorization)
			r.Get("/admin/transaction", th.FetchAllTransaction)
		})
	})

	log.Println("[server] is running on port", config.NewAppConfig().AppPort)
	http.ListenAndServe(":"+config.NewAppConfig().AppPort, r)
}
