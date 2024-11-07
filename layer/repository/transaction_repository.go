package repository

import (
	"database/sql"
	"fmt"
	"math"
	"storichallenge_layer/models"
)

type TransactionRepository struct {
	DB          *sql.DB
	BalanceRepo *BalanceRepository
}

func (repo *TransactionRepository) Create(transaction models.Transaction) error {
	query := "INSERT INTO transaction (account_id, month, dt, amt) VALUES (?,?,?,?)"
	_, err := repo.DB.Exec(query, transaction.AccountID, transaction.Month, transaction.DateTime, transaction.Amount)
	if err != nil {
		return fmt.Errorf("error while creating transaction: %v", err)
	}
	err = repo.BalanceRepo.UpdateAmountArithmetically(transaction.AccountID, transaction.Month, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepository) GetByAccountID(accountID int64) ([]models.Transaction, error) {
	query := "SELECT id, account_id, month, dt, amt FROM transaction WHERE account_id = ? ORDER BY dt DESC"
	rows, err := repo.DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Month, &transaction.DateTime, &transaction.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (repo *TransactionRepository) GetByAccountIDMonth(accountID int64, month string) ([]models.Transaction, error) {
	query := "SELECT id, account_id, month, dt, amt FROM transaction WHERE account_id = ? AND month = ? ORDER BY dt DESC"
	rows, err := repo.DB.Query(query, accountID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Month, &transaction.DateTime, &transaction.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (repo *TransactionRepository) GetNumberOfTransactions(accountID int64, month string) (int64, error) {
	query := " SELECT SUM(amt) FROM transaction WHERE account_id = ? AND month = ? "

	var transactionAmount int64
	err := repo.DB.QueryRow(query, accountID, month).Scan(&transactionAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return transactionAmount, nil
}

func (repo *TransactionRepository) GetAverageDebitAmount(accountID int64, month string) (float64, error) {
	query := "SELECT AVG(amt) FROM transaction WHERE account_id = ? AND month = ? AND amt < 0"

	var avgDebit float64
	err := repo.DB.QueryRow(query, accountID, month).Scan(&avgDebit)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	avgDebit = math.Round(float64(avgDebit)) / 100

	return avgDebit, nil
}

func (repo *TransactionRepository) GetAverageCreditAmount(accountID int64, month string) (float64, error) {
	query := "SELECT AVG(amt) FROM transaction WHERE account_id = ? AND month = ? AND amt > 0"

	var avgCredit float64
	err := repo.DB.QueryRow(query, accountID, month).Scan(&avgCredit)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	avgCredit = math.Round(float64(avgCredit)) / 100

	return avgCredit, nil
}
