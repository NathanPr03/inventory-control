package pkg

import (
	"github.com/joho/godotenv"
	"net/smtp"
	"os"
)

func SendEmail(subject string, body string) error {
	_ = godotenv.Load("../../.env.development.local")

	emailUsername := os.Getenv("EMAIL_USERNAME")

	auth := smtp.PlainAuth(
		"",
		emailUsername,
		os.Getenv("EMAIL_PASSWORD"),
		"smtp.gmail.com",
	)

	to := []string{"doroyi1377@jzexport.com"}
	msg := []byte("To: doroyi1377@jzexport.com\r\n" +
		"Subject: " + subject +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, emailUsername, to, msg)
	if err != nil {
		println(err)
		return err
	}

	return nil
}
