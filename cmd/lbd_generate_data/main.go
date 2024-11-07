package main

import (
	"context"
	"log"
	"math"
	"math/rand"
	"time"

	"storichallenge_layer/models"
	"storichallenge_layer/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize the account service
	accountService, err := services.NewAccountService()
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	accounts, err := getSampleAccounts()

	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// Create each account in the database
	for i := range accounts {
		accountID, err := accountService.CreateAccount(accounts[i])
		if err != nil {
			log.Fatalf("Failed to save account in DB: %v", err)
		}
		accounts[i].ID = accountID
	}

	startOfYear := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(2024, time.December, 31, 23, 59, 59, 0, time.UTC)

	// Create random transactions for each account in database
	for i := 0; i < 1000; i++ {
		transactionAmount := int64(randomFloat2Decimals(-1000, 1000) * 100)
		transactionDate := randomDateTime(startOfYear, endOfYear)
		anyAccount := accounts[rand.Intn(len(accounts))]

		transaction, err := models.NewTransaction(transactionAmount, transactionDate, anyAccount.ID)

		if err != nil {
			log.Fatalf("Failed to create transaction: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		err = accountService.CreateTransaction(transaction)

		if err != nil {
			log.Fatalf("Failed to save transaction in DB: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}
	}

	// Return successfull response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Accounts and transactions successfully created.",
	}, nil
}

func getSampleAccounts() ([]models.Account, error) {
	account1, err := models.NewAccount("Max", "Verstappen", 28, "ignaciomatiasaramburudeveloper@gmail.com")
	if err != nil {
		return nil, err
	}
	account2, err := models.NewAccount("Charles", "Leclerc", 27, "ignaciomatiasaramburudeveloper@gmail.com")
	if err != nil {
		return nil, err
	}
	account3, err := models.NewAccount("Valtteri", "Bottas", 26, "ignaciomatiasaramburudeveloper@gmail.com")
	if err != nil {
		return nil, err
	}

	// Create some sample accounts
	accounts := []models.Account{account1, account2, account3}

	return accounts, nil
}

// randomFloat2Decimals generates a random float64 between min and max, rounded to 2 decimals
func randomFloat2Decimals(min, max float64) float64 {
	var result float64
	for result == 0 {
		result = min + rand.Float64()*(max-min)
		result = math.Round(result*100) / 100
	}
	return result
}

// randomDateTime generates a random date between min and max
func randomDateTime(min, max time.Time) time.Time {
	diff := max.Unix() - min.Unix()
	secondsOffset := rand.Int63n(diff)
	return min.Add(time.Duration(secondsOffset) * time.Second)
}

func main() {
	lambda.Start(HandleRequest)
}
