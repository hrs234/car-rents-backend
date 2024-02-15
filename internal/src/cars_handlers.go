package src

import (
	"api/internal/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) CarsListHandler(c *gin.Context) {
	var listRequest models.RequestListsGeneral
	err := c.BindQuery(&listRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.CarsResponseList{
			Message: err.Error(),
		})
		return
	}

	resp, err := s.listCarsController(c, &listRequest)
	if err != nil {
		if strings.Contains(err.Error(), "missing") {
			c.JSON(http.StatusBadRequest, &models.CarsResponseList{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &models.CarsResponseList{
			Message: err.Error(),
		})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.JSON(http.StatusOK, resp)
}

func (s *Server) CarsCreateHandler(c *gin.Context) {
	var carsItem models.CarsRequestCreate
	err := c.ShouldBindJSON(&carsItem)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}

	resp, err := s.createCarsController(c, &carsItem)
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

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.JSON(http.StatusOK, resp)
}

func (s *Server) CarsUpdateHandler(c *gin.Context) {
	var carsItem models.CarsRequestUpdate
	err := c.ShouldBindJSON(&carsItem)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &models.ResponseGeneral{
			Message: err.Error(),
		})
		return
	}
	carsItem.Id = c.Param("id")

	resp, err := s.updateCarsController(c, &carsItem)
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

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.JSON(http.StatusOK, resp)
}

func (s *Server) CarsDeleteHandler(c *gin.Context) {
	carId := c.Param("id")

	resp, err := s.deleteCarsController(c, carId)
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

func (s *Server) CarsGetByIdHandler(c *gin.Context) {
	carId := c.Param("id")
	log.Println(carId)

	resp, err := s.getCarsByIdController(c, carId)
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

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.JSON(http.StatusOK, resp)
}
