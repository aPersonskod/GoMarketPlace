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
	GetReportById(reportId, userId string, userService IUserService, orderService IOrderService) (*types.BuyReportDto, error)
	GetReportsByUserId(userId string, userService IUserService, orderService IOrderService) ([]types.BuyReportDto, error)
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

// IT DOES NOT WORK, IDK WHY DO I NEED THIS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func (s BuyService) GetReportById(reportId, userId string, userService IUserService, orderService IOrderService) (*types.BuyReportDto, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", s.tableName())
	rows, err := s.DB.Query(query, reportId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	r := types.BuyReport{}
	for rows.Next() {
		err = rows.Scan(&r.Id, &r.CartId, &r.SaleDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	reportDto, err := s.getBuyReportDto(r, types.CartDto{}, userService, orderService) // get error, do not use it !!!
	if err != nil {
		return nil, err
	}
	return reportDto, nil
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

func (s BuyService) GetReportsByUserId(userId string, userService IUserService, orderService IOrderService) ([]types.BuyReportDto, error) {
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
		rDto, err := s.getBuyReportDto(*r, boughtCart, userService, orderService)
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

// TODO need to add SAGA
func (s BuyService) BuyCart(cart types.CartDto) error {
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

	userDto, err := s.UserService.SpendMoney(cart.UserId, cart.AmountToPay)
	if err != nil {
		return err
	}
	if userDto == nil {
		return fmt.Errorf("Can not find user and buy cart !!!")
	}

	err = s.addBuyReport(cart.Id)
	if err != nil {
		return err
	}
	err = s.OrderService.MarkCartAsBought(cart.Id)
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

func (s BuyService) getBuyReportDto(buyReport types.BuyReport, cart types.CartDto,
	userService IUserService, orderService IOrderService) (*types.BuyReportDto, error) {
	userDto, err := userService.GetUser()
	if err != nil {
		return nil, err
	}
	placeDto, err := orderService.GetPlace(cart.PlaceId)
	if err != nil {
		return nil, err
	}
	orders, err := orderService.GetOrders(cart.Id)
	if err != nil {
		return nil, err
	}
	orderDtos := []types.OrderEntity{}
	for _, order := range orders {
		productDto, err := orderService.GetProduct(order.OrderedProductId)
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
