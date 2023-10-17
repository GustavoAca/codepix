package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"notnull"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same ")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewTransactional(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transactional := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	transactional.ID = uuid.NewV4().String()
	transactional.CreateAt = time.Now()

	err := transactional.isValid()
	if err != nil {
		return nil, err
	}

	return &transactional, nil
}

func (t *Transaction) Completed() error {
	t.Status = TransactionCompleted
	t.UpdateAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.UpdateAt = time.Now()
	t.Description = description
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdateAt = time.Now()
	err := t.isValid()
	return err
}
