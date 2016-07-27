package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Colors
var titleColor = color.New(color.FgBlue).Add(color.Bold).Add(color.Underline)
var okColor = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
var warnColor = color.New(color.FgYellow).Add(color.Bold).SprintFunc()
var errColor = color.New(color.FgRed).Add(color.Bold).SprintFunc()

// Item states
var OK = "OK"
var NeedsRenewing = "NEEDS RENEWING"
var Late = "!!LATE!!"

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
	fmt.Printf("Autorenewing is not implemented yet")
	return
}
