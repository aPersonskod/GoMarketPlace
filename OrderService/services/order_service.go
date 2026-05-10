package services

import (
	"database/sql"
	"fmt"
	"order_service/types"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type IOrderService interface {
	GetOrders(cartId string) ([]types.Order, error)                                                                   // TODO need to get from redis db
	AddOrder(orderDto types.OrderDto, cartService ICartService, productService IProductService) (*types.Order, error) // TODO need to get from redis db
	DeleteOrder(productId, cartId string, cartService ICartService, productService IProductService) error             // TODO need to get from redis db
}

type OrderService struct {
	DB *sql.DB
}

func (s OrderService) tableName() string {
	return "public.\"Orders\""
}

func (s OrderService) GetOrders(cartId string) ([]types.Order, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"CartId\" = $1 AND \"Quantity\" > 0", s.tableName())
	rows, err := s.DB.Query(query, cartId)
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

func (s OrderService) AddOrder(orderDto types.OrderDto, cartService ICartService, productService IProductService) (*types.Order, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"CartId\" = $1 AND \"OrderedProductId\" = $2", s.tableName())
	rows, err := s.DB.Query(query, orderDto.CartId, orderDto.OrderedProductId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	o := types.Order{}
	for rows.Next() {
		err = rows.Scan(&o.Id, &o.CartId, &o.OrderedProductId, &o.Quantity)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	if o.Id == "" {
		createdOrder, err := s.createOrder(orderDto)
		if err != nil {
			return nil, err
		}

		err = s.updateAmountToPay(createdOrder.CartId, cartService, productService)
		if err != nil {
			return nil, err
		}
		return createdOrder, nil
	}
	addedOrder, err := s.appendOrder(orderDto, o.Id)
	if err != nil {
		return nil, err
	}
	err = s.updateAmountToPay(addedOrder.CartId, cartService, productService)
	if err != nil {
		return nil, err
	}
	return addedOrder, nil
}

func (s OrderService) updateAmountToPay(cartId string, cartService ICartService, productService IProductService) error {
	amountToPay, err := s.getAmountToPay(cartId, productService)
	if err != nil {
		return err
	}
	err = cartService.UpdateAmountToPay(cartId, amountToPay)
	if err != nil {
		return err
	}
	return nil
}

func (s OrderService) createOrder(orderDto types.OrderDto) (*types.Order, error) {
	newId := fmt.Sprintf("%s", uuid.New())
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"CartId\", \"OrderedProductId\", \"Quantity\") VALUES ($1, $2, $3, $4)", s.tableName())
	_, err := s.DB.Exec(query, newId, orderDto.CartId, orderDto.OrderedProductId, orderDto.Quantity)
	if err != nil {
		return nil, err
	}
	o := types.Order{
		Id:               newId,
		CartId:           orderDto.CartId,
		OrderedProductId: orderDto.OrderedProductId,
		Quantity:         orderDto.Quantity,
	}
	return &o, nil
}

func (s OrderService) appendOrder(orderDto types.OrderDto, orderId string) (*types.Order, error) {
	query := fmt.Sprintf("UPDATE %s SET \"Quantity\" = $1 WHERE \"OrderedProductId\" = $2 AND \"CartId\" = $3", s.tableName())
	_, err := s.DB.Exec(query, orderDto.Quantity, orderDto.OrderedProductId, orderDto.CartId)
	if err != nil {
		return nil, err
	}
	o := types.Order{
		Id:               orderId,
		CartId:           orderDto.CartId,
		OrderedProductId: orderDto.OrderedProductId,
		Quantity:         orderDto.Quantity,
	}
	return &o, nil
}

func (s OrderService) getAmountToPay(cartId string, productService IProductService) (int, error) {
	orders, err := s.GetOrders(cartId)
	if err != nil {
		return 0, err
	}
	amountToPay := 0
	for _, order := range orders {
		product, err := productService.GetProduct(order.OrderedProductId)
		if err != nil {
			continue
		}
		amountToPay += product.Cost * order.Quantity
	}
	return amountToPay, nil
}

func (s OrderService) DeleteOrder(productId string, cartId string, cartService ICartService, productService IProductService) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE \"OrderedProductId\" = $1 AND \"CartId\" = $2", s.tableName())
	_, err := s.DB.Exec(query, productId, cartId)
	if err != nil {
		return err
	}
	err = s.updateAmountToPay(cartId, cartService, productService)
	if err != nil {
		return err
	}
	return nil
}
