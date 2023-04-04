package handlers

import (
	"errors"
	"net/http"
	"os"
	"practica-goWeb/estructura-api/internal/domain"
	"practica-goWeb/estructura-api/internal/product"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	s product.Service
}

// constructor de controller
func NewProductHandler(s product.Service) *productHandler {
	return &productHandler{
		s: s,
	}
}

func (h *productHandler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Product
		token := c.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			c.JSON(401, gin.H{"error": "token invalido"})
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		valid, err := Valid(&req)
		if !valid {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		product, err := h.s.Save(req)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "success", "data": product})
	}
}

func Valid(p *domain.Product) (bool, error) {
	if p.Name == "" || p.Code_value == "" || p.Quantity == 0 || p.Price == 0 || p.Expiration == "" {
		return false, errors.New("Hay uno o mas campos vacios")
	}
	_, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return false, errors.New("La fecha de expiracion debe ser ingresada en el formato:XX/XX/XXXX")
	}
	return true, nil
}

func (h *productHandler) MaxPrice() gin.HandlerFunc {
	return func(c *gin.Context) {
		price := c.Query("price")
		priceFloat, _ := strconv.ParseFloat(price, 64)
		products, err := h.s.GetMaxPrice(priceFloat)
		if err != nil {
			c.JSON(404, gin.H{"message": "no product found", "data": nil})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "success"})
		c.IndentedJSON(200, products)

	}
}

// func (ct *ControllerProduct) getProduct(ctx *gin.Context) {
// 	nombre := ctx.Query("nombre")
// 	for _, p := range products {
// 		if p.Name == nombre {
// 			ctx.IndentedJSON(http.StatusOK, p)
// 			return
// 		}
// 	}
// 	ctx.IndentedJSON(http.StatusNotFound, nil)
// }

func (h *productHandler) FindId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, _ := strconv.Atoi(id)
		product, err := h.s.GetId(idInt)
		if err != nil {
			c.JSON(404, gin.H{"message": "product not found", "data": nil})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "success"})
		c.IndentedJSON(200, product)

	}
}

func (h *productHandler) AllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, _ := h.s.GetAllProducts()
		c.IndentedJSON(200, products)
	}
}

// Delete elimina un producto
func (h *productHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "token invalido"})
			return
		}

		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "product deleted"})
	}
}

// Put actualiza un producto
func (h *productHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "token invalido"})
			return
		}

		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		var product domain.Product
		err = ctx.ShouldBindJSON(&product)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid product"})
			return
		}
		valid, err := Valid(&product)
		if !valid {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		p, err := h.s.Update(id, product)
		if err != nil {
			ctx.JSON(409, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
}

// Patch update selected fields of a product WIP
func (h *productHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Name        string  `json:"name,omitempty"`
		Quantity    int     `json:"quantity,omitempty"`
		CodeValue   string  `json:"code_value,omitempty"`
		IsPublished bool    `json:"is_published,omitempty"`
		Expiration  string  `json:"expiration,omitempty"`
		Price       float64 `json:"price,omitempty"`
	}
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "token invalido"})
			return
		}

		var r Request
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		if err := ctx.ShouldBindJSON(&r); err != nil {
			ctx.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		update := domain.Product{
			Name:         r.Name,
			Quantity:     r.Quantity,
			Code_value:   r.CodeValue,
			Is_published: r.IsPublished,
			Expiration:   r.Expiration,
			Price:        r.Price,
		}
		if update.Expiration != "" {
			valid, err := Valid(&update)
			if !valid {
				ctx.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
		p, err := h.s.Update(id, update)
		if err != nil {
			ctx.JSON(409, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
}
