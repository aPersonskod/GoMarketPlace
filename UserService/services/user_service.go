package services

import (
	"database/sql"
	"errors"
	"fmt"
	"user_service/types"

	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers() ([]types.UserDto, error)
	GetUserById(id string) (*types.UserDto, error)
	GetUserByEmail(email string) (*types.User, error)
	AddUser(u *types.User) (*types.UserDto, error)
	UpdateUser(u *types.UpdateUserDto) (*types.UserDto, error)
	DeleteUser(id string) error
	WalletReplenishment(id string, money int) (*types.UserDto, error)
	SpendMoney(id string, money int) (*types.UserDto, error)
}

type UserService struct {
	ConnStr string
}

func (service UserService) GetUsers() ([]types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM public.\"Users\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []types.UserDto{}
	for rows.Next() {
		u := types.User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u.GetUserDto())
	}
	return users, nil
}

func (service UserService) GetUserById(id string) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM public.\"Users\" WHERE \"Id\" = $1"
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := types.User{}
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

func (service UserService) GetUserByEmail(email string) (*types.User, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM public.\"Users\" WHERE \"Email\" = $1 limit 1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	u := types.User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	if u.Id == "" {
		//err = errors.New("Invalid user response")
		return nil, errors.New("Invalid user response")
	}
	return &u, nil
}

func (service UserService) AddUser(u *types.User) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	newId := uuid.New()
	u.Id = fmt.Sprintf("%s", newId)
	query := "INSERT INTO public.\"Users\" (\"Id\", \"Name\", \"Email\", \"Password\", \"Wallet\", \"Role\") VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = db.Exec(query, u.Id, u.Name, u.Email, u.Password, u.Wallet, u.Role)
	if err != nil {
		return nil, err
	}
	userDto := u.GetUserDto()
	return &userDto, nil
}

func (service UserService) UpdateUser(u *types.UpdateUserDto) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "UPDATE public.\"Users\" SET \"Name\" = $1, \"Email\" = $2, \"Role\" = $3 WHERE \"Id\" = $4"
	_, err = db.Exec(query, u.Name, u.Email, u.Role, u.Id)
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
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM public.\"Users\" WHERE \"Id\" = $1"
	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (service UserService) WalletReplenishment(id string, money int) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := "SELECT * FROM public.\"Users\" WHERE \"Id\" = $1"
	rows, err := db.Query(getQuery, id)
	if err != nil {
		return nil, err
	}

	u := types.User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	query := "UPDATE public.\"Users\" SET \"Wallet\" = $1"
	_, err = db.Exec(query, u.Wallet+money)
	if err != nil {
		return nil, err
	}

	userDto, err := service.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

func (service UserService) SpendMoney(id string, money int) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := "SELECT * FROM public.\"Users\" WHERE \"Id\" = $1"
	rows, err := db.Query(getQuery, id)
	if err != nil {
		return nil, err
	}

	u := types.User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.Wallet, &u.Role)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	query := "UPDATE public.\"Users\" SET \"Wallet\" = $1"
	_, err = db.Exec(query, u.Wallet-money)
	if err != nil {
		return nil, err
	}

	userDto, err := service.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}
