package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

var listadoProductos []Producto

const TOKEN string = "1A2B3C45D6"

type Producto struct {
	Id            int     `json:"id"`
	Nombre        string  `json:"nombre" binding:"required"`
	Color         string  `json:"color" binding:"required"`
	Precio        float64 `json:"precio" binding:"required"`
	Stock         int     `json:"stock" binding:"required"`
	Codigo        string  `json:"codigo" binding:"required"`
	Publicado     bool    `json:"publicado"` //TODO consultar
	FechaCreacion string  `json:"fechaCreacion" binding:"required"`
}

func HandlerGetAll(c *gin.Context) {
	/* Con productos hardcodeados:
	productos := []Producto{{1, "prod1", "rojo", 100.5, 200, "K4RF", true, "12-12-2021"}, {2, "prod2", "verde", 30.75, 100, "RK43", true, "13-12-2021"}}
	response, _ := json.Marshal(productos)
	c.JSON(200, string(response)) */

	productos, err := readProductsFromFile()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, productos)
}

func HandlerFilterProducts(c *gin.Context) {
	productos, err := readProductsFromFile()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var resultado []Producto
	filtros := c.Request.URL.Query()

	for _, producto := range productos {
		cumpleCondiciones := true

		for key, value := range filtros {
			switch key {
			case "id":
				castedValue, _ := strconv.Atoi(value[0])
				if producto.Id != castedValue {
					cumpleCondiciones = false
				}
			case "nombre":
				if producto.Nombre != value[0] {
					cumpleCondiciones = false
				}
			case "color":
				if producto.Color != value[0] {
					cumpleCondiciones = false
				}
			case "precio":
				castedValue, _ := strconv.ParseFloat(value[0], 32)
				if producto.Precio != castedValue {
					cumpleCondiciones = false
				}
			case "stock":
				castedValue, _ := strconv.Atoi(value[0])
				if producto.Stock != castedValue {
					cumpleCondiciones = false
				}
			case "codigo":
				if producto.Codigo != value[0] {
					cumpleCondiciones = false
				}
			case "publicado":
				castedValue, _ := strconv.ParseBool(value[0])
				if producto.Publicado == castedValue {
					cumpleCondiciones = false
				}
			case "fechaCreacion":
				if producto.FechaCreacion != value[0] {
					cumpleCondiciones = false
				}
			}
		}

		if cumpleCondiciones {
			resultado = append(resultado, producto)
		}
	}

	c.JSON(200, resultado)
}

func HandlerGetById(c *gin.Context) {
	productos, err := readProductsFromFile()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	idCasted, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "se produjo un error"})
		return
	}

	for _, prod := range productos {
		if prod.Id == idCasted {
			c.JSON(200, prod)
			return
		}
	}

	c.JSON(404, gin.H{"error": "no se encontró el producto con ese id"})
}

func HandlerCrearProducto(c *gin.Context) {
	var nuevoProducto Producto

	if err := c.ShouldBindJSON(&nuevoProducto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	nuevoProducto.Id = generateId()
	listadoProductos = append(listadoProductos, nuevoProducto)

	c.JSON(200, nuevoProducto)
}

func readProductsFromFile() ([]Producto, error) {
	var productos []Producto

	data, err := ioutil.ReadFile("./productos.json")
	if err != nil {
		return productos, err
	}

	err = json.Unmarshal(data, &productos)
	if err != nil {
		return productos, err
	}

	return productos, nil
}

func generateId() int {
	productosLength := len(listadoProductos)
	if productosLength == 0 {
		return 1
	}
	return listadoProductos[productosLength-1].Id + 1
}

func ValidateToken(c *gin.Context) {
	t := c.GetHeader("token")

	if t != TOKEN {
		c.AbortWithStatusJSON(401, gin.H{"error": "no tiene permisos para realizar la petición solicitada"})
	}
}

func main() {
	router := gin.Default()

	router.GET("/hola/:nombre", func(c *gin.Context) {
		saludo := "Hola " + c.Param("nombre")
		c.JSON(200, gin.H{"message": saludo})
	})

	productos := router.Group("/productos")
	{
		productos.GET("/", HandlerGetAll)
		productos.GET("/filter", HandlerFilterProducts)
		productos.GET("/:id", HandlerGetById)
		productos.POST("/new", ValidateToken, HandlerCrearProducto)
	}

	router.Run()
}
