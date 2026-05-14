package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"user_service/configs"
	"user_service/middleware"
	"user_service/services"
	"user_service/types"

	docs "user_service/docs"

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

func getConnStr(dbHost, dbPort, dbUser, dbPassword, dbName string) string {
	//return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func initServices(store *MainStore) error {
	user_service = services.UserService{DB: store.DB}

	return nil
}

var user_service services.IUserService

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
		log.Fatalf("Migration up failed: %v", err)
	}
	fmt.Println("Migration up completed successfully")

	initServices(&MainStore{DB: db})

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
			eg.POST("/login", Login)
			eg.PUT("/", AddUser)
			eg.GET("/test", TestApi)
			//need auth
			wa := eg.Use(middleware.JwtAuthMiddleware())
			wa.GET("/get-all", GetUsers)
			wa.GET("/", GetUser)
			wa.PATCH("/", UpdateUser)
			wa.DELETE("/", DeleteUser)
			wa.POST("/wallet-replenishment", WalletReplenishment)
			wa.POST("/spend-money", SpendMoney)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%s", configs.Env.Port))
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
	role, _ := ctx.Get("role") // protected data
	if role != services.AdminRole {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have sufficient permissions"})
		return
	}
	users, err := user_service.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// @BasePath /api
// @Description description of function that get user by id
// @Tags user-service
// @Accept json
// @Produce json
/* // @Param   id	path	string		true	"Some ID" */
// @Success 200 {string} idk_WTF
// @Router /user-service/ [get]
func GetUser(ctx *gin.Context) {
	//id := ctx.Param("id")
	id, _ := ctx.Get("id") // protected data
	u, err := user_service.GetUserById(fmt.Sprint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description user registaration
// @Tags user-service
// @Accept json
// @Produce json
// @Param user	body	types.User	true	"User data"
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PUT]
func AddUser(ctx *gin.Context) {
	var user types.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	u, err := user_service.AddUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description update user
// @Tags user-service
// @Accept json
// @Produce json
// @Param user	body	types.UpdateUserDto	true	"User data"
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PATCH]
func UpdateUser(ctx *gin.Context) {
	var user types.UpdateUserDto
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	u, err := user_service.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
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
// @Router /user-service/ [DELETE]
func DeleteUser(ctx *gin.Context) {
	//id := ctx.Param("id")
	id, _ := ctx.Get("id")
	err := user_service.DeleteUser(fmt.Sprint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, "Successfuly deleted")
}

// @BasePath /api
// @Description add money
// @Tags user-service
// @Accept json
// @Param   money	query	int		false	"Some money"
// @Success 200 {string} Ok
// @Router /user-service/wallet-replenishment [POST]
func WalletReplenishment(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	money, err := strconv.Atoi(ctx.Query("money"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	u, err := user_service.WalletReplenishment(fmt.Sprint(id), money)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description spend money
// @Tags user-service
// @Accept json
// @Param   money	query	int		false	"Some money"
// @Success 200 {string} Ok
// @Router /user-service/spend-money [POST]
func SpendMoney(ctx *gin.Context) {
	//id := ctx.Query("id")
	id, _ := ctx.Get("id")
	money, err := strconv.Atoi(ctx.Query("money"))

	u, err := user_service.SpendMoney(fmt.Sprint(id), money)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description add money
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials	body	types.UserCredentials	true	"User data"
// @Success 200 {string} Ok
// @Router /user-service/login [POST]
func Login(ctx *gin.Context) {
	var credentials types.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := user_service.GetUserByEmail(credentials.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.Email != credentials.Email || user.Password != credentials.Password {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong email or password"})
		return
	}

	token, err := services.GenerateToken(user.Id, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
