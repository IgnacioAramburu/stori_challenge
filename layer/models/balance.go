package models

import (
	"fmt"
	"storichallenge_layer/utils"
	"storichallenge_layer/validation"
	"time"
)

type Balance struct {
	Month        string
	Amount       int64
	Transactions []Transaction
	AccountID    int64
}

func NewBalance(accountID int64, amount int64, month string) (Balance, error) {

	if accountID == 0 {
		return Balance{}, fmt.Errorf(validation.ErrFieldRequired, "Account ID")
	}

	if month == "" {
		month = utils.GetMonth(time.Now())
	}

	balance := Balance{
		Month:     month,
		Amount:    amount,
		AccountID: accountID,
	}

	return balance, nil
}
