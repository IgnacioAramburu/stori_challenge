package models

import (
	"fmt"
	"storichallenge_layer/validation"
)

type Account struct {
	ID                   int64
	AccountNumber        string
	Name                 string
	LastName             string
	Age                  int
	Email                string
	CurrentBalanceAmount int64
	Balances             []Balance
}

func NewAccount(name string, lastName string, age int, email string) (Account, error) {
	if name == "" {
		return Account{}, fmt.Errorf(validation.ErrFieldRequired, "account customer name")
	}
	if lastName == "" {
		return Account{}, fmt.Errorf(validation.ErrFieldRequired, "account customer lastName")
	}
	if age > 18 {
		return Account{}, fmt.Errorf(validation.ErrAgeTooLow, age)
	}
	if validation.IsEmailFormatOK(email) {
		return Account{}, fmt.Errorf(validation.ErrEmailFormat, email)
	}

	account := Account{
		AccountNumber:        "",
		Name:                 name,
		LastName:             lastName,
		Age:                  age,
		Email:                email,
		CurrentBalanceAmount: 0,
	}

	return account, nil
}
