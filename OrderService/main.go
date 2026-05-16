package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"order_service/configs"
	"order_service/middleware"
	"order_service/services"
	"order_service/types"

	docs "order_service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type MainStore struct {
	DB *sql.DB
}

var placeService services.IPlaceService
var orderService services.IOrderService
var cartService services.ICartService
var userService services.IUserService
var productService services.IProductService

func getConnStr(dbHost, dbPort, dbUser, dbPassword, dbName string) string {
	//return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func initServices(ctx *gin.Context, store *MainStore) error {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return err
	}
	userService = services.UserService{AuthHeader: authHeader}
	placeService = services.PlaceService{DB: store.DB}
	orderService = services.OrderService{DB: store.DB}
	cartService = services.CartService{DB: store.DB}
	productService = services.ProductService{}

	return nil
}

func getAuthHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Auth header required")
	}
	return authHeader, nil
}

func main() {
	r := gin.Default()
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
		gr := v1.Group("/order-service")
		{
			gr.GET("/get-places", s.GetPlaces)
			gr.GET("/get-place/:id", s.GetPlace)

			gr.GET("/get-cart-orders/:id", s.GetCartOrders)
			//with auth
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.GET("/get-cart", s.GetCart)
			wa.GET("/get-bought-carts", s.GetBoughtCarts)
			wa.POST("/confirm-cart", s.ConfirmCart)
			wa.POST("/unconfirm-cart", s.UnconfirmCart)
			wa.POST("/mark-cart-as-bought", s.MarkCartAsBought)
			wa.POST("/mark-cart-as-not-bought", s.MarkCartAsNotBought)
			wa.DELETE("/delete-cart/:id", s.DeleteCart)

			wa.PUT("/add-order", s.AddOrder)
			wa.DELETE("/delete-order", s.DeleteOrder)
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
// @Summary GetPlaces
// @Schemes
// @Description description of function that get all places from DB
// @Tags order-service/place
// @Accept json
// @Produce json
// @Success 200 {string} s
// @Router /order-service/get-places [get]
func (s *MainStore) GetPlaces(ctx *gin.Context) {
	initServices(ctx, s)
	places, err := placeService.GetPlaces()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, places)
}

// @BasePath /api
// @Description description of function that get place by id
// @Tags order-service/place
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Place ID"
// @Success 200 {string} idk_WTF
// @Router /order-service/get-place/{id} [get]
func (s *MainStore) GetPlace(ctx *gin.Context) {
	initServices(ctx, s)
	id := ctx.Param("id")
	place, err := placeService.GetPlace(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, place)
}

// @BasePath /api
// @Description description of function that get cart to auth user
// @Tags order-service/cart
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /order-service/get-cart [get]
func (s *MainStore) GetCart(ctx *gin.Context) {
	initServices(ctx, s)
	userId, _ := ctx.Get("id") // protected data
	cart, err := cartService.GetCart(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

// @BasePath /api
// @Description description of function that get bought carts to auth user
// @Tags order-service/cart
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /order-service/get-bought-carts [get]
func (s *MainStore) GetBoughtCarts(ctx *gin.Context) {
	initServices(ctx, s)
	userId, _ := ctx.Get("id") // protected data
	carts, err := cartService.GetBoughtCarts(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, carts)
}

// @BasePath /api
// @Description confirm cart
// @Tags order-service/cart
// @Accept json
// @Param   placeId	query	string	false	"Place id"
// @Success 200 {string} Ok
// @Router /order-service/confirm-cart [POST]
func (s *MainStore) ConfirmCart(ctx *gin.Context) {
	initServices(ctx, s)
	userId, _ := ctx.Get("id") // protected data
	placeId := ctx.Query("placeId")
	cart, err := cartService.ConfirmCart(placeId, fmt.Sprint(userId), orderService, userService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

// @BasePath /api
// @Description Compensating actions for cart confirm
// @Tags order-service/cart
// @Accept json
// @Success 200 {string} Ok
// @Router /order-service/unconfirm-cart [POST]
func (s *MainStore) UnconfirmCart(ctx *gin.Context) {
	initServices(ctx, s)
	userId, _ := ctx.Get("id") // protected data
	cart, err := cartService.UnconfirmCart(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

// @BasePath /api
// @Description mark cart as bought
// @Tags order-service/cart
// @Accept json
// @Param   cartId	query	string	false	"Cart id"
// @Success 200 {string} Ok
// @Router /order-service/mark-cart-as-bought [POST]
func (s *MainStore) MarkCartAsBought(ctx *gin.Context) {
	initServices(ctx, s)
	cartId := ctx.Query("cartId")
	err := cartService.MarkCartAsBought(cartId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart successfuly bought")
}

// @BasePath /api
// @Description compensation actions for marking cart as bought
// @Tags order-service/cart
// @Accept json
// @Param   cartId	query	string	false	"Cart id"
// @Success 200 {string} Ok
// @Router /order-service/mark-cart-as-not-bought [POST]
func (s *MainStore) MarkCartAsNotBought(ctx *gin.Context) {
	initServices(ctx, s)
	cartId := ctx.Query("cartId")
	err := cartService.MarkCartAsNotBought(cartId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart successfuly unmark")
}

// @BasePath /api
// @Description delete cart
// @Tags order-service/cart
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Cart ID"
// @Success 200 {string} Ok
// @Router /order-service/delete-cart/{id} [DELETE]
func (s *MainStore) DeleteCart(ctx *gin.Context) {
	initServices(ctx, s)
	cartId := ctx.Param("id")
	err := cartService.DeleteCart(cartId, orderService, cartService, productService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart successfuly deleted")
}

// @BasePath /api
// @Summary GetOrders
// @Schemes
// @Description Get orders by cart id from db
// @Tags order-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Cart ID"
// @Success 200 {string} s
// @Router /order-service/get-cart-orders/{id} [get]
func (s *MainStore) GetCartOrders(ctx *gin.Context) {
	initServices(ctx, s)
	cartId := ctx.Param("id")
	orders, err := orderService.GetOrders(cartId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// @BasePath /api
// @Description add order
// @Tags order-service
// @Accept json
// @Produce json
// @Param order	body	types.OrderDto	true	"Order data"
// @Success 200 {string} Ok
// @Router /order-service/add-order [PUT]
func (s *MainStore) AddOrder(ctx *gin.Context) {
	initServices(ctx, s)
	orderDto := types.OrderDto{}
	if err := ctx.ShouldBindJSON(&orderDto); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	order, err := orderService.AddOrder(orderDto, cartService, productService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// @BasePath /api
// @Description delete order
// @Tags order-service
// @Accept json
// @Produce json
// @Param   productId	query	string		true	"Product ID"
// @Param   cartId		query	string		true	"Cart ID"
// @Success 200 {string} Ok
// @Router /order-service/delete-order [DELETE]
func (s *MainStore) DeleteOrder(ctx *gin.Context) {
	initServices(ctx, s)
	productId := ctx.Query("productId")
	cartId := ctx.Query("cartId")
	if productId == "" || cartId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product id and cart id params required !!!"})
		return
	}
	err := orderService.DeleteOrder(productId, cartId, cartService, productService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Order successfuly deleted")
}
