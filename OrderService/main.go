package main

import (
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
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var placeService services.IPlaceService
var orderService services.IOrderService
var cartService services.ICartService
var userService services.IUserService

func getConnStr(dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
}

func initServices(ctx *gin.Context) error {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return err
	}
	userService = services.UserService{AuthHeader: authHeader}
	placeService = services.PlaceService{ConnStr: getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName)}
	orderService = services.OrderService{ConnStr: getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName)}
	cartService = services.CartService{
		ConnStr:      getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName),
		OrderService: orderService,
		UserService:  userService,
	}
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
		gr := v1.Group("/order-service")
		{
			gr.GET("/get-places", GetPlaces)
			gr.GET("/get-place/:id", GetPlace)

			gr.GET("/get-cart-orders/:id", GetCartOrders)
			//with auth
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.GET("/get-cart", GetCart)
			wa.GET("/get-bought-carts", GetBoughtCarts)
			wa.POST("/confirm-cart", ConfirmCart)
			wa.POST("/mark-cart-as-bought", MarkCartAsBought)

			wa.PUT("/add-order", AddOrder)
			wa.DELETE("/delete-order/:id", DeleteOrder)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

var connStr string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName)

// @BasePath /api
// @Summary GetPlaces
// @Schemes
// @Description description of function that get all places from DB
// @Tags order-service/place
// @Accept json
// @Produce json
// @Success 200 {string} s
// @Router /order-service/get-places [get]
func GetPlaces(ctx *gin.Context) {
	initServices(ctx)
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
func GetPlace(ctx *gin.Context) {
	initServices(ctx)
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
func GetCart(ctx *gin.Context) {
	initServices(ctx)
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
func GetBoughtCarts(ctx *gin.Context) {
	initServices(ctx)
	userId, _ := ctx.Get("id") // protected data
	carts, err := cartService.GetBoughtCarts(fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, carts)
}

/* // @BasePath /api
// @Description description of function that get cart by id
// @Tags order-service/cart
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Place ID"
// @Success 200 {string} idk_WTF
// @Router /order-service/get-cart/{id} [get]
func GetCartById(ctx *gin.Context) {
	id := ctx.Param("id")
	cart, err := cartService.(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
} */

// @BasePath /api
// @Description confirm cart
// @Tags order-service/cart
// @Accept json
// @Param   placeId	query	string	false	"Place id"
// @Success 200 {string} Ok
// @Router /order-service/confirm-cart [POST]
func ConfirmCart(ctx *gin.Context) {
	initServices(ctx)
	userId, _ := ctx.Get("id") // protected data
	placeId := ctx.Query("placeId")
	cart, err := cartService.ConfirmAndBuyCart(placeId, fmt.Sprint(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

// @BasePath /api
// @Description confirm cart
// @Tags order-service/cart
// @Accept json
// @Param   cartId	query	string	false	"Cart id"
// @Success 200 {string} Ok
// @Router /order-service/mark-cart-as-bought [POST]
func MarkCartAsBought(ctx *gin.Context) {
	initServices(ctx)
	cartId := ctx.Query("cartId")
	err := cartService.MarkCartAsBought(cartId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Cart successfuly bought")
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
func GetCartOrders(ctx *gin.Context) {
	initServices(ctx)
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
func AddOrder(ctx *gin.Context) {
	initServices(ctx)
	orderDto := types.OrderDto{}
	if err := ctx.ShouldBindJSON(&orderDto); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	order, err := orderService.AddOrder(orderDto)
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
// @Param   id	path	string		true	"Product ID"
// @Success 200 {string} idk_WTF
// @Router /order-service/delete-order/{id} [DELETE]
func DeleteOrder(ctx *gin.Context) {
	initServices(ctx)
	productId := ctx.Param("id")
	err := orderService.DeleteOrder(productId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Order successfuly deleted")
}
