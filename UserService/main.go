package main

import (
	"database/sql"
	"fmt"
	"net/http"

	docs "marketplace/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var connStr string = "user=postgres password=password dbname=marketplace-users-db sslmode=disable"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
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

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/test-api", TestApi)
			eg.GET("/GetAll", GetUsers)
			eg.GET("/Get/:id", GetUserById)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}

// @BasePath /api/v1
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} EndpointTest
// @Router /example/test-api [get]
func TestApi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "This enpoint works!!!")
}

// @BasePath /api/v1
// @Summary GetAll
// @Schemes
// @Description description of function that get all users from DB
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} idk_WTF
// @Router /example/GetAll [get]
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

// @BasePath /api/v1
// @Description description of function that get user by id
// @Tags example
// @Accept json
// @Produce json
// @Param   id	path	string		true	"Some ID"
// @Success 200 {string} idk_WTF
// @Router /example/Get/{id} [get]
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
