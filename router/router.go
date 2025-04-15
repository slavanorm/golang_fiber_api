package router

import (
	"rentincome/handler"
	"rentincome/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "${method} ${path} ${status} ${error}\n"}),
	)
	app.Use(cors.New())
	/*
	   // TODO: Configure CORS middleware
	   app.Use(cors.New(cors.Config{
	       AllowOrigins:     "http://localhost:3000, https://yourfrontend.com", // Allowed origins
	       AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS", // Allowed HTTP methods
	       AllowHeaders:     "Origin, Content-Type, Accept, Authorization",    // Allowed headers
	       AllowCredentials: true, // Allow cookies/credentials
	   }))
	*/
	app.Get("/health", handler.Health)

	auth := app.Group("/auth")

    auth.Post("/login", handler.Login)
	auth.Post("/register", handler.CreateUser)

	protected := app.Group("/api", middleware.JwtProtected())

    user := protected.Group("/user")
   
    user.Get("/:id", handler.GetUser)
    user.Put("/:id", handler.UpdateUser)
    user.Delete("/:id", handler.DeleteUser)

	estate := protected.Group("/estate")

    estate.Post("", handler.CreateEstate)
    estate.Get("/:id", handler.GetEstate)
	estate.Put("/:id", handler.UpdateEstate)
    estate.Delete("/:id", handler.DeleteEstate)
    
    transaction := protected.Group("/transaction")
    
    transaction.Post("", handler.CreateTransaction)
    transaction.Put("/:id", handler.UpdateTransaction)
    transaction.Delete("/:id", handler.DeleteTransaction)
    transaction.Get("/:id", handler.GetTransaction)

}
