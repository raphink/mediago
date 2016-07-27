package main

import (
	"fmt"
	"time"
)

type Item struct {
	Entite   string
	Date     time.Time
	Location string
	Type     string
	Title    string
	Barcode  string
	RentType string
	Booked   string
	State    state
}

func (i *Item) processState(renewBefore time.Duration) (alert bool) {
	now := time.Now()
	renewDate := now.Add(renewBefore)

	if now.After(i.Date) {
		i.State = stateLate
		alert = true
	} else if renewDate.After(i.Date) {
		i.State = stateNeedsRenewing
		alert = true
	} else {
		i.State = stateOK
	}
	return
}

func (i *Item) renew() (err error) {
	fmt.Println("Autorenewing is not implemented yet")
	return
}
