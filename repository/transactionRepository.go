package repository

import "miniWalletExercise/model"

type TransactionRepository struct {
	repository []model.Transaction
}

type TransactionInterface interface {
	ReferenceExist(referenceId string) bool
	Insert(transaction model.Transaction)
}

func CustomTransactionRepository(mockStorage []model.Transaction) TransactionInterface {
	return &TransactionRepository{repository: mockStorage}
}

func NewTransactionRepository() TransactionInterface {
	repo := CustomTransactionRepository(make([]model.Transaction, 0))
	return repo
}

func (tr TransactionRepository) ReferenceExist(referenceId string) bool {
	for _, transaction := range tr.repository {
		if transaction.ReferenceId == referenceId {
			return true
		}
	}
	return false
}

func (tr *TransactionRepository) Insert(transaction model.Transaction) {
	tr.repository = append(tr.repository, transaction)
}
