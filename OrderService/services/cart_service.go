package services

import (
	"database/sql"
	"errors"
	"fmt"
	"order_service/types"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ICartService interface {
	GetCart(userId string) (*types.Cart, error)
	ConfirmCart(placeId string, userId string) (*types.Cart, error)
	MarkCartAsBought(cartid string) error
}

type CartService struct {
	ConnStr      string
	OrderService IOrderService
	UserService  IUserService
}

func (s CartService) tableName() string {
	return "public.\"ShoppingCarts\""
}

func (s CartService) GetCart(userId string) (*types.Cart, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s WHERE \"UserId\" = $1 AND \"IsConfirmed\" = 'false' AND \"IsBought\" = 'false'", s.tableName())
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	c := types.Cart{}
	for rows.Next() {
		err = rows.Scan(&c.Id, &c.UserId, &c.PlaceId, &c.AmountToPay, &c.IsConfirmed, &c.IsBought)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	if c.Id == "" {
		newCart, err := s.addCart(userId)
		if err != nil {
			return nil, fmt.Errorf("can't create new cart: %s", err.Error())
		}
		return newCart, nil
	}
	return &c, nil
}
func (s CartService) ConfirmAndBuyCart(placeId string, userId string) (*types.Cart, error) {
	cart, err := s.GetCart(userId)
	if err != nil {
		return nil, err
	}
	orders, err := s.OrderService.GetOrders(cart.Id)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("Cart has no orders !!!")
	}
	user, err := s.UserService.GetUser()
	if err != nil {
		return nil, err
	}
	if user.Wallet < cart.AmountToPay {
		return nil, fmt.Errorf("You have not enough money")
	}
	confirmedCart, err := s.confirmCart(placeId, userId)
	if err != nil {
		return nil, err
	}
	fmt.Println(confirmedCart.PlaceId)
	// TODO buy cart
	return confirmedCart, nil
}
func (s CartService) MarkCartAsBought(cartid string) error {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE %s SET \"IsBought\" = $1 WHERE \"Id\" = $2", s.tableName())
	_, err = db.Exec(query, true, cartid)
	if err != nil {
		return err
	}
	return nil
}

func (s CartService) confirmCart(placeId string, userId string) (*types.Cart, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE %s SET \"PlaceId\" = $1, \"IsConfirmed\" = $2", s.tableName())
	_, err = db.Exec(query, placeId, true)
	if err != nil {
		return nil, err
	}

	foundCart, err := s.GetCart(userId)
	if err != nil {
		return nil, err
	}
	return foundCart, nil
}
func (s CartService) addCart(userId string) (*types.Cart, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	newId := uuid.New()
	c := types.Cart{
		Id:     fmt.Sprintf("%s", newId),
		UserId: userId,
	}
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"UserId\") VALUES ($1, $2)", s.tableName())
	_, err = db.Exec(query, c.Id, c.UserId)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
