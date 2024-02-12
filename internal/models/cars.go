package models

type CarsItem struct {
	Id        string  `json:"id"`
	CarName   string  `json:"car_name"`
	DayRate   float64 `json:"day_rate"`
	MonthRate float64 `json:"month_rate"`
	Image     string  `json:"image"`
}

type CarsResponseList struct {
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
	Order   string      `json:"order"`
	OrderBy string      `json:"order_by"`
	Items   []*CarsItem `json:"items"`
}

type CarsRequestCreate struct {
	CarName   string  `json:"car_name"`
	DayRate   float64 `json:"day_rate"`
	MonthRate float64 `json:"month_rate"`
	Image     string  `json:"image"`
}

type CarsRequestUpdate struct {
	Id        string  `json:"id"`
	CarName   string  `json:"car_name"`
	DayRate   float64 `json:"day_rate"`
	MonthRate float64 `json:"month_rate"`
	Image     string  `json:"image"`
}

type CarsRequestDelete struct {
	Id string `json:"id"`
}

type CarsResponseGet struct {
	Message string    `json:"message"`
	Item    *CarsItem `json:"item"`
}
