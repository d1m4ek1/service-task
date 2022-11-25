package response

import (
	"github.com/gin-gonic/gin"
	"service-task/models/payment"
)

type Response struct {
	Context    *gin.Context
	StatusCode int
	Error      string
	Result     payment.Payment
}

type SetResponse interface {
	SendResponse()
}

func (r *Response) SendResponse() {
	if r.Error == "" {
		r.Context.JSON(r.StatusCode, r.Result)
		return
	}

	r.Context.JSON(r.StatusCode, gin.H{
		"error": r.Error,
	})
}
