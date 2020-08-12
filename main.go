package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode"
)

var (
	seed          int64
	size          int
	inputPath     string
	outputPath    string
)

func init() {
	flag.Int64Var(&seed, "s", 0, "The seed for shuffle, 0 represents a time related seed")
	flag.IntVar(&size, "z", 3, "The maximum length of characters to shuffle at once")
	flag.StringVar(&inputPath, "i", "", "Input file path, empty represents stdin")
	flag.StringVar(&outputPath, "o", "", "Input file path, empty represents stdout")
}

func newRand() *rand.Rand {
	var r *rand.Rand
	if seed == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(seed))
	}
	return r
}

func openFile(path string, empty *os.File, flag int, name string) *os.File {
	var f *os.File
	if path == "" {
		f = empty
	} else {
		var err error
		f, err = os.OpenFile(path, flag, 0)
		if err != nil {
			log.Fatal("failed to open " + name + " file")
		}
	}
	return f
}

func readFile(file *os.File) []rune {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("failed to read file")
	}
	return []rune(string(bytes))
}

func writeFile(file *os.File, content []rune) {
	_, err := file.Write([]byte(string(content)))
	if err != nil {
		log.Fatal("failed to write file")
	}
}

func shuffleSlice(content []rune, r *rand.Rand) {
	if len(content) == 0 {
		return
	}
	r.Shuffle(len(content), func(i, j int) {
		content[i], content[j] = content[j], content[i]
	})
}

func structureChar(char rune) bool {
	switch char {
	case '是', '和', '与', '我', '你', '她', '会', '不', '好', '用', '有', '的':
		return true
	}
	return false
}

func shuffle(content []rune, r *rand.Rand) {
	lowerBound := 0
	for upperBound, char := range content {
		if !unicode.Is(unicode.Han, char) || !unicode.IsLetter(char) || structureChar(char) || upperBound - lowerBound >= size  {
			shuffleSlice(content[lowerBound:upperBound],  r)
			lowerBound = upperBound + 1
		}
	}
}

func main() {
	flag.Parse()
	r := newRand()
	inputFile := openFile(inputPath, os.Stdin, os.O_RDONLY, "input")
	outputFile := openFile(outputPath, os.Stdout, os.O_WRONLY, "output")
	content := readFile(inputFile)
	shuffle(content, r)
	writeFile(outputFile, content)
}