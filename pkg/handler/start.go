package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sheikh-arman/go-project/pkg/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func Handle() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)
}
