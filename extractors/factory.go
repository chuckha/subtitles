package extractors

import (
	"fmt"
	"io"

	"github.com/chuckha/subtitles/discovery"
	"github.com/chuckha/subtitles/extractors/srt"
)

type Extractor interface {
	Extract(io.Reader) string
}

func ExtracatorFactory(format string) (Extractor, error) {
	switch format {
	case discovery.Srt:
		return &srt.Extractor{}, nil
	default:
		return nil, fmt.Errorf("error: unknown format: %q", format)
	}
}
