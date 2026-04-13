package services

import (
	"database/sql"
	"fmt"
	"user_service/types"

	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers() ([]types.UserDto, error)
	GetUserById(id string) (*types.UserDto, error)
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

	query := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(query)
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
	userDto := u.GetUserDto()
	return &userDto, nil
}

func (service UserService) AddUser(u *types.User) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
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

func (service UserService) UpdateUser(u *types.UpdateUserDto) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
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
	db, err := sql.Open("postgres", service.ConnStr)
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

func (service UserService) WalletReplenishment(id string, money int) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
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

func (service UserService) SpendMoney(id string, money int) (*types.UserDto, error) {
	db, err := sql.Open("postgres", service.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	getQuery := fmt.Sprintf("SELECT * FROM public.\"Users\" WHERE \"Id\" = '%s'", id)
	rows, err := db.Query(getQuery)
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
