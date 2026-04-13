package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	docs "user_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var connStr string = "user=postgres password=password dbname=marketplace-users-db sslmode=disable"

type IUserService interface {
	GetUsers() ([]UserDto, error)
	GetUserById(id string) (*UserDto, error)
	AddUser(u *User) (*UserDto, error)
	UpdateUser(u *UpdateUserDto) (*UserDto, error)
	DeleteUser(id string) error
	WalletReplenishment(id string, money int) (*UserDto, error)
	SpendMoney(id string, money int) (*UserDto, error)
}

var Service IUserService

type UserService struct{}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Wallet   int    `json:"wallet"`
	Role     string `json:"role"`
}

type UpdateUserDto struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserDto struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Wallet int    `json:"wallet"`
	Role   string `json:"role"`
}

type UserDtoAdapter interface {
	GetUserDto() UserDto
}

func (user *User) GetUserDto() UserDto {
	return UserDto{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Wallet: user.Wallet,
		Role:   user.Role,
	}
}

func main() {
	fmt.Println("Hello blya")
	Service = UserService{}
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

func (service UserService) GetUsers() ([]UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM public.\"Users\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []UserDto{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u.GetUserDto())
	}
	return users, nil
}

func (service UserService) GetUserById(id string) (*UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	u := User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	userDto := u.GetUserDto()
	return &userDto, nil
}

func (service UserService) AddUser(u *User) (*UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	newId := uuid.New()
	u.Id = fmt.Sprintf("%s", newId)
	query := fmt.Sprintf("INSERT INTO public.\"Users\" (\"Id\", \"Name\", \"Email\", \"Password\", \"Wallet\", \"Role\") VALUES ('%s','%s', '%s', '%s', %d, '%s')",
		u.Id, u.Name, u.Email, u.Password, u.Wallet, u.Role)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	userDto := u.GetUserDto()
	return &userDto, nil
}
func (service UserService) UpdateUser(u *UpdateUserDto) (*UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE public.\"Users\" SET \"Name\" = '%s', \"Email\" = '%s', \"Role\" = '%s' WHERE \"Id\" = '%s'",
		u.Name, u.Email, u.Role, u.Id)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	userDto, err := service.GetUserById(u.Id)
	if err != nil {
		return userDto, err
	}
	return userDto, nil
}
func (service UserService) DeleteUser(id string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (service UserService) WalletReplenishment(id string, money int) (*UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	userDto, err := service.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}
func (service UserService) SpendMoney(id string, money int) (*UserDto, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	userDto, err := service.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
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

	service := UserService{}
	users, err := service.GetUsers()
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
	var user User
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
	var user UpdateUserDto
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
