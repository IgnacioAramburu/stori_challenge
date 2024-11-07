package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"storichallenge_layer/models"
)

type BalanceRepository struct {
	DB              *sql.DB
	AccountRepo     *AccountRepository
	TransactionRepo *TransactionRepository
}

func (repo *BalanceRepository) Create(balance models.Balance) error {
	query := "INSERT INTO balance (account_id, month, amt) VALUES (?,?,?)"
	_, err := repo.DB.Exec(query, balance.AccountID, balance.Month, balance.Amount)
	if err != nil {
		return fmt.Errorf("error while creating balance: %v", err)
	}

	return nil
}

func (repo *BalanceRepository) GetByAccountID(accountID int64, includeTransactions bool) ([]models.Balance, error) {
	query := "SELECT account_id, month, amt FROM balance WHERE account_id = ? ORDER BY month DESC"
	rows, err := repo.DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var balance models.Balance
		err := rows.Scan(&balance.AccountID, &balance.Month, &balance.Amount)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	if includeTransactions {
		transactions, err := repo.TransactionRepo.GetByAccountID(accountID)
		if err != nil {
			return nil, err
		}
		for i := len(balances) - 1; i > 0; i-- {
			for j := len(transactions) - 1; j > 0; j-- {
				if transactions[j].Month == balances[i].Month {
					balances[i].Transactions = append(balances[i].Transactions, transactions[j])
					transactions = append(transactions[:j], transactions[j+1:]...)
					continue
				}
				break
			}
		}
	}
	return balances, nil
}

func (repo *BalanceRepository) GetByAccountIDMonth(accountID int64, month string, includeTransactions bool) (models.Balance, error) {
	query := "SELECT account_id, month, amt FROM balance WHERE account_id = ? AND month = ?"
	var balance models.Balance
	err := repo.DB.QueryRow(query, accountID).Scan(
		&balance.AccountID, &balance.Month, &balance.Amount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Balance{}, errors.New("balance not found")
		}
		return models.Balance{}, err
	}
	if includeTransactions {
		transactions, err := repo.TransactionRepo.GetByAccountIDMonth(accountID, month)
		if err != nil {
			return models.Balance{}, err
		}
		balance.Transactions = transactions
	}
	return balance, nil
}

func (repo *BalanceRepository) UpdateAmountArithmetically(accountID int64, month string, amountToAdd int64) error {
	query := "UPDATE balance SET amount = amount + ? WHERE account_id = ? AND month = ?"
	result, err := repo.DB.Exec(query, amountToAdd, accountID, month)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		newBalance, err := models.NewBalance(accountID, 0, month)
		if err != nil {
			return err
		}
		repo.Create(newBalance)
		return repo.UpdateAmountArithmetically(accountID, month, amountToAdd)
	}

	err = repo.AccountRepo.UpdateCurrentBalanceAmountArithmetrically(accountID, amountToAdd)

	if err != nil {
		return err
	}

	return nil
}
