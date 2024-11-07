package models

import (
	"fmt"
	"storichallenge_layer/utils"
	"storichallenge_layer/validation"
	"time"
)

type Transaction struct {
	ID        int64
	AccountID int64
	Month     string
	DateTime  time.Time
	Amount    int64
}

func NewTransaction(amount int64, dateTime time.Time, accountID int64) (Transaction, error) {
	if amount == 0 {
		return Transaction{}, fmt.Errorf(validation.ErrFieldRequired, "Transaction Amount")
	}

	if dateTime.IsZero() {
		dateTime = time.Now()
	}

	month := utils.GetMonth(dateTime)

	transaction := Transaction{
		AccountID: accountID,
		Month:     month,
		DateTime:  dateTime,
		Amount:    amount,
	}

	return transaction, nil
}
