package main

import (
	"fmt"
	"net/http"

	"product_service/configs"
	docs "product_service/docs"
	"product_service/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	// Use Default() for basic "allow all origins"
	r.Use(cors.Default())

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
			//with auth
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.PUT("/", AddProduct)
			wa.PATCH("/", UpdateProduct)
			wa.DELETE("/:id", DeleteProduct)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

var connStr string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName)

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

// @BasePath /api
// @Description add new product
// @Tags product-service
// @Accept json
// @Produce json
// @Param product	body	Product	true	"Product data"
// @Success 200 {string} idk_WTF
// @Router /product-service/ [PUT]
func AddProduct(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	if role != middleware.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have sufficient permissions"})
		return
	}

	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect database"})
		return
	}
	newId := fmt.Sprint(uuid.New())
	newProduct := &Product{
		Id:   newId,
		Name: product.Name,
		Cost: product.Cost,
	}
	err = gorm.G[Product](db).Create(ctx, newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, newProduct)
}

// @BasePath /api
// @Description update product
// @Tags product-service
// @Accept json
// @Produce json
// @Param product	body	Product	true	"Product data"
// @Success 200 {string} idk_WTF
// @Router /product-service/ [PATCH]
func UpdateProduct(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	if role != middleware.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have sufficient permissions"})
		return
	}

	var product Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect database"})
		return
	}
	_, err = gorm.G[Product](db).Where(`"Id" = ?`, product.Id).Updates(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// @BasePath /api
// @Description delete product
// @Tags product-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Some ID"
// @Success 200 {string} idk_WTF
// @Router /product-service/{id} [DELETE]
func DeleteProduct(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	if role != middleware.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have sufficient permissions"})
		return
	}

	id := ctx.Param("id")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect database"})
		return
	}
	rowsAffected, err := gorm.G[Product](db).Where(`"Id" = ?`, id).Delete(ctx)
	if rowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Succesfuly deleted")
}
