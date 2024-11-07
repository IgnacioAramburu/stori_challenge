package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"storichallenge_layer/models"
)

type AccountRepository struct {
	DB          *sql.DB
	BalanceRepo *BalanceRepository
}

func (repo *AccountRepository) Create(account models.Account) (int64, error) {
	query := "INSERT INTO account (account_number, name, last_name, age, email, cur_balance_amt) VALUES (?,?,?,?,?,?)"
	result, err := repo.DB.Exec(query, account.AccountNumber, account.Name, account.LastName, account.Age, account.Email, account.CurrentBalanceAmount)
	if err != nil {
		return 0, fmt.Errorf("error while creating account: %v", err)
	}
	accountID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("error occured when getting last inserted account id")
	}

	initBalance, err := models.NewBalance(accountID, 0, "")

	if err != nil {
		return accountID, errors.New("error while generating new balance")
	}

	err = repo.BalanceRepo.Create(initBalance)

	if err != nil {
		return accountID, errors.New("error while inserting in DB initial balance of account")
	}

	return accountID, nil
}

func (repo *AccountRepository) GetByID(id int64, includeBalances, includeTransactions bool) (models.Account, error) {
	query := "SELECT id, account_number, name, last_name, age, email, cur_balance_amt FROM account WHERE id = ?"
	var account models.Account
	err := repo.DB.QueryRow(query, id).Scan(
		&account.ID, &account.AccountNumber, &account.Name, &account.LastName,
		&account.Age, &account.Email, &account.CurrentBalanceAmount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}
	if includeBalances {
		balances, err := repo.BalanceRepo.GetByAccountID(id, includeTransactions)
		if err != nil {
			return models.Account{}, err
		}
		account.Balances = balances
	}
	return account, nil
}

func (repo *AccountRepository) GetByAccountNumber(accountNumber string, includeBalances, includeTransactions bool) (models.Account, error) {
	query := "SELECT id, account_number, name, last_name, age, email, cur_balance_amt FROM account WHERE account_number = ?"
	var account models.Account
	err := repo.DB.QueryRow(query, accountNumber).Scan(
		&account.ID, &account.AccountNumber, &account.Name, &account.LastName,
		&account.Age, &account.Email, &account.CurrentBalanceAmount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}
	if includeBalances {
		balances, err := repo.BalanceRepo.GetByAccountID(account.ID, includeTransactions)
		if err != nil {
			return models.Account{}, err
		}
		account.Balances = balances
	}
	return account, nil
}

func (repo *AccountRepository) GetAll() ([]models.Account, error) {
	query := `SELECT id, account_number, name, last_name, age, email, current_balance_amount 
			  FROM accounts`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID, &account.AccountNumber, &account.Name, &account.LastName,
			&account.Age, &account.Email, &account.CurrentBalanceAmount,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (repo *AccountRepository) UpdateCurrentBalanceAmountArithmetrically(accountID int64, amountToAdd int64) error {
	query := "UPDATE accounts SET cur_balance_amt = cur_balance_amt + ? WHERE id = ?"
	result, err := repo.DB.Exec(query, amountToAdd)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("account current balance could not be updated")
	}

	if rowsAffected == 0 {
		return errors.New("no balance found with the specified account ID and month")
	}

	return nil
}
