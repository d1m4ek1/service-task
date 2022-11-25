package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"service-task/internal/response"
	"service-task/models/payment"
	"service-task/pkg/seterror"
)

func findPaymentByID(paymentItems []payment.Payment, id string) (payment.Payment, bool) {
	for _, item := range paymentItems {
		if item.Id == id {
			return item, true
		}
	}

	return payment.Payment{}, false
}

func readLocalDataPaymentsKind() (payment.PaymentsKind, error) {
	var paymentsKind payment.PaymentsKind

	file, err := os.ReadFile("./localData/payments-kind.json")
	if err != nil {
		seterror.SetAppError("os.ReadFile", err)
		return payment.PaymentsKind{}, err
	}

	if err := json.Unmarshal(file, &paymentsKind); err != nil {
		seterror.SetAppError("json.Unmarshal", err)
		return payment.PaymentsKind{}, err
	}

	return paymentsKind, nil
}

func getKindData(kind string) (int64, error) {
	paymentsKind, err := readLocalDataPaymentsKind()
	if err != nil {
		seterror.SetAppError("readLocalDataPaymentsKind()", err)
		return 0, err
	}

	switch kind {
	case "designer":
		return paymentsKind.Designer, nil
	case "accountant":
		return paymentsKind.Accountant, nil
	case "editor":
		return paymentsKind.Editor, nil
	}

	return 0, nil
}

func Create() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		var sendResponse = &response.Response{
			Context:    context,
			StatusCode: 200,
		}

		var kind payment.Kind
		var paymentsAmount int64
		var err error

		if err := json.NewDecoder(context.Request.Body).Decode(&kind); err != nil {
			seterror.SetAppError("json.NewDecoder().Decoder()", err)
			sendResponse.Error = "bad_request"
			sendResponse.StatusCode = http.StatusBadRequest
			sendResponse.SendResponse()
			return
		}

		paymentsAmount, err = getKindData(kind.Kind)
		if err != nil {
			seterror.SetAppError("getKindData()", err)
			sendResponse.Error = "server_error"
			sendResponse.StatusCode = http.StatusInternalServerError
			sendResponse.SendResponse()
			return
		}

		var setPayment = payment.NewPayment(paymentsAmount, kind.Kind)
		if err := setPayment.AddPayment(); err != nil {
			seterror.SetAppError("setPayment.AddPayment()", err)
			sendResponse.Error = "server_error"
			sendResponse.StatusCode = http.StatusInternalServerError
			sendResponse.SendResponse()
			return
		}
	})
}

func GetPayment() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		var sendResponse = &response.Response{
			Context:    context,
			StatusCode: 200,
		}

		var id string = context.Param("id")

		if id == "" {
			sendResponse.StatusCode = 400
			sendResponse.Error = "miss_payment_id"

			sendResponse.SendResponse()
			return
		}

		paymentItems, err := payment.ReadLocalDataPayments()
		if err != nil {
			seterror.SetAppError("readLocalDataPayments()", err)
			sendResponse.Error = "server_error"
			sendResponse.StatusCode = 500
			sendResponse.SendResponse()
			return
		}

		paymentItem, found := findPaymentByID(paymentItems, id)
		if !found {
			sendResponse.Error = "not_found"
			sendResponse.StatusCode = 404
			sendResponse.SendResponse()
			return
		}

		sendResponse.Result = paymentItem
		sendResponse.SendResponse()
		return
	})
}
