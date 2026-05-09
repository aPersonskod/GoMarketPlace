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
	GetBoughtCarts(userId string) ([]types.Cart, error)
	ConfirmAndBuyCart(placeId string, userId string, orderService IOrderService, userService IUserService, buyService IBuyService) (*types.Cart, error)
	MarkCartAsBought(cartId string) error
	UpdateAmountToPay(cartId string, amountToPay int) error
}

type CartService struct {
	DB *sql.DB
}

func (s CartService) tableName() string {
	return "public.\"ShoppingCarts\""
}

func (s CartService) GetCart(userId string) (*types.Cart, error) {
	c, err := s.getCartFromDB(userId, false)
	if err != nil {
		return nil, err
	}
	if c.Id == "" {
		newCart, err := s.addCart(userId)
		if err != nil {
			return nil, fmt.Errorf("can't create new cart: %s", err.Error())
		}
		return newCart, nil
	}
	return c, nil
}

func (s CartService) GetBoughtCart(userId string) (*types.Cart, error) {
	c, err := s.getCartFromDB(userId, false)
	if err != nil {
		return nil, err
	}
	if c.Id == "" {
		newCart, err := s.addCart(userId)
		if err != nil {
			return nil, fmt.Errorf("can't create new cart: %s", err.Error())
		}
		return newCart, nil
	}
	return c, nil
}

func (s CartService) getCartFromDB(userId string, isConfirmed bool) (*types.Cart, error) {
	query := fmt.Sprintf("SELECT \"Id\", \"UserId\", \"PlaceId\", \"AmountToPay\", \"IsConfirmed\", \"IsBought\" FROM %s WHERE \"UserId\" = $1 AND \"IsConfirmed\" = $2 AND \"IsBought\" = 'false'", s.tableName())
	rows, err := s.DB.Query(query, userId, isConfirmed)
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
	return &c, nil
}

func (s CartService) GetBoughtCarts(userId string) ([]types.Cart, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"UserId\" = $1 AND \"IsBought\" = 'true'", s.tableName())
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carts := []types.Cart{}
	for rows.Next() {
		c := types.Cart{}
		err = rows.Scan(&c.Id, &c.UserId, &c.PlaceId, &c.AmountToPay, &c.IsConfirmed, &c.IsBought)
		if err != nil {
			fmt.Println(err)
			continue
		}
		carts = append(carts, c)
	}
	return carts, nil
}

func (s CartService) ConfirmAndBuyCart(placeId string, userId string, orderService IOrderService, userService IUserService, buyService IBuyService) (*types.Cart, error) {
	cart, err := s.GetCart(userId)
	if err != nil {
		return nil, err
	}
	orders, err := orderService.GetOrders(cart.Id)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("Cart has no orders !!!")
	}
	user, err := userService.GetUser()
	if err != nil {
		return nil, err
	}
	if user.Wallet < cart.AmountToPay {
		return nil, fmt.Errorf("You have not enough money")
	}
	confirmedCart, err := s.confirmCart(placeId, userId, cart.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println(*confirmedCart.PlaceId)
	err = buyService.BuyCart(*confirmedCart)
	if err != nil {
		return nil, err
	}
	return confirmedCart, nil
}
func (s CartService) MarkCartAsBought(cartid string) error {
	query := fmt.Sprintf("UPDATE %s SET \"IsBought\" = $1 WHERE \"Id\" = $2", s.tableName())
	_, err := s.DB.Exec(query, true, cartid)
	if err != nil {
		return err
	}
	return nil
}

func (s CartService) confirmCart(placeId, userId, cartId string) (*types.Cart, error) {
	query := fmt.Sprintf("UPDATE %s SET \"PlaceId\" = $1 WHERE \"UserId\" = $2 AND \"Id\" = $3", s.tableName())
	res, err := s.DB.Exec(query, placeId, userId, cartId)
	if err != nil {
		return nil, err
	}
	if r, e := res.RowsAffected(); r == 0 || e != nil {
		return nil, fmt.Errorf("Place does not set!!!")
	}

	foundCart, err := s.getCartFromDB(userId, false)
	if err != nil {
		return nil, err
	}

	confirmQuery := fmt.Sprintf("UPDATE %s SET \"IsConfirmed\" = $1 WHERE \"UserId\" = $2", s.tableName())
	_, err = s.DB.Exec(confirmQuery, true, userId)
	if err != nil {
		return nil, err
	}

	foundCart, err = s.getCartFromDB(userId, true)
	if err != nil {
		return nil, err
	}

	return foundCart, nil
}
func (s CartService) addCart(userId string) (*types.Cart, error) {
	newId := uuid.New()
	c := types.Cart{
		Id:     fmt.Sprintf("%s", newId),
		UserId: userId,
	}
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"UserId\") VALUES ($1, $2)", s.tableName())
	_, err := s.DB.Exec(query, c.Id, c.UserId)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s CartService) UpdateAmountToPay(cartId string, amountToPay int) error {
	query := fmt.Sprintf("UPDATE %s SET \"AmountToPay\" = $1 WHERE \"Id\" = $2", s.tableName())
	_, err := s.DB.Exec(query, amountToPay, cartId)
	if err != nil {
		return err
	}
	return nil
}
