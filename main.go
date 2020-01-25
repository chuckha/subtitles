package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuckha/subtitles/discovery"
	"github.com/chuckha/subtitles/extractors"
	"github.com/chuckha/subtitles/types"
)

func main() {
	os.Exit(run())
}

func run() int {
	for _, file := range os.Args[1:] {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(file)
			fmt.Println(err)
			return 1
		}
		ext := filepath.Ext(file)
		newFilename := strings.TrimSuffix(file, ext) + ".txt"
		subtitles, err := convert(contents)
		output := subtitles.String()
		if err != nil {
			fmt.Println(file)
			fmt.Println(err)
			return 1
		}
		if err := ioutil.WriteFile(newFilename, []byte(output), 0644); err != nil {
			fmt.Println(file)
			fmt.Println(err)
			return 1
		}
	}
	return 0
}

func convert(contents []byte) (types.Subtitles, error) {
	disc := &discovery.Discoverer{}
	format := disc.Discover(contents)
	fmt.Println("Found subtitle type: ", format)
	extractor, err := extractors.ExtracatorFactory(format)
	if err != nil {
		return types.Subtitles{}, err
	}
	return extractor.Extract(bytes.NewReader(contents))
}
