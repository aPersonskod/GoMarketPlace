package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	docs "marketplace/docs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var connStr string = "user=postgres password=password dbname=marketplace-users-db sslmode=disable"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   int    `json:"wallet"`
	Role     string `json:"role"`
}

func main() {
	fmt.Println("Hello blya")
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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM public.\"Users\"")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	u := User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	ctx.JSON(http.StatusOK, u)
}

// @BasePath /api
// @Description add user
// @Tags user-service
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PUT]
func AddUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	newId := uuid.New()
	user.Id = fmt.Sprintf("%s", newId)
	query := fmt.Sprintf("INSERT INTO public.\"Users\" (\"Id\", \"Name\", \"Email\", \"Password\", \"Wallet\", \"Role\") VALUES ('%s','%s', '%s', '%s', %d, '%s')",
		user.Id, user.Name, user.Email, user.Password, user.Wallet, user.Role)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, user)
}

// @BasePath /api
// @Description update user
// @Tags user-service
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /user-service/ [PATCH]
func UpdateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE public.\"Users\" SET \"Name\" = '%s', \"Email\" = '%s', \"Password\" = '%s', \"Wallet\" = %d, \"Role\" = '%s'",
		user.Name, user.Email, user.Password, user.Wallet, user.Role)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, user)
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, "")
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
	if err != nil {
		panic(err)
	}

	u := User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	query := fmt.Sprintf("UPDATE public.\"Users\" SET \"Wallet\" = %d", u.Wallet+money)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	result := fmt.Sprintf("wallet before: %d, wallet after: %d", u.Wallet, u.Wallet+money)
	ctx.JSON(http.StatusOK, result)
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

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
	if err != nil {
		panic(err)
	}

	u := User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	query := fmt.Sprintf("UPDATE public.\"Users\" SET \"Wallet\" = %d", u.Wallet-money)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	result := fmt.Sprintf("wallet before: %d, wallet after: %d", u.Wallet, u.Wallet-money)
	ctx.JSON(http.StatusOK, result)
}
