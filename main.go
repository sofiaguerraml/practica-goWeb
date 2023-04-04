package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	Code_value   string  `json:"code_value" binding:"required"`
	Is_published bool    `json:"is_published" binding:"required"`
	Expiration   string  `json:"expiration" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
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

	pr := router.Group("/products")
	{
		pr.GET("", func(c *gin.Context) {
			c.IndentedJSON(200, products)
		})
		pr.GET("/:id", getId)

		pr.GET("/search", getMaxPrice)

		pr.POST("", Save())
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

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
			ctx.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, nil)
}

func Save() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "123456" {
			c.JSON(401, gin.H{
				"error": "token invalido",
			})
			return
		}

		type request struct {
			Name         string  `json:"name"`
			Quantity     int     `json:"quantity"`
			Code_value   string  `json:"code_value"`
			Is_published bool    `json:"is_published"`
			Expiration   string  `json:"expiration"`
			Price        float64 `json:"price"`
		}

		var req request
		//recibe los datos y los pasa a nuestra estructura
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := Valid(req.Name, req.Quantity, req.Code_value, req.Is_published, req.Expiration, req.Price)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		pr := SaveProduct(req.Name, req.Quantity, req.Code_value, req.Is_published, req.Expiration, req.Price)
		c.JSON(200, pr)
	}
}

func exist(code string) bool {
	for _, c := range products {
		if c.Code_value == code {
			return true
		}
	}
	return false
}

func Valid(name string, quantity int, code_value string, is_published bool, expiration string, price float64) error {
	if name == "" || code_value == "" || quantity == 0 || price == 0 {
		return errors.New("Hay uno o mas campos vacios")
	}

	_, err := time.Parse("02/01/2006", expiration)
	if err != nil {
		return errors.New("La fecha de expiracion debe ser ingresada en el siguente formato:XX/XX/XXXX")
	}

	if exist(code_value) {
		return errors.New("Ya existe un producto con ese codigo")
	}
	return nil
}

func SaveProduct(name string, quantity int, code_value string, is_published bool, expiration string, price float64) Product {
	pr := Product{
		Id:           len(products) + 1,
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}
	products = append(products, pr)
	return pr
}
