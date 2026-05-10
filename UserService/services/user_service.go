package services

import (
	"database/sql"
	"errors"
	"fmt"
	"user_service/types"

	"github.com/google/uuid"
)

const UserRole string = "user"
const AdminRole string = "admin"
const tableName string = "public.\"Users\""

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
	DB *sql.DB
}

func (s UserService) GetUsers() ([]types.UserDto, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
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

func (s UserService) GetUserById(id string) (*types.UserDto, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", tableName)
	rows, err := s.DB.Query(query, id)
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

func (s UserService) GetUserByEmail(email string) (*types.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"Email\" = $1 limit 1", tableName)
	rows, err := s.DB.Query(query, email)
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

func (s UserService) AddUser(u *types.User) (*types.UserDto, error) {
	newId := uuid.New()
	u.Id = fmt.Sprintf("%s", newId)
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"Name\", \"Email\", \"Password\", \"Wallet\", \"Role\") VALUES ($1, $2, $3, $4, $5, $6)", tableName)
	_, err := s.DB.Exec(query, u.Id, u.Name, u.Email, u.Password, u.Wallet, u.Role)
	if err != nil {
		return nil, err
	}
	userDto := u.GetUserDto()
	return &userDto, nil
}

func (s UserService) UpdateUser(u *types.UpdateUserDto) (*types.UserDto, error) {
	query := "UPDATE public.\"Users\" SET \"Name\" = $1, \"Email\" = $2, \"Role\" = $3 WHERE \"Id\" = $4"
	_, err := s.DB.Exec(query, u.Name, u.Email, u.Role, u.Id)
	if err != nil {
		return nil, err
	}

	userDto, err := s.GetUserById(u.Id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

func (s UserService) DeleteUser(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE \"Id\" = $1", tableName)
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) WalletReplenishment(id string, money int) (*types.UserDto, error) {
	getQuery := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", tableName)
	rows, err := s.DB.Query(getQuery, id)
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

	query := fmt.Sprintf("UPDATE %s SET \"Wallet\" = $1", tableName)
	_, err = s.DB.Exec(query, u.Wallet+money)
	if err != nil {
		return nil, err
	}

	userDto, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

func (s UserService) SpendMoney(id string, money int) (*types.UserDto, error) {
	getQuery := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", tableName)
	rows, err := s.DB.Query(getQuery, id)
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

	query := fmt.Sprintf("UPDATE %s SET \"Wallet\" = $1", tableName)
	_, err = s.DB.Exec(query, u.Wallet-money)
	if err != nil {
		return nil, err
	}

	userDto, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return userDto, nil
}
