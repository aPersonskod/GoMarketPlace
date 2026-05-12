package main

import (
	"fmt"
	"net/http"
	"orchestrator_service/configs"
	"orchestrator_service/middleware"
	"orchestrator_service/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "orchestrator_service/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var userService services.IUserService
var orderService services.IOrderService
var buyService services.IBuyService

func initServices(ctx *gin.Context) error {
	authHeader, err := getAuthHeader(ctx)
	if err != nil {
		return err
	}
	userService = services.UserService{AuthHeader: authHeader}
	orderService = services.OrderService{AuthHeader: authHeader}
	buyService = services.BuyService{AuthHeader: authHeader}
	return nil
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
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api"
	v1 := r.Group("/api")
	{
		gr := v1.Group("/buy-actions")
		{
			wa := gr.Use(middleware.JwtAuthMiddleware())
			wa.POST("/buy-cart", BuyCart)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

// @BasePath /api
// @Description buy cart
// @Tags buy-service
// @Accept json
// @Produce json
// @Param   placeId	query	string	false	"Place id"
// @Success 200 {string} Ok
// @Router /buy-actions/buy-cart [POST]
func BuyCart(ctx *gin.Context) {
	// order-service.ConfirmCart
	// buy-service.BuyCart
	// user-service.SpendMoney
	// order-service.MarkCartAsBought
	// buy-service.AddBuyReport
	initServices(ctx)
	placeId := ctx.Query("placeId")

	confirmedCart, err := orderService.ConfirmCart(placeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = userService.SpendMoney(confirmedCart.AmountToPay)
	if err != nil {
		orderService.UnconfirmCart()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = orderService.MarkCartAsBought(confirmedCart.Id)
	if err != nil {
		userService.WalletReplenishment(confirmedCart.AmountToPay)
		orderService.UnconfirmCart()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	confirmedCart.IsBought = true
	err = buyService.BuyCart(*confirmedCart)
	if err != nil {
		orderService.MarkCartAsNotBought(confirmedCart.Id)
		userService.WalletReplenishment(confirmedCart.AmountToPay)
		orderService.UnconfirmCart()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart was successfully bought!!!"})
}
