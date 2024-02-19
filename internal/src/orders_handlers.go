package src

import (
	"api/internal/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) OrdersListHandler(c *gin.Context) {
	var listRequest models.RequestListsGeneral
	err := c.BindQuery(&listRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.OrdersResponseList{
			Message: err.Error(),
		})
		return
	}

	resp, err := s.listOrdersController(c, &listRequest)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.OrdersResponseList{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.OrdersResponseList{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) OrdersCreateHandler(c *gin.Context) {
	var orderItems models.OrdersRequestCreate
	err := c.ShouldBindJSON(&orderItems)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	resp, err := s.createOrdersController(c, &orderItems)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) OrdersUpdateHandler(c *gin.Context) {
	var ordersItem models.OrdersRequestUpdate
	err := c.ShouldBindJSON(&ordersItem)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}
	ordersItem.Id = c.Param("id")

	resp, err := s.updateOrdersController(c, &ordersItem)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) OrdersDeleteHandler(c *gin.Context) {
	orderId := c.Param("id")

	resp, err := s.deleteOrderController(c, orderId)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) OrdersGetByIdHandler(c *gin.Context) {
	orderId := c.Param("id")

	resp, err := s.getOrderByIdController(c, orderId)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) OrdersCheckCarsHandler(c *gin.Context) {
	carId := c.Param("car_id")
	pickupDate := c.Param("pickup_date")

	resp, err := s.checkCarsIsAlreadyOccupied(c, &models.RequestOrdersCheckOcupiedCars{
		CarId:      carId,
		PickupDate: pickupDate,
	})
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
