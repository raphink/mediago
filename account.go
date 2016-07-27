package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type account struct {
	Name     string
	Login    string
	Password string
	Items    []*Item
	Alert    bool
}

func (a *account) alerts(colored bool) (alerts string) {
	var state string
	for _, i := range a.Items {
		if colored {
			state = i.State.ColoredString()
		} else {
			state = i.State.String()
		}
		alerts += fmt.Sprintf("[%s]\t%s\t%s\n", state, i.Date.Format("02/01/2006"), i.Title)
	}
	return
}

func (a *account) report(cfg *config) {
	titleColor.Println(a.Name)
	fmt.Println(a.alerts(true))

	if a.Alert && cfg.Report == "smtp" {
		fmt.Printf("Sending SMTP report using %s@%s\n", cfg.Smtp.Username, cfg.Smtp.Hostname)
		auth := smtp.PlainAuth("",
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Hostname,
		)
		msg := fmt.Sprintf("To: %s\r\n", strings.Join(cfg.Smtp.Recipients, ","))
		msg += fmt.Sprintf("Subject: Mediathèque books for %s\r\n\r\n", a.Name)
		msg += a.alerts(false)
		err := smtp.SendMail(cfg.Smtp.Hostname+":"+strconv.Itoa(cfg.Smtp.Port),
			auth,
			cfg.Smtp.Username,
			cfg.Smtp.Recipients,
			[]byte(msg),
		)
		if err != nil {
			log.Print("ERROR: attempting to send a mail ", err)
		}
	}
}
