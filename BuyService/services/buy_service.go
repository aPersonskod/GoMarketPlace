package services

import (
	"buy_service/types"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type IBuyService interface {
	GetReportsByUserId(userId string) ([]types.BuyReportDto, error)
	BuyCart(cart types.CartDto) error
}

type BuyService struct {
	DB           *sql.DB
	OrderService IOrderService
	UserService  IUserService
}

func (s BuyService) tableName() string {
	return "public.\"BuyReports\""
}

func (s BuyService) getBuyReportByCart(cart *types.CartDto) (*types.BuyReport, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"CartId\" = $1", s.tableName())
	rows, err := s.DB.Query(query, cart.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := types.BuyReport{}
	for rows.Next() {
		err = rows.Scan(&r.Id, &r.CartId, &r.SaleDate)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	return &r, nil
}

func (s BuyService) GetReportsByUserId(userId string) ([]types.BuyReportDto, error) {
	boughtCarts, err := s.OrderService.GetBoughtCarts(userId)
	if err != nil {
		return nil, err
	}
	reports := []types.BuyReportDto{}
	for _, boughtCart := range boughtCarts {
		if boughtCart.PlaceId == "" {
			continue
		}
		r, err := s.getBuyReportByCart(&boughtCart)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		rDto, err := s.getBuyReportDto(*r, boughtCart)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *rDto)
	}
	slices.SortFunc(reports, func(a, b types.BuyReportDto) int {
		if a.SaleDate.Before(b.SaleDate) {
			return 1
		}
		if b.SaleDate.Before(a.SaleDate) {
			return -1
		}
		return 0
	})
	return reports, nil
}

func (s BuyService) BuyCart(cart types.CartDto) error {
	if cart.IsConfirmed != true {
		return fmt.Errorf("Can't buy cart, cart is not confirmed !!!")
	}
	if cart.PlaceId == "" {
		return fmt.Errorf("Can't buy cart, place is required !!!")
	}
	orders, err := s.OrderService.GetOrders(cart.Id)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return fmt.Errorf("Can't buy cart, cart have not orders !!!")
	}

	// something important and very slow
	time.Sleep(time.Second * 5)

	err = s.addBuyReport(cart.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s BuyService) addBuyReport(cartId string) error {
	newId := fmt.Sprint(uuid.New())
	r := types.BuyReport{
		Id:       newId,
		CartId:   cartId,
		SaleDate: time.Now(),
	}
	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"CartId\", \"SaleDate\") VALUES ($1, $2, $3)", s.tableName())
	_, err := s.DB.Exec(query, r.Id, r.CartId, r.SaleDate)
	if err != nil {
		return err
	}
	return nil
}

func (s BuyService) getBuyReportDto(buyReport types.BuyReport, cart types.CartDto) (*types.BuyReportDto, error) {
	userDto, err := s.UserService.GetUser()
	if err != nil {
		return nil, err
	}
	placeDto, err := s.OrderService.GetPlace(cart.PlaceId)
	if err != nil {
		return nil, err
	}
	orders, err := s.OrderService.GetOrders(cart.Id)
	if err != nil {
		return nil, err
	}
	orderDtos := []types.OrderEntity{}
	for _, order := range orders {
		productDto, err := s.OrderService.GetProduct(order.OrderedProductId)
		if err != nil {
			return nil, err
		}
		orderDto := types.OrderEntity{
			Id:             order.Id,
			OrderedProduct: *productDto,
			Quantity:       order.Quantity,
		}
		orderDtos = append(orderDtos, orderDto)
	}

	cartEntity := types.CartEntity{
		Id:          cart.Id,
		User:        *userDto,
		Place:       *placeDto,
		Orders:      orderDtos,
		AmountToPay: cart.AmountToPay,
		IsConfirmed: cart.IsConfirmed,
		IsBought:    cart.IsBought,
	}
	reportDto := types.BuyReportDto{Id: buyReport.Id, Cart: cartEntity, SaleDate: buyReport.SaleDate}
	return &reportDto, nil
}
