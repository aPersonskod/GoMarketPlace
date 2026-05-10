package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "orchestrator_service/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
			gr.POST("/buy-cart", BuyCart)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8084")
}

type Cart struct{}

// @BasePath /api
// @Description buy cart
// @Tags buy-service
// @Accept json
// @Produce json
// @Param cart	body	Cart	true	"Cart data"
// @Success 200 {string} Ok
// @Router /buy-actions/buy-cart [POST]
func BuyCart(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Cart was bought!!!"})
}
