package srt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/chuckha/subtitles/types"
	strip "github.com/grokify/html-strip-tags-go"
)

type Extractor struct{}

// 0 == subtitles + newline found, starting over
// 1 == number found, expecting duration next
// 2 == duration found, expecting subtitles next
func (e *Extractor) Extract(input io.Reader) (types.Subtitles, error) {
	subs := make(types.Subtitles, 0)
	scanner := bufio.NewScanner(input)
	state := 0
	lineno := 0
	var sub types.Subtitle
	for scanner.Scan() {
		lineno++
		line := scanner.Bytes()
		switch state {
		case 0:
			// expect an infinite number of blank lines or a number
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
			num, err := readNumber(bytes.TrimSpace(line))
			if err != nil {
				return subs, fmt.Errorf("error reading number on line: %v", lineno)
			}
			sub.Number = num
			state = 1
		case 1:
			// read duration
			from, to, err := readDuration(bytes.TrimSpace(line))
			if err != nil {
				fmt.Println("line: ", string(line))
				return subs, fmt.Errorf("error reading duration on line: %v", lineno)
			}
			sub.From = from
			sub.To = to
			state = 2
		case 2:
			// read subtitles or a blank line
			// a blank line resets the whole machine
			if len(bytes.TrimSpace(line)) == 0 {
				state = 0
				switch {
				case len(sub.Contents) > 0:
					subs = append(subs, sub)
				case len(sub.Contents) == 0:
				default:
				}
				sub = types.Subtitle{}
				continue
			}
			item := clean(readSubtitle(line))
			if len(item) == 0 {
				continue
			}
			sub.Contents = append(sub.Contents, item)
		}
	}
	return subs, nil
}

func readNumber(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}

// 00:02:07,840 --> 00:02:09,650
func readDuration(b []byte) (time.Duration, time.Duration, error) {
	durations := bytes.Split(bytes.TrimSpace(b), []byte(" --> "))
	if len(durations) != 2 {
		return 0, 0, errors.New("error parsing duration")
	}
	from, err := duration(durations[0])
	if err != nil {
		return 0, 0, errors.New("error parsing first duration")
	}
	to, err := duration(durations[1])
	if err != nil {
		return 0, 0, errors.New("error parsing second duration")
	}
	return from, to, nil
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

func readSubtitle(b []byte) string {
	return strip.StripTags(string(b))
}

func clean(sub string) string {
	return strings.TrimSpace(strings.ReplaceAll(sub, "{\\an8}", ""))
}
