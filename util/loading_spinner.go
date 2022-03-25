package util

import (
	"github.com/briandowns/spinner"
	"time"
)

type LoadingSpinner interface {
	Start(text string)
	Stop()
}

type loadingSpinner struct {
	spinner *spinner.Spinner
}

func NewLoadingSpinner() LoadingSpinner {
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	_ = s.Color("green")
	return &loadingSpinner{spinner: s}
}

func (s *loadingSpinner) Start(text string) {
	s.spinner.Suffix = text
	s.spinner.Start()
}

func (s *loadingSpinner) Stop() {
	s.spinner.Stop()
}
