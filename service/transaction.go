package service

import (
	"context"
	"fp-rpl/dto"
	"fp-rpl/entity"
	"fp-rpl/repository"

	"github.com/jinzhu/copier"
)

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

type TransactionService interface {
	CreateNewTransaction(ctx context.Context, transactionDTO dto.TransactionMakeRequest) (entity.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]entity.Transaction, error)
	GetTransactionByID(ctx context.Context, id uint64) (entity.Transaction, error)
}

func NewTransactionService(transactionR repository.TransactionRepository) TransactionService {
	return &transactionService{transactionRepository: transactionR}
}

func (transactionS *transactionService) CreateNewTransaction(ctx context.Context, transactionDTO dto.TransactionMakeRequest) (entity.Transaction, error) {
	// Copy TransactionDTO to empty newly created transaction var
	var transaction entity.Transaction
	copier.Copy(&transaction, &transactionDTO)

	// create new transaction
	newTransaction, err := transactionS.transactionRepository.CreateNewTransaction(ctx, nil, transaction)
	if err != nil {
		return entity.Transaction{}, err
	}
	return newTransaction, nil
}

func (transactionS *transactionService) GetAllTransactions(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := transactionS.transactionRepository.GetAllTransactions(ctx, nil)
	if err != nil {
		return []entity.Transaction{}, err
	}
	return transactions, nil
}

func (transactionS *transactionService) GetTransactionByID(ctx context.Context, id uint64) (entity.Transaction, error) {
	transaction, err := transactionS.transactionRepository.GetTransactionByID(ctx, nil, id)
	if err != nil {
		return entity.Transaction{}, err
	}
	return transaction, nil
}