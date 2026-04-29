package main

import (
	"fmt"
	"net/http"
	"order_service/configs"
	"order_service/middleware"
	"order_service/services"
	"strconv"

	docs "order_service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cartService services.ICartService
var orderService services.IOrderService
var placeService services.IPlaceService
var userService services.IUserService

func getConnStr(dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
}

func main() {
	/* 	userService = services.UserService{} // TODO fix services
	   	cartService = services.CartService{
	   		ConnStr: getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName),
	   		OrderService: nil,
	   		UserService: nil,
	   	} */
	placeService = services.PlaceService{ConnStr: getConnStr(configs.Env.DbUser, configs.Env.DbPassword, configs.Env.DbName)}
	r := gin.Default()

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

			gr.POST("/mark-cart-as-bought", MarkCartAsBought)

			gr.GET("/get-cart-orders", GetCartOrders)
			//with auth
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.GET("/get-cart", GetCart)
			wa.POST("/confirm-cart", ConfirmCart)

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
	id, _ := ctx.Get("id") // protected data
	cart, err := cartService.GetCart(fmt.Sprint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cart)
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
	userId, _ := ctx.Get("id") // protected data
	placeId := ctx.Query("placeId")
	cart, err := cartService.ConfirmCart(placeId, fmt.Sprint(userId))
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
// @Param   placeId	query	string	false	"Place id"
// @Success 200 {string} Ok
// @Router /order-service/mark-cart-as-bought [POST]
func MarkCartAsBought(ctx *gin.Context) {
	placeId := ctx.Query("placeId")
	err := cartService.MarkCartAsBought(placeId)
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
// @Param   productId	query	string	false	"Product id"
// @Param   cartId	query	string	false	"Cart id"
// @Param   quantity	query	int	false	"Quantity"
// @Success 200 {string} Ok
// @Router /order-service/add-order [PUT]
func AddOrder(ctx *gin.Context) {
	productId := ctx.Query("productId")
	cartId := ctx.Query("cartId")
	quantity, err := strconv.Atoi(ctx.Query("quantity"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := orderService.AddOrder(productId, cartId, quantity)
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
	productId := ctx.Param("id")
	err := orderService.DeleteOrder(productId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Order successfuly deleted")
}
