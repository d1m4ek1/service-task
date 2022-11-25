package server

import (
	"github.com/gin-gonic/gin"
	"service-task/internal/handlers"
	"service-task/pkg/seterror"
)

func InitServer() error {
	gin.SetMode(gin.ReleaseMode)
	rtr := gin.Default()

	handlers.InitAPIHandlers(rtr)

	if err := rtr.Run(":8029"); err != nil {
		seterror.SetAppError("gin.Run", err)
		return err
	}

	return nil
}
