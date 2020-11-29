package srt

import (
	"bufio"
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

const uefeff = "\uefeff"
const efbbbf = "\xef\xbb\xbf"

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
		line := scanner.Text()
		line = strings.TrimSpace(line)
		switch state {
		case 0:
			// expect an infinite number of blank lines or a number
			if len(line) == 0 {
				continue
			}
			line = strings.Replace(line, uefeff, "", 1)
			line = strings.Replace(line, efbbbf, "", -1)
			num, err := strconv.Atoi(line)
			if err != nil {
				return subs, fmt.Errorf("error reading %q on line number %d", line, lineno)
			}
			sub.Number = num
			state = 1
		case 1:
			// read duration
			from, to, err := readDuration(line)
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
			if len(line) == 0 {
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

// 00:02:07,840 --> 00:02:09,650
func readDuration(b string) (time.Duration, time.Duration, error) {
	durations := strings.Split(b, " --> ")
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

func duration(d string) (time.Duration, error) {
	// replace the first : with 'h'
	out := strings.Replace(d, ":", "h", 1)
	// replace the second : with 'm'
	out = strings.Replace(out, ":", "m", 1)
	// replace the next ,  with 's'
	out = strings.Replace(out, ",", "s", 1)
	// add an ms to the end of it
	out = out + "ms"
	return time.ParseDuration(out)
}

func readSubtitle(b string) string {
	return strip.StripTags(b)
}

func clean(sub string) string {
	return strings.TrimSpace(strings.ReplaceAll(sub, "{\\an8}", ""))
}
