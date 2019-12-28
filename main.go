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
		output, err := convert(contents)
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

func convert(contents []byte) (string, error) {
	disc := &discovery.SRTDiscoverer{}
	format := disc.Discover(contents)
	extractor, err := extractors.ExtracatorFactory(format)
	if err != nil {
		return "", err
	}
	return extractor.Extract(bytes.NewReader(contents))
}
