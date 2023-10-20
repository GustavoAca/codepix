package factory

import (
	"github.com/codepix/imersao/codepix-go/application/usecase"
	"github.com/codepix/imersao/codepix-go/infra/repository"
	"github.com/jinzhu/gorm"
)

func TransactionUSeCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}

	transactionUseCase := usecase.TransactionUseCase{
		TransactionRepository: &transactionRepository,
		PixRepository:         pixRepository,
	}

	return transactionUseCase

}
