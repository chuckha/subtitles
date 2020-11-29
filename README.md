# Subtitles

A CLI to get plain text out of subtitles

## Supported formats

* `.srt`
* `.ass`

## Usage

Use the binaries found on the release page. Look for the [latest release](https://github.com/chuckha/subtitles/releases/latest)!

```
subex file1.srt file2.ass file3.ass

# output files:
# file1.txt
# file2.txt
# file3.txt
```

# Release

Run goreleaser with a valid environment.

```
GITHUB_TOKEN=<github token> goreleaser --rm-dist
```
