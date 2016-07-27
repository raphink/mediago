package main

import (
	"fmt"
	"os"
)

var confFile = fmt.Sprintf("%s/.mediago.conf", os.Getenv("HOME"))

func main() {
	cfg := loadConfig()
	for _, a := range cfg.Account {
		items := getAccountItems(a.Name, a.Login, a.Password)
		for _, i := range items {
			alert := i.processState(cfg.RenewBefore.Duration)
			if i.State == stateNeedsRenewing && cfg.AutoRenew {
				_ = i.renew()
			}
			a.Alert = alert || a.Alert
			a.Items = append(a.Items, i)
		}
		a.report(cfg)
	}
}
