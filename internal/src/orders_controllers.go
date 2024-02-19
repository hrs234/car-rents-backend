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
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) listOrdersController(c *gin.Context, req *models.RequestListsGeneral) (*models.OrdersResponseList, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Order == "" || strings.ToUpper(req.Order) != "ASC" {
		req.Order = "DESC"
	}

	if req.OrderBy == "" {
		req.OrderBy = "order_date"
	}

	query := `
		SELECT
			order_id,
			orders.car_id,
			cars.car_name,
			order_date,
			pickup_date,
			dropoff_date,
			pickup_location,
			dropoff_location
		FROM orders JOIN cars ON orders.car_id=cars.car_id
	`
	var params []interface{}

	cmdQuery := ""
	count := 0
	if len(req.Search) > 0 {
		count++
		cmdQuery = fmt.Sprintf("%s WHERE LOWER(pickup_location) LIKE LOWER($%d)", cmdQuery, count)
		params = append(params, "%"+utils.Sanitize(req.Search)+"%")
	}

	// count all of search result
	total := 0
	err := s.db.QueryRow(c, fmt.Sprintf("SELECT COUNT(*) AS total FROM orders %s", cmdQuery), params...).Scan(&total)
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

	rows, err := s.db.Query(c, fmt.Sprintf("%s %s", query, cmdQuery), params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var id, carId sql.NullInt64
	var orderDate, pickupDate, dropoffDate sql.NullTime
	var pickupLocation, dropoffLocation, carName sql.NullString
	ordersData := []*models.OrdersItem{}
	for rows.Next() {
		item := models.OrdersItem{}

		err = rows.Scan(
			&id,
			&carId,
			&carName,
			&orderDate,
			&pickupDate,
			&dropoffDate,
			&pickupLocation,
			&dropoffLocation,
		)

		item.Id = int(id.Int64)
		item.CarId = int(carId.Int64)
		item.CarName = strings.TrimSpace(carName.String)
		item.OrderDate = orderDate.Time.Format("2006-01-02")
		item.PickupDate = pickupDate.Time.Format("2006-01-02")
		item.DropoffDate = dropoffDate.Time.Format("2006-01-02")
		item.PickupLocation = strings.TrimSpace(pickupLocation.String)
		item.DropoffLocation = strings.TrimSpace(dropoffLocation.String)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		ordersData = append(ordersData, &item)
	}

	return &models.OrdersResponseList{
		Total:   total,
		OrderBy: req.OrderBy,
		Order:   req.Order,
		Page:    req.Page,
		Limit:   req.Limit,
		Items:   ordersData,
		Message: "success",
	}, nil
}

func (s *Server) createOrdersController(c *gin.Context, req *models.OrdersRequestCreate) (*models.ResponseGeneral, error) {
	errMsg := ""
	if req.CarId == "" {
		errMsg = "missing-car-id"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	carIdNum, err := strconv.Atoi(req.CarId)
	if err != nil {
		errMsg = "wrong-cars-id-type"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.OrderDate == "" {
		errMsg = "missing-order-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	_, err = time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		errMsg = "failed-parsing-order-date"
		log.Println(err)
		return nil, errors.New(errMsg)
	}

	if req.PickupDate == "" {
		errMsg = "missing-pickup-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	_, err = time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		errMsg = "failed-parsing-pickup-date"
		log.Println(err)
		return nil, errors.New(errMsg)
	}

	if req.DropoffDate == "" {
		errMsg = "missing-dropoff-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	_, err = time.Parse("2006-01-02", req.DropoffDate)
	if err != nil {
		errMsg = "failed-parsing-dropoff-date"
		log.Println(err)
		return nil, errors.New(errMsg)
	}

	if req.PickupLocation == "" {
		errMsg = "missing-pickup-location"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.DropoffLocation == "" {
		errMsg = "missing-dropoff-location"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	resCheckCars, err := s.checkCarsIsAlreadyOccupied(c, &models.RequestOrdersCheckOcupiedCars{
		CarId:      req.CarId,
		PickupDate: req.PickupDate,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resCheckCars.Message == "car-already-occupied" {
		return nil, errors.New(resCheckCars.Message)
	}

	var orderId int
	err = s.db.QueryRow(c, "INSERT INTO orders (car_id, order_date, pickup_date, dropoff_date, pickup_location, dropoff_location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING order_id", carIdNum, req.OrderDate, req.PickupDate, req.DropoffDate, req.PickupLocation, req.DropoffLocation).Scan(&orderId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      orderId,
		Message: "success",
	}, nil
}

func (s *Server) updateOrdersController(c *gin.Context, req *models.OrdersRequestUpdate) (*models.ResponseGeneral, error) {
	errorMsg := ""
	if req.Id == "" {
		errorMsg = "missing-orders-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	orderId, err := strconv.Atoi(req.Id)
	if err != nil {
		errorMsg = "wrong-orders-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	query := "UPDATE orders SET"
	var params []interface{}
	var set []string
	count := 0

	if req.CarId != "" {
		carIdNum, err := strconv.Atoi(req.CarId)
		if err != nil {
			errorMsg = "wrong-car-id-type"
			log.Println(errorMsg)
			return nil, errors.New(errorMsg)
		}

		count++
		set = append(set, fmt.Sprintf("car_id=$%d", count))
		params = append(params, carIdNum)
	}

	if req.OrderDate != "" {
		_, err = time.Parse("2006-01-02", req.OrderDate)
		if err != nil {
			errorMsg = "failed-parsing-order-date"
			log.Println(err)
			return nil, errors.New(errorMsg)
		}

		count++
		set = append(set, fmt.Sprintf("order_date=$%d", count))
		params = append(params, req.OrderDate)
	}

	if req.PickupDate != "" {
		_, err = time.Parse("2006-01-02", req.PickupDate)
		if err != nil {
			errorMsg = "failed-parsing-pickup-date"
			log.Println(err)
			return nil, errors.New(errorMsg)
		}

		count++
		set = append(set, fmt.Sprintf("pickup_date=$%d", count))
		params = append(params, req.PickupDate)
	}

	if req.DropoffDate != "" {
		_, err = time.Parse("2006-01-02", req.DropoffDate)
		if err != nil {
			errorMsg = "failed-parsing-dropoff-date"
			log.Println(err)
			return nil, errors.New(errorMsg)
		}

		count++
		set = append(set, fmt.Sprintf("dropoff_date=$%d", count))
		params = append(params, req.DropoffDate)
	}

	if req.PickupLocation != "" {
		count++
		set = append(set, fmt.Sprintf("pickup_location=$%d", count))
		params = append(params, req.PickupLocation)
	}

	if req.DropoffLocation != "" {
		count++
		set = append(set, fmt.Sprintf("dropoff_location=$%d", count))
		params = append(params, req.DropoffLocation)
	}

	count++
	query = fmt.Sprintf("%s %s WHERE order_id=$%d", query, strings.Join(set, ","), count)
	params = append(params, orderId)

	if req.CarId != "" && req.PickupDate != "" {
		resCheckCars, err := s.checkCarsIsAlreadyOccupied(c, &models.RequestOrdersCheckOcupiedCars{
			CarId:      req.CarId,
			PickupDate: req.PickupDate,
		})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if resCheckCars.Message == "car-already-occupied" {
			return nil, errors.New(resCheckCars.Message)
		}
	}

	_, err = s.db.Exec(c, query, params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      orderId,
		Message: "success",
	}, nil
}

func (s *Server) deleteOrderController(c *gin.Context, id string) (*models.ResponseGeneral, error) {
	errorMsg := ""
	if id == "" {
		errorMsg = "missing-order-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		errorMsg = "wrong-order-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	_, err = s.db.Exec(c, "DELETE FROM orders WHERE order_id=$1", orderId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      orderId,
		Message: "success",
	}, nil
}

func (s *Server) checkCarsIsAlreadyOccupied(c *gin.Context, req *models.RequestOrdersCheckOcupiedCars) (*models.ResponseGeneral, error) {
	errorMsg := ""

	if req.CarId == "" {
		errorMsg = "missing-car-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	if req.PickupDate == "" {
		errorMsg = "missing-pickup-date"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	_, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		errorMsg = "failed-parsing-pickup-date"
		log.Println(err)
		return nil, errors.New(errorMsg)
	}

	// check the is the cars is ordered in order_date ?
	var usedCars int
	err = s.db.QueryRow(c, "SELECT COUNT(*) FROM orders WHERE dropoff_date >= $1 AND car_id=$2", req.PickupDate, req.CarId).Scan(&usedCars)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if usedCars > 0 {
		return &models.ResponseGeneral{
			Message: "car-already-occupied",
		}, nil
	}

	return &models.ResponseGeneral{
		Message: "car-available",
	}, nil
}

func (s *Server) getOrderByIdController(c *gin.Context, id string) (*models.OrdersResponseGet, error) {
	errorMsg := ""
	if id == "" {
		errorMsg = "missing-order-id"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	carId, err := strconv.Atoi(id)
	if err != nil {
		errorMsg = "wrong-order-id-type"
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	}

	var resp models.OrdersResponseGet
	var resId, resCarId sql.NullInt64
	var orderDate, pickupDate, dropoffDate sql.NullTime
	var pickupLocation, dropoffLocation, carName sql.NullString
	err = s.db.QueryRow(c, `
		SELECT 
			order_id,
			orders.car_id,
			cars.car_name,
			order_date, 
			pickup_date, 
			dropoff_date,
			pickup_location,
			dropoff_location
		FROM orders JOIN cars ON orders.car_id=cars.car_id WHERE orders.order_id = $1
		`, carId).Scan(
		&resId,
		&resCarId,
		&carName,
		&orderDate,
		&pickupDate,
		&dropoffDate,
		&pickupLocation,
		&dropoffLocation,
	)

	resp.Item = &models.OrdersItem{
		Id:              int(resId.Int64),
		CarId:           int(resCarId.Int64),
		CarName:         strings.TrimSpace(carName.String),
		OrderDate:       orderDate.Time.Format("2006-01-02"),
		PickupDate:      pickupDate.Time.Format("2006-01-02"),
		DropoffDate:     dropoffDate.Time.Format("2006-01-02"),
		PickupLocation:  strings.TrimSpace(pickupLocation.String),
		DropoffLocation: strings.TrimSpace(dropoffLocation.String),
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp.Message = "success"

	return &resp, nil
}
