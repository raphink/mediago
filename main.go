package main

import (
	"fmt"
	"os"
)

var confFile = fmt.Sprintf("%s/.mediago.conf", os.Getenv("HOME"))

func main() {
	cfg := loadConfig()
	for _, a := range cfg.Account {
		items := a.getItems()
		for _, i := range items {
			alert := i.processState(cfg.RenewBefore.Duration)
			if i.State == stateNeedsRenewing && cfg.AutoRenew {
				err := i.renew(a.Client)
				if err != nil {
					fmt.Printf("Failed to renew %s: %s", i.Title, err)
				}
			}
			a.Alert = alert || a.Alert
			a.Items = append(a.Items, i)
		}
		a.report(cfg)
	}
}
