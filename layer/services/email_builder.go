package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"storichallenge_layer/config"
	"text/template"
)

type EmailBuilder struct {
	AccountService *AccountService
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	AssetsPath     string
}

func NewEmailBuilder(accountService *AccountService) *EmailBuilder {
	return &EmailBuilder{
		AccountService: accountService,
		SMTPHost:       config.SMTP_HOST,
		SMTPPort:       config.SMTP_PORT,
		SMTPUser:       config.SMTP_USERNAME,
		SMTPPassword:   config.SMTP_PASSWORD,
		AssetsPath:     "../assets",
	}
}

type EmailTemplate struct {
	AccountNumber    string
	CurrentBalance   float64
	TransactionsInfo []TransactionsMonthData
	LogoBase64       string
}

type TransactionsMonthData struct {
	Month     string
	Qty       int64
	AvgDebit  float64
	AvgCredit float64
}

func (e *EmailBuilder) SendAccountSummaryEmail(accountNumber string, months []string) error {
	account, err := e.AccountService.GetAccountByAccountNumber(accountNumber, false, false)

	if err != nil {
		return err
	}

	currentBalance := float64(account.CurrentBalanceAmount) / 100

	var transactionsInfo []TransactionsMonthData

	for _, month := range months {
		transactionNum, err := e.AccountService.GetNumberOfTransactions(account.ID, month)

		if err != nil {
			return err
		}

		avgDebit, err := e.AccountService.GetAverageDebitAmount(account.ID, month)

		if err != nil {
			return nil
		}

		avgCredit, err := e.AccountService.GetAverageCreditAmount(account.ID, month)
		transactionsInfo = append(transactionsInfo, TransactionsMonthData{
			Month:     month,
			Qty:       transactionNum,
			AvgDebit:  avgDebit,
			AvgCredit: avgCredit,
		})
	}

	logoBase64, err := e.encodeImageToBase64("stori_logo.png")

	emailData := EmailTemplate{
		AccountNumber:    accountNumber,
		CurrentBalance:   currentBalance,
		TransactionsInfo: transactionsInfo,
		LogoBase64:       logoBase64,
	}

	body, err := e.buildAccountSummaryEmailBody(emailData)

	if err != nil {
		return err
	}

	return e.sendEmail(account.Email, "Stori: Account Summary", body)

}

func (e *EmailBuilder) encodeImageToBase64(filename string) (string, error) {
	path := filepath.Join(e.AssetsPath, filename)
	imageData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imageData), nil
}

func (e *EmailBuilder) buildAccountSummaryEmailBody(data EmailTemplate) (string, error) {
	tmpl := `<html>
	<head>
		<title>Account Summary</title>
	</head>
	<body>
		<div style="text-align: center;">
			<img src="data:image/png;base64,{{.LogoBase64}}" alt="Company Logo" style="width: 150px; height: auto;">
		</div>
		<h2>Account Summary for {{.AccountNumber}}</h2>
		<p>Total Balance: ${{.TotalBalance}}</p>
		<p>Number of Transactions (last 2 months): {{.NumberOfTransactions}}</p>
		<p>Average Debit Amount: ${{.AverageDebitAmount}}</p>
		<p>Average Credit Amount: ${{.AverageCreditAmount}}</p>
	</body>
	</html>`

	t, err := template.New("emailTemplate").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return buf.String(), nil

}

// sendEmail sends the email using SMTP
func (e *EmailBuilder) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.SMTPUser, e.SMTPPassword, e.SMTPHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s", to, subject, body))

	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	if err := smtp.SendMail(addr, auth, e.SMTPUser, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s successfully", to)
	return nil
}
