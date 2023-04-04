package main

import (
	"encoding/json"
	"io/ioutil"
	"practica-goWeb/estructura-api/cmd/server/handlers"
	"practica-goWeb/estructura-api/internal/domain"
	"practica-goWeb/estructura-api/internal/product"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	productList := []domain.Product{}

	data, err := ioutil.ReadFile("/Users/sguerra/Documents/practica-goWeb/products.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &productList)
	if err != nil {
		panic(err)
	}

	_ = godotenv.Load()
	repo := product.NewRepository(productList)
	service := product.NewService(repo)
	productHandler := handlers.NewProductHandler(service)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	products := router.Group("/products")
	{
		products.GET("", productHandler.AllProducts())
		products.GET(":id", productHandler.FindId())
		products.GET("/search", productHandler.MaxPrice())
		products.POST("", productHandler.CreateProduct())
		products.DELETE(":id", productHandler.Delete())
		products.PATCH(":id", productHandler.Patch())
		products.PUT(":id", productHandler.Put())
	}

	router.Run(":8081")
}

//Listo, falta imports cuando este todo completo
