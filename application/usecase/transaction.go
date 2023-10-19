package usecase

import (
	"errors"
	"log"

	"github.com/codepix/imersao/codepix-go/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepositoryInterface
	PixRepository         model.PixKeyRepositoryInterface
}

func (t *TransactionUseCase) Register(accountId string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}
	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}
	transactional, err := model.NewTransactional(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Save(transactional)
	if transactional.ID != "" {
		return transactional, nil
	}

	return nil, errors.New("unable to process this transaction")
}

func (t *TransactionUseCase) Confirm(transactionId string) (*model.Transaction, error) {
	transactional, err := t.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transactional not found", transactionId)
		return nil, err
	}

	transactional.Status = model.TransactionConfirmed
	err = t.TransactionRepository.Save(transactional)
	if err != nil {
		return nil, err
	}

	return transactional, nil
}

func (t *TransactionUseCase) Complete(transactionId string) (*model.Transaction, error) {
	transactional, err := t.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transactional not found", transactionId)
		return nil, err
	}

	transactional.Status = model.TransactionCompleted
	err = t.TransactionRepository.Save(transactional)
	if err != nil {
		return nil, err
	}

	return transactional, nil
}

func (t *TransactionUseCase) Error(transactionId string, reason string) (*model.Transaction, error) {
	transactional, err := t.TransactionRepository.Find(transactionId)
	if err != nil {
		return nil, err
	}

	transactional.Status = model.TransactionError
	transactional.CancelDescription = reason
	err = t.TransactionRepository.Save(transactional)
	if err != nil {
		return nil, err
	}

	return transactional, nil
}
