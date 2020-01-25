package extractors

import (
	"fmt"
	"io"

	"github.com/chuckha/subtitles/discovery"
	"github.com/chuckha/subtitles/extractors/ass"
	"github.com/chuckha/subtitles/extractors/srt"
	"github.com/chuckha/subtitles/types"
)

type Extractor interface {
	Extract(io.Reader) (types.Subtitles, error)
}

func ExtracatorFactory(format string) (Extractor, error) {
	switch format {
	case discovery.Srt:
		return &srt.Extractor{}, nil
	case discovery.Ass:
		return &ass.Extractor{}, nil
	default:
		return nil, fmt.Errorf("error: unknown format: %q", format)
	}
}
