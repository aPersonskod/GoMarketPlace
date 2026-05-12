package main

import (
	"buy_service/configs"
	"buy_service/middleware"
	"buy_service/services"
	"buy_service/types"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"

	docs "buy_service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type MainStore struct {
	DB *sql.DB
}

var buyService services.IBuyService
var orderService services.IOrderService
var userService services.IUserService

func getConnStr(dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
}

func getAuthHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Auth header required")
	}
	return authHeader, nil
}

func initServices(ctx *gin.Context, store *MainStore) error {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return err
	}

	userService = services.UserService{AuthHeader: authHeader}
	orderService = services.OrderService{AuthHeader: authHeader}
	buyService = services.BuyService{
		DB:           store.DB,
		OrderService: orderService,
		UserService:  userService,
	}

	return nil
}

func main() {
	r := gin.Default()
	// Use Default() for basic "allow all origins"
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	db, err := sql.Open("postgres", getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName))
	if err != nil {
		panic(err.Error())
	}
	s := MainStore{DB: db}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	docs.SwaggerInfo.BasePath = "/api"
	v1 := r.Group("/api")
	{
		gr := v1.Group("/buy-service")
		{
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.GET("/get-reports-by-userid", s.GetReportsByUserId)
			wa.POST("/buy-cart", s.BuyCart)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

// @BasePath /api
// @Description description of function that get reports by user's id
// @Tags buy-service
// @Accept json
// @Produce json
// @Success 200 {string} Ok
// @Router /buy-service/get-reports-by-userid [get]
func (s MainStore) GetReportsByUserId(ctx *gin.Context) {
	initServices(ctx, &s)
	userId, _ := ctx.Get("id") // protected data

	reports, err := buyService.GetReportsByUserId(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reports)
}

// @BasePath /api
// @Description buy cart
// @Tags buy-service
// @Accept json
// @Produce json
// @Param cart	body	types.Cart	true	"Cart data"
// @Success 200 {string} Ok
// @Router /buy-service/buy-cart [POST]
func (s MainStore) BuyCart(ctx *gin.Context) {
	initServices(ctx, &s)
	cart := types.CartDto{}
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := buyService.BuyCart(cart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart was successfuly bought")
}
