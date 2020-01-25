package discovery

import "bytes"

const (
	Srt = "SRT"
	Ass = "ASS"
)

type Discoverer struct{}

func (d *Discoverer) Discover(input []byte) string {
	if bytes.Contains(input, []byte("aegisub.org")) {
		return Ass
	}
	return Srt
}
