package services

import (
	"database/sql"
	"fmt"
	"order_service/types"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type IOrderService interface {
	GetOrders(cartId string) ([]types.Order, error)                               // TODO need to get from redis db
	AddOrder(productId string, cartId string, quantity int) (*types.Order, error) // TODO need to get from redis db
	DeleteOrder(productId string) error                                           // TODO need to get from redis db
}

type OrderService struct {
	ConnStr string
}

func (s OrderService) tableName() string {
	return "public.\"Orders\""
}

func (s OrderService) GetOrders(cartId string) ([]types.Order, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s WHERE \"CartId\" = $1 AND \"Quantity\" > 0", s.tableName())
	rows, err := db.Query(query, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []types.Order{}
	for rows.Next() {
		o := types.Order{}
		err = rows.Scan(&o.Id, &o.CartId, &o.OrderedProductId, &o.Quantity)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s OrderService) AddOrder(productId string, cartId string, quantity int) (*types.Order, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	newId := fmt.Sprintf("%s", uuid.New())
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"CartId\", \"OrderedProductId\", \"Quantity\") VALUES ($1, $2, $3, $4)", s.tableName())
	_, err = db.Exec(query, newId, cartId, productId, quantity)
	if err != nil {
		return nil, err
	}
	o := types.Order{
		Id:               newId,
		CartId:           cartId,
		OrderedProductId: productId,
		Quantity:         quantity,
	}
	return &o, nil
}
func (s OrderService) DeleteOrder(productId string) error {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM %s WHERE \"Id\" = $1", s.tableName())
	_, err = db.Exec(query, productId)
	if err != nil {
		return err
	}
	return nil
}
