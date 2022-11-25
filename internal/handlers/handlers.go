package handlers

import (
	"github.com/gin-gonic/gin"
	"service-task/internal/controllers/api"
)

func InitAPIHandlers(rtr *gin.Engine) {
	payment := rtr.Group("/api/payment")
	{
		payment.POST("/", api.Create())
		payment.GET("/:id", api.GetPayment())
	}
}
