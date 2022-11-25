package payment

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"service-task/pkg/seterror"
	"time"
)

type Payment struct {
	Id        string `json:"id"`
	Amount    int64  `json:"amount"`
	Kind      string `json:"kind"`
	CreatedAt int64  `json:"createdAt"`
}

type PaymentsKind struct {
	Designer   int64 `json:"designer"`
	Accountant int64 `json:"accountant"`
	Editor     int64 `json:"editor"`
}

type Kind struct {
	Kind string `json:"kind"`
}

func createDate() int64 {
	var dateTime time.Time = time.Now()
	return dateTime.UnixNano() / int64(time.Millisecond)
}

func generateId() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-")
	var characters = 30

	b := make([]rune, characters)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			seterror.SetAppError("rand.Int", err)
			return ""
		}

		b[i] = letterRunes[n.Int64()]
	}
	return string(b)
}

func ReadLocalDataPayments() ([]Payment, error) {
	var paymentItems []Payment

	file, err := os.ReadFile("./localData/payments.json")
	if err != nil {
		seterror.SetAppError("os.ReadFile", err)
		return nil, err
	}

	if err := json.Unmarshal(file, &paymentItems); err != nil {
		seterror.SetAppError("json.Unmarshal", err)
		return nil, err
	}

	return paymentItems, nil
}

func NewPayment(amount int64, kind string) *Payment {
	return &Payment{
		Id:        generateId(),
		Amount:    amount,
		Kind:      kind,
		CreatedAt: createDate(),
	}
}

func (p *Payment) AddPayment() error {
	var payment = Payment{
		Id:        p.Id,
		Amount:    p.Amount,
		Kind:      p.Kind,
		CreatedAt: p.CreatedAt,
	}

	paymentsItems, err := ReadLocalDataPayments()
	if err != nil {
		seterror.SetAppError("ReadLocalDataPayments()", err)
		return err
	}

	paymentsItems = append(paymentsItems, payment)

	jsonContent, err := json.Marshal(paymentsItems)
	if err != nil {
		seterror.SetAppError("json.Marshal", err)
		return err
	}

	file, err := os.OpenFile("localData/payments.json", os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		seterror.SetAppError("os.OpenFile", err)
		return err
	}

	if _, err := file.Write(jsonContent); err != nil {
		seterror.SetAppError("file.Write", err)
		return err
	}

	return nil
}
