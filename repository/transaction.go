package repository

import (
	"context"
	"errors"
	"fp-rpl/entity"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepository interface {
	// db transaction
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	// functional
	CreateNewTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) (entity.Transaction, error)
	GetAllTransactions(ctx context.Context, tx *gorm.DB) ([]entity.Transaction, error)
	GetTransactionByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (transactionR *transactionRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := transactionR.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (transactionR *transactionRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Commit().Error
	if err == nil {
		return err
	}
	return nil
}

func (transactionR *transactionRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.WithContext(ctx).Debug().Rollback()
}

func (transactionR *transactionRepository) CreateNewTransaction(ctx context.Context, tx *gorm.DB, transaction entity.Transaction) (entity.Transaction, error) {
	var err error
	if tx == nil {
		tx = transactionR.db.WithContext(ctx).Debug().Create(&transaction)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&transaction).Error
	}

	if err != nil {
		return entity.Transaction{}, err
	}
	return transaction, nil
}

func (transactionR *transactionRepository) GetAllTransactions(ctx context.Context, tx *gorm.DB) ([]entity.Transaction, error) {
	var err error
	var transactions []entity.Transaction

	if tx == nil {
		tx = transactionR.db.WithContext(ctx).Debug().Preload("Spots").Find(&transactions)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Preload("Spots").Find(&transactions).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return transactions, err
	}
	return transactions, nil
}

func (transactionR *transactionRepository) GetTransactionByID(ctx context.Context, tx *gorm.DB, id uint64) (entity.Transaction, error) {
	var err error
	var transaction entity.Transaction
	if tx == nil {
		tx = transactionR.db.WithContext(ctx).Debug().Where("id = $1", id).Preload("Spots").Take(&transaction)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Preload("Spots").Take(&transaction).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return transaction, err
	}
	return transaction, nil
}