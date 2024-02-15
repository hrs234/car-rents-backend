package src

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.Default())
	r.GET("/health", s.healthHandler)
	v1 := r.Group("/api/v1/")
	{
		v1.GET("/cars", s.CarsListHandler)
		v1.GET("/cars/:id", s.CarsGetByIdHandler)
		v1.POST("/cars", s.CarsCreateHandler)
		v1.PUT("/cars/:id", s.CarsUpdateHandler)
		v1.DELETE("/cars/:id", s.CarsDeleteHandler)

		v1.GET("/orders", s.OrdersListHandler)
		v1.GET("/orders/:id", s.OrdersGetByIdHandler)
		v1.POST("/orders", s.OrdersCreateHandler)
		v1.PUT("/orders/:id", s.OrdersUpdateHandler)
		v1.DELETE("/orders/:id", s.OrdersDeleteHandler)
	}
	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
