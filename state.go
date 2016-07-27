package main

type state struct {
	Message string
}

// Item states
var OK = "OK"
var NeedsRenewing = "NEEDS RENEWING"
var Late = "!!LATE!!"

var stateOK = state{Message: OK}
var stateNeedsRenewing = state{Message: NeedsRenewing}
var stateLate = state{Message: Late}

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
	}
	return ""
}
