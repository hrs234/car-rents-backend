package models

import "time"

type OrdersItem struct {
	Id              string     `json:"id"`
	CarId           string     `json:"car_id"`
	OrderDate       *time.Time `json:"order_date"`
	PickupDate      *time.Time `json:"pickup_date"`
	DropoffDate     *time.Time `json:"dropoff_date"`
	PickupLocation  string     `json:"pickup_location"`
	DropoffLocation string     `json:"dropoff_location"`
}

type OrdersResponseList struct {
	Page    int           `json:"page"`
	Limit   int           `json:"limit"`
	Total   int           `json:"total"`
	Order   string        `json:"order"`
	OrderBy string        `json:"order_by"`
	Items   []*OrdersItem `json:"items"`
}

type OrdersRequestCreate struct {
	CarId           string     `json:"car_id"`
	OrderDate       *time.Time `json:"order_date"`
	PickupDate      *time.Time `json:"pickup_date"`
	DropoffDate     *time.Time `json:"dropoff_date"`
	PickupLocation  string     `json:"pickup_location"`
	DropoffLocation string     `json:"dropoff_location"`
}

type OrdersRequestUpdate struct {
	Id              string     `json:"id"`
	CarId           string     `json:"car_id"`
	OrderDate       *time.Time `json:"order_date"`
	PickupDate      *time.Time `json:"pickup_date"`
	DropoffDate     *time.Time `json:"dropoff_date"`
	PickupLocation  string     `json:"pickup_location"`
	DropoffLocation string     `json:"dropoff_location"`
}

type OrdersRequestDelete struct {
	Id string `json:"id"`
}

type OrdersResponseGet struct {
	Message string      `json:"message"`
	Item    *OrdersItem `json:"item"`
}
