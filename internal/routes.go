package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.GET("/health", s.healthHandler)
	v1 := r.Group("/api/v1/")
	{
		v1.GET("/cars", s.ListCarsHandler)
		v1.GET("/cars/:id", s.GetCarsByIDHandler)
		v1.POST("/cars", s.CreateCarsHandler)
		v1.PUT("/cars/:id", s.UpdateCarsHandler)
		v1.DELETE("/cars/:id", s.DeleteCarsHandler)

		v1.GET("/orders", s.ListOrdersHandler)
		v1.GET("/orders/:id", s.GetOrdersByIDHandler)
		v1.POST("/orders", s.CreateOrdersHandler)
		v1.PUT("/orders/:id", s.UpdateOrdersHandler)
		v1.DELETE("/orders/:id", s.DeleteOrdersHandler)
	}
	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
