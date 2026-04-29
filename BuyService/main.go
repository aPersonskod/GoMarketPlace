package main

import (
	"buy_service/configs"
	"fmt"
	"net/http"

	docs "buy_service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			/* 			gr.GET("/:id", GetById)
			   			//with auth
			   			wa := gr.Use(middleware.JwtAuthMiddleware())
			   			wa.PUT("/", AddProduct)
			   			wa.PATCH("/", UpdateProduct)
			   			wa.DELETE("/:id", DeleteProduct) */
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
}

func GetAll(ctx *gin.Context) {}
