package controller

import (
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"miniWalletExercise/model"
	"miniWalletExercise/repository"
	"miniWalletExercise/util/token"
	"net/http"
	"strconv"
	"time"
)

const (
	transactionDeposit    = 0
	transactionWithdrawal = 1
)

type RepositoryHandler struct {
	Wallet      repository.WalletInterface
	Transaction repository.TransactionInterface
}

type InitializeAccountInput struct {
	CustomerXID string `json:"customer_xid"`
}

type InitializeAccountOutput struct {
	Token string `json:"token"`
}

func randomId() string {
	var stdChars = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	s1 := uniuri.NewLenChars(8, stdChars)
	s2 := uniuri.NewLenChars(4, stdChars)
	s3 := uniuri.NewLenChars(4, stdChars)
	s4 := uniuri.NewLenChars(4, stdChars)
	s5 := uniuri.NewLenChars(12, stdChars)
	return fmt.Sprintf("%s-%s-%s-%s-%s", s1, s2, s3, s4, s5)
}

func (rh RepositoryHandler) InitializeAccount(ctx *gin.Context) {
	input := ctx.PostForm("customer_xid")
	walletId := randomId()
	newWallet := model.Wallet{
		Id:        walletId,
		OwnedBy:   input,
		Status:    false,
		EnabledAt: time.UnixMilli(0),
		Balance:   0,
	}
	rh.Wallet.InitializeAccount(newWallet)
	generatedToken, _ := token.GenerateToken(walletId)
	output := InitializeAccountOutput{Token: generatedToken}
	ctx.JSON(http.StatusCreated, gin.H{"data": output, "status": "success"})
}

type WalletOutput struct {
	Id        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   int       `json:"balance"`
}

func (rh RepositoryHandler) EnableWallet(ctx *gin.Context) {
	decoded, _ := token.ExtractWalletIdFromToken(ctx)
	wallet := rh.Wallet.EnableWallet(decoded)
	if wallet.Id == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed"})
		return
	}
	output := WalletOutput{
		Id:        decoded,
		OwnedBy:   wallet.OwnedBy,
		Status:    "enabled",
		EnabledAt: time.Now(),
		Balance:   wallet.Balance,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": map[string]WalletOutput{"wallet": output}})
}

func (rh RepositoryHandler) ViewWallet(ctx *gin.Context) {
	decoded, _ := token.ExtractWalletIdFromToken(ctx)
	wallet := rh.Wallet.ViewWallet(decoded)
	if wallet.Id == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "wallet not found"})
		return
	}
	status := "disabled"
	if wallet.Status == true {
		status = "enabled"
	}
	output := WalletOutput{
		Id:        decoded,
		OwnedBy:   wallet.OwnedBy,
		Status:    status,
		EnabledAt: wallet.EnabledAt,
		Balance:   wallet.Balance,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": map[string]WalletOutput{"wallet": output}})
}

type DepositWalletOutput struct {
	Id          string    `json:"id"`
	DepositedBy string    `json:"deposited_by"`
	Status      string    `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      int       `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

func (rh RepositoryHandler) DepositWallet(ctx *gin.Context) {
	decoded, _ := token.ExtractWalletIdFromToken(ctx)
	amount, _ := strconv.Atoi(ctx.PostForm("amount"))
	referenceId := ctx.PostForm("reference_id")
	wallet := rh.Wallet.ViewWallet(decoded)
	if wallet.Id == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "wallet not found"})
		return
	}
	isEnabled := wallet.Status
	isExist := rh.Transaction.ReferenceExist(referenceId)
	if isExist == true || !isEnabled {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "wallet is not enabled or reference_id already used"})
		return
	}
	transactionId := randomId()
	rh.Transaction.Insert(model.Transaction{
		Id:              transactionId,
		TransactionType: transactionDeposit,
		TransactionBy:   decoded,
		Status:          true,
		TransactionAt:   time.Now(),
		Amount:          amount,
		ReferenceId:     referenceId,
	})
	rh.Wallet.UpdateBalance(decoded, amount)
	deposit := DepositWalletOutput{
		Id:          transactionId,
		DepositedBy: decoded,
		Status:      "success",
		DepositedAt: time.Now(),
		Amount:      amount,
		ReferenceId: referenceId,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": map[string]DepositWalletOutput{"deposit": deposit}})
}

type WithdrawalWalletOutput struct {
	Id          string    `json:"id"`
	WithdrawnBy string    `json:"withdrawn_by"`
	Status      string    `json:"status"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
	Amount      int       `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

func (rh RepositoryHandler) WithdrawalWallet(ctx *gin.Context) {
	decoded, _ := token.ExtractWalletIdFromToken(ctx)
	amount, _ := strconv.Atoi(ctx.PostForm("amount"))
	referenceId := ctx.PostForm("reference_id")
	wallet := rh.Wallet.ViewWallet(decoded)
	if wallet.Id == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "wallet not found"})
		return
	}
	if rh.Transaction.ReferenceExist(referenceId) || !wallet.Status {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "either wallet is not found, not enabled, or reference_id already used"})
		return
	}
	if wallet.Balance < amount {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": "failed", "message": "not enough balance"})
		return
	}
	transactionId := randomId()
	rh.Transaction.Insert(model.Transaction{
		Id:              transactionId,
		TransactionType: transactionWithdrawal,
		TransactionBy:   decoded,
		Status:          true,
		TransactionAt:   time.Now(),
		Amount:          amount,
		ReferenceId:     referenceId,
	})
	withdrawal := WithdrawalWalletOutput{
		Id:          transactionId,
		WithdrawnBy: decoded,
		Status:      "success",
		WithdrawnAt: time.Now(),
		Amount:      amount,
		ReferenceId: referenceId,
	}
	rh.Wallet.UpdateBalance(decoded, -amount)
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": map[string]WithdrawalWalletOutput{"withdrawal": withdrawal}})
}

type DisableWalletOutput struct {
	Id         string    `json:"id"`
	OwnedBy    string    `json:"owned_by"`
	Status     string    `json:"status"`
	DisabledAt time.Time `json:"disabled_at"`
	Balance    int       `json:"balance"`
}

func (rh RepositoryHandler) DisableWallet(ctx *gin.Context) {
	decoded, _ := token.ExtractWalletIdFromToken(ctx)
	isDisabled, err := strconv.ParseBool(ctx.PostForm("is_disabled"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "parse failed"})
		return
	}
	wallet := rh.Wallet.ViewWallet(decoded)
	if isDisabled != true {
		return
	}
	rh.Wallet.DisableWallet(decoded)
	output := DisableWalletOutput{
		Id:         decoded,
		OwnedBy:    wallet.OwnedBy,
		Status:     "disabled",
		DisabledAt: time.Now(),
		Balance:    wallet.Balance,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": map[string]DisableWalletOutput{"wallet": output}})
}
