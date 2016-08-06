package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func gistReport(cfg *config, a *account) (err error) {
	log.Debug("Saving report to Gist")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Gist.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	gist, _, err := client.Gists.Get(cfg.Gist.GistID)
	if err != nil {
		err = fmt.Errorf("could not retrieve Gist %s: %v\n", cfg.Gist.GistID, err)
		return
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
	content := fmt.Sprintf("# %s\n\n", a.Name)
	newContent := a.alerts(false, true)
	if newContent == "" {
		content += fmt.Sprintf("**No books for %s**", a.Name)
	} else {
		content += newContent
	}
	file.Content = &content
	gist.Files[github.GistFilename(*file.Filename)] = *file

	_, _, err = client.Gists.Edit(cfg.Gist.GistID, gist)
	if err != nil {
		err = fmt.Errorf("failed to save Gist %s: %v\n", cfg.Gist.GistID, err)
		return
	}
	return
}
