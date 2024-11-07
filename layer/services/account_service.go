package services

import (
	"storichallenge_layer/config"
	"storichallenge_layer/models"
	"storichallenge_layer/repository"
)

type AccountService struct {
	AccountRepo     *repository.AccountRepository
	BalanceRepo     *repository.BalanceRepository
	TransactionRepo *repository.TransactionRepository
}

func NewAccountService() (*AccountService, error) {
	db, err := config.ConnectToDB()
	if err != nil {
		return nil, err
	}
	accountRepo := &repository.AccountRepository{DB: db}
	balanceRepo := &repository.BalanceRepository{DB: db}
	transactionRepo := &repository.TransactionRepository{DB: db}

	return &AccountService{
		AccountRepo:     accountRepo,
		BalanceRepo:     balanceRepo,
		TransactionRepo: transactionRepo,
	}, nil
}

func (svc *AccountService) CreateAccount(account models.Account) (int64, error) {
	accountID, err := svc.AccountRepo.Create(account)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}

func (svc *AccountService) GetAccountByAccountNumber(accountNumber string, includeBalances, includeTransactions bool) (models.Account, error) {
	account, err := svc.AccountRepo.GetByAccountNumber(accountNumber, includeBalances, includeTransactions)
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (svc *AccountService) CreateBalance(balance models.Balance) error {
	err := svc.BalanceRepo.Create(balance)
	if err != nil {
		return err
	}
	return nil
}

func (svc *AccountService) CreateTransaction(transaction models.Transaction) error {
	err := svc.TransactionRepo.Create(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (svc *AccountService) GetNumberOfTransactions(accountID int64, month string) (int64, error) {
	transactionNum, err := svc.TransactionRepo.GetNumberOfTransactions(accountID, month)
	if err != nil {
		return -1, err
	}
	return transactionNum, err
}

func (svc *AccountService) GetAverageDebitAmount(accountID int64, month string) (float64, error) {
	avgDebitAmount, err := svc.TransactionRepo.GetAverageDebitAmount(accountID, month)
	if err != nil {
		return -1, err
	}
	return avgDebitAmount, nil
}

func (svc *AccountService) GetAverageCreditAmount(accountID int64, month string) (float64, error) {
	avgCreditAmount, err := svc.TransactionRepo.GetAverageCreditAmount(accountID, month)
	if err != nil {
		return -1, err
	}
	return avgCreditAmount, nil
}
