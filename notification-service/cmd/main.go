package main

import (
	"fmt"

	"github.com/Gergenus/commerce/notification-service/internal/config"
	"github.com/Gergenus/commerce/notification-service/internal/email"
	"github.com/Gergenus/commerce/notification-service/pkg/logger"
)

func main() {
	cfg := config.InitConfig()
	logger := logger.SetUp(cfg.LogLevel)
	mail := email.NewEmailSender(cfg.FromEmail, cfg.FromEmailPassword, cfg.FromEmailSMTP, cfg.SMTPAddr, logger)

	err := mail.SendVerificationEmail("khakimoff.dima@mail.ru", "registration", "verification_email", map[string]string{"VerificationLink": "http://localhost:8081/verification=1488gfdnjlgkfndhbnb"})

	fmt.Println(err)

}
