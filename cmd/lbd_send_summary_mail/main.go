package main

import (
	"context"
	"log"
	"net/http"
	"storichallenge_layer/services"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Initialize the account service
	accountService, err := services.NewAccountService()
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// Initialize email builder
	emailBuilder := services.NewEmailBuilder(accountService)

	accountNumber := request.QueryStringParameters["accountNumber"]

	if accountNumber == "" {
		log.Println("Account number is missing from query parameters")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Account number is required",
		}, nil
	}

	monthsParam := request.QueryStringParameters["months"]
	var months []string
	if monthsParam == "" {
		// Default months if no months are provided
		log.Println("Month is missing from query parameters")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Month is required",
		}, nil
	} else {
		// Split months parameter by commas (e.g., "2024/01,2024/02,2024/03")
		months = strings.Split(monthsParam, ",")
	}

	err = emailBuilder.SendAccountSummaryEmail(accountNumber, months)

	if err != nil {
		log.Printf("Failed to send account summary email: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to send account summary email",
		}, nil
	}

	// Return successfull response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Mail was successfully sent.",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
