package services

import (
	"database/sql"
	"errors"
	"fmt"
	"order_service/types"

	_ "github.com/lib/pq"
)

type IPlaceService interface {
	GetPlaces() ([]types.Place, error)
	GetPlace(placeId string) (*types.Place, error)
}

type PlaceService struct {
	DB *sql.DB
}

func (s PlaceService) tableName() string {
	return "public.\"Places\""
}

func (s PlaceService) GetPlaces() ([]types.Place, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM %s", s.tableName()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places := []types.Place{}
	for rows.Next() {
		p := types.Place{}
		err = rows.Scan(&p.Id, &p.Address, &p.WorkingTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		places = append(places, p)
	}
	return places, nil
}

func (s PlaceService) GetPlace(placeId string) (*types.Place, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE \"Id\" = $1", s.tableName())
	rows, err := s.DB.Query(query, placeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p := types.Place{}
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Address, &p.WorkingTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	if p.Id == "" {
		return nil, errors.New("Place not found")
	}
	return &p, nil
}
