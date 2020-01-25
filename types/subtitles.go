package types

import (
	"strings"
	"time"
)

type Subtitles []Subtitle

func (s Subtitles) String() string {
	subs := make([]string, 0)
	for _, s := range s {
		subs = append(subs, s.String())
	}
	return strings.Join(subs, "\n")
}

type Subtitle struct {
	Number   int
	From, To time.Duration
	Contents []string
}

func (s *Subtitle) String() string {
	return strings.Join(s.Contents, "\n")
}
