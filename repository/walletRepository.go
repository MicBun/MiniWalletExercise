package repository

import (
	"miniWalletExercise/model"
	"sync"
	"time"
)

type WalletRepository struct {
	Mu         sync.Mutex
	Repository []model.Wallet
}

type WalletInterface interface {
	InitializeAccount(wallet model.Wallet)
	EnableWallet(id string) model.Wallet
	ViewWallet(walletId string) model.Wallet
	UpdateBalance(walletId string, balance int)
	DisableWallet(id string) model.Wallet
}

func CustomWalletRepository(wallet []model.Wallet) WalletInterface {
	return &WalletRepository{Repository: wallet}
}

func NewWalletRepository() WalletInterface {
	repo := CustomWalletRepository(make([]model.Wallet, 0))
	return repo
}

func (wr *WalletRepository) InitializeAccount(wallet model.Wallet) {
	wr.Repository = append(wr.Repository, wallet)
}

func (wr *WalletRepository) EnableWallet(id string) model.Wallet {
	for i, m := range wr.Repository {
		if m.Id == id {
			wr.Repository[i].Status = true
			wr.Repository[i].EnabledAt = time.Now()
			return m
		}
	}
	return model.Wallet{}
}

func (wr *WalletRepository) ViewWallet(walletId string) model.Wallet {
	for _, wallet := range wr.Repository {
		if wallet.Id == walletId {
			return wallet
		}
	}
	return model.Wallet{}
}

func (wr *WalletRepository) UpdateBalance(walletId string, balance int) {
	wr.Mu.Lock()
	defer wr.Mu.Unlock()
	for i, wallet := range wr.Repository {
		if wallet.Id == walletId {
			wr.Repository[i].Balance += balance
		}
	}
}

func (wr *WalletRepository) DisableWallet(id string) model.Wallet {
	for i, m := range wr.Repository {
		if m.Id == id {
			wr.Repository[i].Status = false
			wr.Repository[i].DisabledAt = time.Now()
			return m
		}
	}
	return model.Wallet{}
}
