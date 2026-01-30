package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

// Stats содержит статистику файла
type Stats struct {
	Lines int64
	Words int64
	Bytes int64
	Chars int64
}

// Flags содержит флаги командной строки
type Flags struct {
	Lines bool
	Words bool
	Bytes bool
	Chars bool
}

func CountStats(r io.Reader) (Stats, error) {
	res := Stats{}
	reader := bufio.NewReader(r)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return Stats{}, err
		}

		stringData := string(line)

		res.Lines += 1
		res.Chars += int64(utf8.RuneCountInString(stringData))
		res.Bytes += int64(len(line))
		res.Words += int64(len(strings.Fields(stringData)))
	}

	return res, nil
}

func ProcessFile(path string) (Stats, error) {
	f, err := os.Open(path)
	if err != nil {
		return Stats{}, err
	}
	defer f.Close()

	stats, err := CountStats(f)
	if err != nil {
		return Stats{}, err
	}

	return stats, nil
}

func FormatOutput(stats Stats, filename string, flags Flags) string {
	res := ""
	if flags.Lines {
		res += fmt.Sprintf("%7d ", stats.Lines)
	}
	if flags.Words {
		res += fmt.Sprintf("%7d ", stats.Words)
	}
	if flags.Bytes {
		res += fmt.Sprintf("%7d ", stats.Bytes)
	}
	if flags.Chars {
		res += fmt.Sprintf("%7d ", stats.Chars)
	}
	res += fmt.Sprintf("%s\n", filename)
	return res
}

func main() {

	var (
		includeLines      = flag.Bool("l", false, "count lines")
		includeWords      = flag.Bool("w", false, "count words")
		includeBytes      = flag.Bool("c", false, "count bytes")
		includeCharacters = flag.Bool("m", false, "count characters")
	)

	flag.Parse()
	args := flag.Args()

	flags := Flags{
		Lines: *includeLines,
		Words: *includeWords,
		Bytes: *includeBytes,
		Chars: *includeCharacters,
	}

	if flags.Chars && flags.Bytes {
		flags.Chars = false
	}

	if !flags.Lines &&
		!flags.Chars &&
		!flags.Bytes &&
		!flags.Words {
		flags.Lines = true
		flags.Words = true
		flags.Bytes = true
	}

	if len(args) == 0 {
		log.Fatal("empty args")
	}

	total := Stats{
		Lines: 0,
		Words: 0,
		Chars: 0,
		Bytes: 0,
	}

	res := ""

	for _, fileName := range args {
		stats, err := ProcessFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		total.Lines += stats.Lines
		total.Words += stats.Words
		total.Bytes += stats.Bytes
		total.Chars += stats.Chars

		res += FormatOutput(stats, fileName, flags)
	}

	if len(args) > 1 {
		res += FormatOutput(total, "total", flags)
	}
	fmt.Print(res)
}
