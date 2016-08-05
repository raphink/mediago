package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func SMTPAlert(cfg *config, a *account) {
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

func gistAlert(cfg *config, a *account) {
	log.Println("Saving report to Gist")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Gist.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	gist, _, err := client.Gists.Get(cfg.Gist.GistID)
	if err != nil {
		log.Printf("ERROR: could not retrieve Gist %s: %v\n", cfg.Gist.GistID, err)
	}
	var file *github.GistFile
	for _, f := range gist.Files {
		if *f.Filename == fmt.Sprintf("%s.md", a.Name) {
			file = &f
			break
		}
	}
	if file == nil {
		file = new(github.GistFile)
		filename := fmt.Sprintf("%s.md", a.Name)
		file.Filename = &filename
	}
	content := a.alerts(false)
	file.Content = &content
	gist.Files[github.GistFilename(*file.Filename)] = *file

	_, _, err = client.Gists.Edit(cfg.Gist.GistID, gist)
	if err != nil {
		log.Printf("ERROR: failed to save Gist %s: %v\n", cfg.Gist.GistID, err)
	}
}
