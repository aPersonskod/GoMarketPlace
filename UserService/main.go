package main

import (
	"fmt"
	"net/http"
	"strconv"
	"user_service/services"
	"user_service/types"

	docs "user_service/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var connStr string = "user=postgres password=password dbname=marketplace-users-db sslmode=disable"

var Service services.IUserService

func main() {
	fmt.Println("Hello blya")
	Service = services.UserService{ConnStr: connStr}
	createGin()
}

func createGin() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	docs.SwaggerInfo.BasePath = "/api"
	v1 := r.Group("/api")
	{
		eg := v1.Group("/user-service")
		{
			eg.GET("/test", TestApi)
			eg.GET("/get-all", GetUsers)
			eg.GET("/:id", GetUserById)
			eg.PUT("/", AddUser)
			eg.PATCH("/", UpdateUser)
			eg.DELETE("/:id", DeleteUser)
			eg.POST("/wallet-replenishment", WalletReplenishment)
			eg.POST("/spend-money", SpendMoney)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}

// @BasePath /api
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags user-service
// @Accept json
// @Produce json
// @Success 200 {string} EndpointTest
// @Router /user-service/test [get]
func TestApi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "This enpoint works!!!")
}

// @BasePath /api
// @Summary GetAll
// @Schemes
// @Description description of function that get all users from DB
// @Tags user-service
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /user-service/get-all [get]
func GetUsers(ctx *gin.Context) {

	users, err := Service.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, users)
}

// @BasePath /api
// @Description description of function that get user by id
// @Tags user-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Some ID"
// @Success 200 {string} idk_WTF
// @Router /user-service/{id} [get]
func GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	u, err := Service.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description add user
// @Tags user-service
// @Accept json
// @Produce json
// @Param user	body	User	true	"User data"
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PUT]
func AddUser(ctx *gin.Context) {
	var user types.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	u, err := Service.AddUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description update user
// @Tags user-service
// @Accept json
// @Produce json
// @Param user	body	UpdateUserDto	true	"User data"
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PATCH]
func UpdateUser(ctx *gin.Context) {
	var user types.UpdateUserDto
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	u, err := Service.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description delete user
// @Tags user-service
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Some ID"
// @Success 200 {string} idk_WTF
// @Router /user-service/{id} [DELETE]
func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := Service.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, "Successfuly deleted")
}

// @BasePath /api
// @Description add money
// @Tags user-service
// @Accept json
// @Param   id		query	string	false	"Some ID"
// @Param   money	query	int		false	"Some money"
// @Success 200 {string} Ok
// @Router /user-service/wallet-replenishment [POST]
func WalletReplenishment(ctx *gin.Context) {
	id := ctx.Query("id")
	money, err := strconv.Atoi(ctx.Query("money"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	u, err := Service.WalletReplenishment(id, money)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description spend money
// @Tags user-service
// @Accept json
// @Param   id		query	string	false	"Some ID"
// @Param   money	query	int		false	"Some money"
// @Success 200 {string} Ok
// @Router /user-service/spend-money [POST]
func SpendMoney(ctx *gin.Context) {
	id := ctx.Query("id")
	money, err := strconv.Atoi(ctx.Query("money"))

	u, err := Service.SpendMoney(id, money)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, u)
}
