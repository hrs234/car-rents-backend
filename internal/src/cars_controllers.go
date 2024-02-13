package src

import (
	"api/internal/models"
	"api/internal/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) listCarsController(c *gin.Context, req *models.RequestListsGeneral) (*models.CarsResponseList, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Order == "" || strings.ToUpper(req.Order) != "DESC" {
		req.Order = "ASC"
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}

	query := `
		SELECT
			car_id,
			car_name,
			day_rate,
			month_rate,
			image
		FROM cars
	`
	var params []interface{}

	cmdQuery := ""
	if len(req.Search) > 0 {
		cmdQuery = fmt.Sprintf("%s WHERE LOWER(car_name) LIKE LOWER(?)", cmdQuery)
		params = append(params, "%"+utils.Sanitize(req.Search)+"%")
	}

	// count all of search result
	total := 0
	err := s.db.QueryRow(c, fmt.Sprintf("SELECT COUNT(*) AS total FROM cars %s", cmdQuery), params...).Scan(&total)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.OrderBy != "" {
		cmdQuery = fmt.Sprintf("%s ORDER BY %s %s", cmdQuery, utils.Sanitize(req.OrderBy), req.Order)
	}

	cmdQuery = fmt.Sprintf("%s LIMIT ? OFFSET ?", cmdQuery)

	rows, err := s.db.Query(c, fmt.Sprintf("%s %s", query, cmdQuery), params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	carsData := []*models.CarsItem{}
	for rows.Next() {
		item := models.CarsItem{}

		err = rows.Scan(
			&item.Id,
			&item.CarName,
			&item.DayRate,
			&item.MonthRate,
			&item.Image,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		carsData = append(carsData, &item)
	}

	return &models.CarsResponseList{
		Total:   total,
		OrderBy: req.OrderBy,
		Order:   req.Order,
		Page:    req.Page,
		Limit:   req.Limit,
		Items:   carsData,
	}, nil
}

func (s *Server) createCarsController(c *gin.Context, req *models.CarsRequestCreate) (*models.ResponseGeneral, error) {
	errMsg := ""
	if req.CarName == "" {
		errMsg = "missing-car-name"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.DayRate == 0 {
		errMsg = "missing-day-rate"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.MonthRate == 0 {
		errMsg = "missing-month-rate"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	var carsId int
	err := s.db.QueryRow(c, "INSERT INTO cars (car_name, day_rate, month_rate, image) VALUES (?, ?, ?, ?) RETURNING car_id", req.CarName, req.DayRate, req.MonthRate, req.Image).Scan(&carsId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id: carsId,
	}, nil
}

func (s *Server) updateCarsController(c *gin.Context, req *models.CarsRequestUpdate) (*models.ResponseGeneral, error) {
	errorMsg := ""
	if req.Id == "" {
		errorMsg = "missing-cars-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	carId, err := strconv.Atoi(req.Id)
	if err != nil {
		errorMsg = "wrong-cars-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	query := "UPDATE cars SET"
	var params []interface{}
	var set []string

	if req.CarName != "" {
		set = append(set, "car_name=?")
		params = append(params, req.CarName)
	}

	if req.Image != "" {
		set = append(set, "image=?")
		params = append(params, req.Image)
	}

	if req.DayRate != 0 {
		set = append(set, "day_rate=?")
		params = append(params, req.DayRate)
	}

	if req.MonthRate != 0 {
		set = append(set, "month_rate=?")
		params = append(params, req.MonthRate)
	}

	query = fmt.Sprintf("%s %s WHERE car_id=?", query, strings.Join(set, ","))
	params = append(params, carId)

	_, err = s.db.Exec(c, query, params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      carId,
		Message: "success",
	}, nil
}

func (s *Server) deleteCarsController(c *gin.Context, id string) (*models.ResponseGeneral, error) {
	errorMsg := ""
	if id == "" {
		errorMsg = "missing-cars-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	carId, err := strconv.Atoi(id)
	if err != nil {
		errorMsg = "wrong-cars-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	_, err = s.db.Exec(c, "DELETE FROM cars WHERE car_id=?", carId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      carId,
		Message: "success",
	}, nil
}

func (s *Server) getCarsByIdController(c *gin.Context, id string) (*models.CarsResponseGet, error) {
	errorMsg := ""
	if id == "" {
		errorMsg = "missing-cars-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	carId, err := strconv.Atoi(id)
	if err != nil {
		errorMsg = "wrong-cars-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	var item models.CarsResponseGet
	err = s.db.QueryRow(c, `
		SELECT 
			id, 
			car_name, 
			day_rate, 
			month_rate,
			image
		FROM cars WHERE car_id = ?
		`, carId).Scan(
		&item.Item.Id,
		&item.Item.CarName,
		&item.Item.DayRate,
		&item.Item.MonthRate,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	item.Message = "success"

	return &item, nil
}
