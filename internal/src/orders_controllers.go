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
			order_id
			car_id,
			order_date,
			pickup_date,
			dropoff_date,
			pickup_location,
			dropoff_location
		FROM orders
	`
	var params []interface{}

	cmdQuery := ""
	if len(req.Search) > 0 {
		cmdQuery = fmt.Sprintf("%s WHERE LOWER(order_date) LIKE LOWER(?)", cmdQuery)
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

	cmdQuery = fmt.Sprintf("%s LIMIT ? OFFSET ?", cmdQuery)

	rows, err := s.db.Query(c, fmt.Sprintf("%s %s", query, cmdQuery), params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	ordersData := []*models.OrdersItem{}
	for rows.Next() {
		item := models.OrdersItem{}

		err = rows.Scan(
			&item.Id,
			&item.CarId,
			&item.OrderDate,
			&item.PickupDate,
			&item.DropoffDate,
			&item.PickupLocation,
			&item.DropoffLocation,
		)

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
	}, nil
}

func (s *Server) createOrdersController(c *gin.Context, req *models.OrdersRequestCreate) (*models.ResponseGeneral, error) {
	errMsg := ""
	if req.CarId == "" {
		errMsg = "missing-car-id"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.OrderDate == nil {
		errMsg = "missing-order-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.PickupDate == nil {
		errMsg = "missing-pickup-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.DropoffDate == nil {
		errMsg = "missing-dropoff-date"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.PickupLocation == "" {
		errMsg = "missing-pickup-location"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if req.DropoffDate == nil {
		errMsg = "missing-dropoff-location"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	var carsId int
	err := s.db.QueryRow(c, "INSERT INTO orders (car_id, order_date, pickup_date, dropoff_date, pickup_location, dropoff_location) VALUES (?, ?, ?, ?, ?, ?) RETURNING car_id", req.CarId, req.OrderDate, req.PickupDate, req.DropoffDate, req.PickupLocation, req.DropoffLocation).Scan(&carsId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id: carsId,
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

	if req.OrderDate != nil {
		set = append(set, "order_date=?")
		params = append(params, req.OrderDate)
	}

	if req.PickupDate != nil {
		set = append(set, "pickup_date=?")
		params = append(params, req.PickupDate)
	}

	if req.DropoffDate != nil {
		set = append(set, "dropoff_date=?")
		params = append(params, req.DropoffDate)
	}

	if req.PickupLocation != "" {
		set = append(set, "pickup_location=?")
		params = append(params, req.PickupLocation)
	}

	if req.DropoffLocation != "" {
		set = append(set, "dropoff_location=?")
		params = append(params, req.DropoffLocation)
	}

	query = fmt.Sprintf("%s %s WHERE order_id=?", query, strings.Join(set, ","))
	params = append(params, orderId)

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

	_, err = s.db.Exec(c, "DELETE FROM orders WHERE order_id=?", orderId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.ResponseGeneral{
		Id:      orderId,
		Message: "success",
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

	var item models.OrdersResponseGet
	err = s.db.QueryRow(c, `
		SELECT 
			order_id,
			car_id, 
			order_date, 
			pickup_date, 
			dropoff_date,
			pickup_location,
			dropoff_location
		FROM orders WHERE car_id = ?
		`, carId).Scan(
		&item.Item.Id,
		&item.Item.CarId,
		&item.Item.OrderDate,
		&item.Item.PickupDate,
		&item.Item.DropoffDate,
		&item.Item.PickupLocation,
		&item.Item.DropoffLocation,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	item.Message = "success"

	return &item, nil
}
