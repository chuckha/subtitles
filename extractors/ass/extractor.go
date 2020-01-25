package ass

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/chuckha/subtitles/types"
)

type Extractor struct{}

var (
	dialoguePrefix = []byte("Dialogue:")
	formatPrefix   = []byte("Format:")
)

const (
	startFormat = "Start"
	endFormat   = "End"
	textFormat  = "Text"
)

type format struct {
	startPos int
	endPos   int
	textPos  int
}

func newFormat(formatLine []byte) format {
	f := format{}
	assFormat := bytes.TrimSpace(bytes.TrimPrefix(formatLine, formatPrefix))
	formatPositions := bytes.Split(assFormat, []byte(", "))
	for i, formatPosition := range formatPositions {
		switch string(formatPosition) {
		case startFormat:
			f.startPos = i
		case endFormat:
			f.endPos = i
		case textFormat:
			f.textPos = i
		default:
			continue
		}
	}
	return f
}

// Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
// Dialogue: 0,0:01:31.27,0:01:33.79,白熊日文,,0,0,0,,シロクマ君の不眠症

func (e *Extractor) Extract(input io.Reader) (types.Subtitles, error) {
	subs := make(types.Subtitles, 0)
	scanner := bufio.NewScanner(input)
	lineno := 0
	var f format
	for scanner.Scan() {
		lineno++
		line := scanner.Bytes()
		if bytes.HasPrefix(line, formatPrefix) {
			f = newFormat(line)
			continue
		}
		if !bytes.HasPrefix(line, dialoguePrefix) {
			continue
		}
		dialogue := bytes.TrimSpace(bytes.TrimPrefix(line, dialoguePrefix))
		dialoguePositions := bytes.Split(dialogue, []byte(","))
		from, err := duration(dialoguePositions[f.startPos])
		if err != nil {
			fmt.Println("line: ", string(line))
			return subs, fmt.Errorf("error parsing start duration on lineno: %v", lineno)
		}
		to, err := duration(dialoguePositions[f.endPos])
		if err != nil {
			fmt.Println("line: ", string(line))
			return subs, fmt.Errorf("error parsing end duration on lineno: %v", lineno)
		}
		content := dialoguePositions[f.textPos]
		if len(bytes.TrimSpace(content)) == 0 {
			continue
		}
		subs = append(subs, types.Subtitle{
			Number:   lineno,
			From:     from,
			To:       to,
			Contents: []string{string(content)},
		})
	}
	return subs, nil
}

func readDuration(b []byte) (time.Duration, error) {
	dur, err := duration(b)
	if err != nil {
		return 0, errors.New("error parsing first duration")
	}
	return dur, nil
}

func duration(d []byte) (time.Duration, error) {
	// replace the first : with 'h'
	out := bytes.Replace(d, []byte(":"), []byte("h"), 1)
	// replace the second : with 'm'
	out = bytes.Replace(out, []byte(":"), []byte("m"), 1)
	// replace the next ,  with 's'
	out = bytes.Replace(out, []byte(","), []byte("s"), 1)
	// add an ms to the end of it
	out = append(out, []byte("ms")...)
	return time.ParseDuration(string(out))
}
