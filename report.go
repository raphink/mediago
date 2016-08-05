package main

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func gistReport(cfg *config, a *account) {
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
	content := a.alerts(false, true)
	if content == "" {
		content = fmt.Sprintf("**No books for %s**", a.Name)
	}
	file.Content = &content
	gist.Files[github.GistFilename(*file.Filename)] = *file

	_, _, err = client.Gists.Edit(cfg.Gist.GistID, gist)
	if err != nil {
		log.Printf("ERROR: failed to save Gist %s: %v\n", cfg.Gist.GistID, err)
	}
}
