package src

import (
	"api/internal/models"
	"api/internal/utils"
	"database/sql"
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
		req.OrderBy = "car_name"
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
	count := 0
	if len(req.Search) > 0 {
		count++
		cmdQuery = fmt.Sprintf("%s WHERE LOWER(car_name) LIKE LOWER($%d)", cmdQuery, count)
		params = append(params, "%"+utils.Sanitize(req.Search)+"%")
	}

	// count all of search result
	total := 0
	queryCounter := fmt.Sprintf("SELECT COUNT(*) AS total FROM cars %s", cmdQuery)
	err := s.db.QueryRow(c, queryCounter, params...).Scan(&total)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if req.OrderBy != "" {
		cmdQuery = fmt.Sprintf("%s ORDER BY %s %s", cmdQuery, utils.Sanitize(req.OrderBy), req.Order)
	}

	count++
	cmdQuery = fmt.Sprintf("%s LIMIT $%d ", cmdQuery, count)
	params = append(params, req.Limit)

	count++
	cmdQuery = fmt.Sprintf("%s OFFSET $%d ", cmdQuery, count)
	params = append(params, (req.Page-1)*req.Limit)

	finalQuery := fmt.Sprintf("%s %s", query, cmdQuery)
	rows, err := s.db.Query(c, finalQuery, params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var id sql.NullInt64
	var dayRate, monthRate sql.NullFloat64
	var carName, image sql.NullString
	carsData := []*models.CarsItem{}
	for rows.Next() {
		item := models.CarsItem{}

		err = rows.Scan(
			&id,
			&carName,
			&dayRate,
			&monthRate,
			&image,
		)

		item.Id = int(id.Int64)
		item.CarName = strings.TrimSpace(carName.String)
		item.DayRate = dayRate.Float64
		item.MonthRate = monthRate.Float64
		item.Image = strings.TrimSpace(image.String)

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
		Message: "success",
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

	if req.Image == "" {
		errMsg = "missing-image"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	// TODO: need image save provider
	var carsId int
	err := s.db.QueryRow(c, "INSERT INTO cars (car_name, day_rate, month_rate, image) VALUES ($1, $2, $3, $4) RETURNING car_id", req.CarName, req.DayRate, req.MonthRate, req.Image).Scan(&carsId)
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

	count := 0
	if req.CarName != "" {
		count++
		set = append(set, fmt.Sprintf("car_name=$%d", count))
		params = append(params, req.CarName)
	}

	if req.Image != "" {
		count++
		set = append(set, fmt.Sprintf("image=$%d", count))
		params = append(params, req.Image)
	}

	if req.DayRate != 0 {
		count++
		set = append(set, fmt.Sprintf("day_rate=$%d", count))
		params = append(params, req.DayRate)
	}

	if req.MonthRate != 0 {
		count++
		set = append(set, fmt.Sprintf("month_rate=$%d", count))
		params = append(params, req.MonthRate)
	}

	count++
	query = fmt.Sprintf("%s %s WHERE car_id=$%d", query, strings.Join(set, ","), count)
	params = append(params, carId)

	log.Println(query)
	log.Println(params)
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

	_, err = s.db.Exec(c, "DELETE FROM cars WHERE car_id=$1", carId)
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

	var resp models.CarsResponseGet
	var idRes sql.NullInt64
	var dayRate, monthRate sql.NullFloat64
	var carName, image sql.NullString
	err = s.db.QueryRow(c, `
		SELECT 
			car_id, 
			car_name, 
			day_rate, 
			month_rate,
			image
		FROM cars WHERE car_id = $1
		`, carId).Scan(
		&idRes,
		&carName,
		&dayRate,
		&monthRate,
		&image,
	)

	resp.Item = &models.CarsItem{
		Id:        int(idRes.Int64),
		CarName:   strings.TrimSpace(carName.String),
		DayRate:   dayRate.Float64,
		MonthRate: monthRate.Float64,
		Image:     strings.TrimSpace(image.String),
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp.Message = "success"

	return &resp, nil
}
