package main

import (
	"net/http"

	docs "product_service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	docs.SwaggerInfo.BasePath = "/api"
	v1 := r.Group("/api")
	{
		gr := v1.Group("/product-service")
		{
			gr.GET("/get-all", GetAll)
			gr.GET("/:id", GetById)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8081")
}

var connStr string = "user=postgres password=password dbname=marketplace-product-catalog-db sslmode=disable"

type Product struct {
	Id   string `gorm:"column:Id"`
	Name string `gorm:"column:Name"`
	Cost int    `gorm:"column:Cost"`
}

func (Product) TableName() string {
	return `public."Products"`
}

// @BasePath /api
// @Summary GetAll
// @Schemes
// @Description description of function that get all products from DB
// @Tags product-service
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /product-service/get-all [get]
func GetAll(ctx *gin.Context) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect database"})
		return
	}

	var products []Product
	db.Find(&products)
	ctx.JSON(http.StatusOK, products)

}

// @BasePath /api
// @Description description of function that get product by id
// @Tags product-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Some ID"
// @Success 200 {string} idk_WTF
// @Router /product-service/{id} [get]
func GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect database"})
		return
	}

	var product Product
	db.First(&product, `"Id" = ?`, id)
	ctx.JSON(http.StatusOK, product)
}
