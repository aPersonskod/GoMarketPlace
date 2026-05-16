package main

import (
	"buy_service/configs"
	"buy_service/middleware"
	"buy_service/services"
	"buy_service/types"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	docs "buy_service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type MainStore struct {
	DB *sql.DB
}

func GetUserService(ctx *gin.Context) (*services.UserService, error) {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return nil, err
	}
	s := services.UserService{AuthHeader: authHeader}
	return &s, nil
}
func GetOrderService(ctx *gin.Context) (*services.OrderService, error) {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return nil, err
	}
	return &services.OrderService{AuthHeader: authHeader}, nil
}
func GetBuyService(db *sql.DB, ctx *gin.Context) (*services.BuyService, error) {
	orderService, err := GetOrderService(ctx)
	if err != nil {
		return nil, err
	}
	userService, err := GetUserService(ctx)
	if err != nil {
		return nil, err
	}
	return &services.BuyService{DB: db, OrderService: orderService, UserService: userService}, nil
}

func getConnStr(dbHost, dbPort, dbUser, dbPassword, dbName string) string {
	//return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func getAuthHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Auth header required")
	}
	return authHeader, nil
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

	db, err := sql.Open("postgres", getConnStr(configs.Env.DbHost, configs.Env.DbPort, configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName))
	if err != nil {
		panic(err.Error())
	}
	err = doMigration(db)
	if err != nil {
		panic(fmt.Sprintf("migration err: %s", err.Error()))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("db ping err: %s", err.Error()))
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

func doMigration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(fmt.Sprintf("Migration up failed: %v", err))
	}
	fmt.Println("Migration up completed successfully")
	return nil
}

// @BasePath /api
// @Description description of function that get reports by user's id
// @Tags buy-service
// @Accept json
// @Produce json
// @Success 200 {string} Ok
// @Router /buy-service/get-reports-by-userid [get]
func (s MainStore) GetReportsByUserId(ctx *gin.Context) {
	buyService, err := GetBuyService(s.DB, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	buyService, err := GetBuyService(s.DB, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cart := types.CartDto{}
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = buyService.BuyCart(cart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart was successfuly bought")
}
