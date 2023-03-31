package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

var products []Product

func main() {

	Copy()
	CreateServer()

}

func Copy() {
	data, err := ioutil.ReadFile("/Users/sguerra/Documents/practica-goWeb/products.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &products)
	if err != nil {
		panic(err)
	}
}

func CreateServer() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/products", func(c *gin.Context) {
		c.IndentedJSON(200, products)
	})

	router.GET("/products/:id", getId)

	router.GET("/products/search/", getMaxPrice)

	router.Run(":8081")
}

func getMaxPrice(ctx *gin.Context) {
	price := ctx.Query("price")
	priceFloat, _ := strconv.ParseFloat(price, 64)
	var productsPrice []Product
	for _, p := range products {
		if p.Price > priceFloat {
			productsPrice = append(productsPrice, p)
		}
	}
	ctx.IndentedJSON(http.StatusOK, productsPrice)
}

func getProduct(ctx *gin.Context) {
	nombre := ctx.Query("nombre")
	for _, p := range products {
		if p.Name == nombre {
			ctx.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, nil)
}

func getId(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.Atoi(id)
	for _, p := range products {
		if p.Id == idInt {
			ctx.JSON(http.StatusOK, p)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, nil)
}
