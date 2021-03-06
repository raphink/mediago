package main

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func SMTPAlert(cfg *config, a *account) (err error) {
	log.Debugf("Sending SMTP report using %s\n", cfg.Smtp.Username)
	auth := smtp.PlainAuth("",
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Hostname,
	)
	msg := fmt.Sprintf("To: %s\r\n", strings.Join(cfg.Smtp.Recipients, ","))
	msg += fmt.Sprintf("Subject: Mediathèque books for %s\r\n\r\n", a.Name)
	msg += a.alerts(false, false)
	if cfg.Gist.GistID != "" {
		msg += fmt.Sprintf("\nSee details on https://gist.github.com/%s\n", cfg.Gist.GistID)
	}
	err = smtp.SendMail(cfg.Smtp.Hostname+":"+strconv.Itoa(cfg.Smtp.Port),
		auth,
		cfg.Smtp.Username,
		cfg.Smtp.Recipients,
		[]byte(msg),
	)
	return
}
