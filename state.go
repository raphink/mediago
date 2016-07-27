package main

import "github.com/fatih/color"

type state struct {
	Message string
}

// Item states
var OK = "OK"
var NeedsRenewing = "NEEDS RENEWING"
var Late = "LATE!"
var Renewed = "RENEWED"
var FailedRenewing = "FAILED RENEWING!"

var stateOK = state{Message: OK}
var stateNeedsRenewing = state{Message: NeedsRenewing}
var stateLate = state{Message: Late}
var stateRenewed = state{Message: Renewed}
var stateFailedRenewing = state{Message: FailedRenewing}

// Colors
var titleColor = color.New(color.FgBlue).Add(color.Bold).Add(color.Underline)
var okColor = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
var warnColor = color.New(color.FgYellow).Add(color.Bold).SprintFunc()
var errColor = color.New(color.FgRed).Add(color.Bold).SprintFunc()

func (s *state) String() string {
	return s.Message
}

func (s *state) ColoredString() string {
	switch s.Message {
	case OK:
		return okColor(s.Message)
	case NeedsRenewing:
		return warnColor(s.Message)
	case Late:
		return errColor(s.Message)
	case Renewed:
		return okColor(s.Message)
	case FailedRenewing:
		return errColor(s.Message)
	}
	return ""
}
