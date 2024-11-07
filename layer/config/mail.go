package config

import "os"

var (
	SMTP_HOST     = os.Getenv("SMTP_HOST")
	SMTP_PORT     = os.Getenv("SMTP_PORT")
	SMTP_USERNAME = os.Getenv("SMTP_USERNAME")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
)
