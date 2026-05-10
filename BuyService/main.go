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

var buy_service services.IBuyService
var order_service services.IOrderService
var user_service services.IUserService

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

	user_service = services.UserService{AuthHeader: authHeader}
	order_service = services.OrderService{AuthHeader: authHeader}
	buy_service = services.BuyService{
		DB:           store.DB,
		OrderService: order_service,
		UserService:  user_service,
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
			wa.GET("/get-report-by-id/:id", s.GetReportById)
			wa.GET("/get-reports-by-userid", s.GetReportsByUserId)
			wa.POST("/buy-cart", s.BuyCart)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

// @BasePath /api
// @Description description of function that get report by id
// @Tags buy-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Report ID"
// @Success 200 {string} Ok
// @Router /buy-service/get-report-by-id/{id} [get]
func (s MainStore) GetReportById(ctx *gin.Context) {
	initServices(ctx, &s)
	reportId := ctx.Param("id")
	userId, _ := ctx.Get("id") // protected data

	reportDto, err := buy_service.GetReportById(reportId, fmt.Sprint(userId), user_service, order_service)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reportDto)
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

	reports, err := buy_service.GetReportsByUserId(fmt.Sprint(userId), user_service, order_service)
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

	err := buy_service.BuyCart(cart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart was successfuly bought")
}
