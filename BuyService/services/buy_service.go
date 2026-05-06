package services

import (
	"buy_service/types"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IBuyService interface {
	GetReportById(reportId, userId string) (*types.BuyReportDto, error)
	GetReportsByUserId(userId string) ([]types.BuyReport, error)
	BuyCart(cart types.Cart) error
}

type BuyService struct {
	ConnStr      string
	OrderService IOrderService
	UserService  IUserService
}

func (s BuyService) tableName() string {
	return "public.\"BuyReports\""
}

func (s BuyService) GetReportById(reportId, userId string) (*types.BuyReportDto, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", s.tableName())
	rows, err := db.Query(query, reportId)
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
	reportDto, err := s.getBuyReportDto(r, userId)
	if err != nil {
		return nil, err
	}
	return reportDto, nil
}

func (s BuyService) getBuyReportByCart(cart *types.Cart) (*types.BuyReport, error) {
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s WHERE \"CartId\" = $1", s.tableName())
	rows, err := db.Query(query, cart.Id)
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

func (s BuyService) GetReportsByUserId(userId string) ([]types.BuyReport, error) {
	boughtCarts, err := s.OrderService.GetBoughtCarts(userId)
	if err != nil {
		return nil, err
	}
	reports := []types.BuyReport{}
	for _, boughtCart := range boughtCarts {
		if boughtCart.PlaceId == "" {
			continue
		}
		r, err := s.getBuyReportByCart(&boughtCart)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		reports = append(reports, *r)
	}
	return reports, nil
}

// TODO need to add SAGA
func (s BuyService) BuyCart(cart types.Cart) error {
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
	db, err := sql.Open("postgres", s.ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO %s (\"Id\", \"CartId\", \"SaleDate\") VALUES ($1, $2, $3)", s.tableName())
	_, err = db.Exec(query, r.Id, r.CartId, r.SaleDate)
	if err != nil {
		return err
	}
	return nil
}

func (s BuyService) getBuyReportDto(buyReport types.BuyReport, userId string) (*types.BuyReportDto, error) {
	cart, err := s.OrderService.GetCart(userId)
	if err != nil {
		return nil, err
	}
	reportDto := types.BuyReportDto{Id: buyReport.Id, Cart: *cart, SaleDate: buyReport.SaleDate}
	return &reportDto, nil
}
