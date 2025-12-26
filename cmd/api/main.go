package main

import (
	"fmt"
	"net/http"

	"example.com/api-GO/config"
	"example.com/api-GO/controllers"
	_ "example.com/api-GO/docs"
	"example.com/api-GO/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
	// Handler Swagger
	// Handler Swagger
)

// @title           Toko Golang API
// @version         1.0
// @description     Ini adalah server API Toko Online belajar Golang.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Support API
// @contact.email   support@tokogolang.com
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
//@BasePath /api
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
func main() {
	
	db := config.ConnectDB()
	defer db.Close()

	productController := controllers.ProductController{DB: db}
	authController := controllers.AuthController{DB: db}	

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/api/register", authController.Register)
	r.Post("/api/login", authController.Login)

	r.Route("/api/products", func(r chi.Router) {
        r.Get("/", productController.GetAll)

		r.Group(func (r chi.Router)  {

		r.Use(middleware.AuthMiddleware)
			r.Post("/", productController.Create)    
			r.Put("/{id}", productController.Update)
			r.Delete("/{id}", productController.Delete)
		})
        
    })
	fmt.Println("Server berjalan di http://localhost:8081")
	http.ListenAndServe(":8081", r)


}