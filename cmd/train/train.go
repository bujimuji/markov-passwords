package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type tuple struct {
	a string
	b string
}

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func loadFile(fileName string) []byte {
	buffBytes1, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Failed to load training file")
	}
	if buffBytes1[len(buffBytes1)-1] != '\n' {
		buffBytes1 = append(buffBytes1, '\n')
	}
	return buffBytes1
}

func countNumLines(fileBuff []byte) int {
	lineSep := []byte{'\n'}
	return bytes.Count(fileBuff, lineSep)
}

func scanLines(fileBytes []byte, lineCount int) []int {
	linePositions := make([]int, lineCount+1)

	linePositions[0] = 0
	var k = 1
	for i := 0; i < len(fileBytes); i++ {
		if fileBytes[i] == '\n' {
			linePositions[k] = i + 1
			k++
		}
	}

	return linePositions
}

func loadAndScanFile(fileName string) ([]byte, int, []int) {
	var fileBytes []byte
	var lineCount int
	var fileLines []int
	fileBytes = loadFile(fileName)
	lineCount = countNumLines(fileBytes)
	fileLines = scanLines(fileBytes, lineCount)
	return fileBytes, lineCount, fileLines
}

func ngram(line string, n int) []tuple {
	output := make([]tuple, 0)
	for i, char := range line {
		if i-n < 0 {
			buff := ""
			for j := 0; j < abs(i-n); j++ {
				buff += "`"
			}
			buff += line[0:i]
			output = append(output, tuple{buff, string(char)})
		} else {
			output = append(output, tuple{line[i-n:i], string(char)})
		}
	}
	return output
}

func freq2prob(stats map[string]map[string]float64) {
	for gram := range stats {
		chars := make([]string, 0)
		occur := make([]float64, 0)
		probs := make([]float64, 0)

		for key, value := range stats[gram] {
			chars = append(chars, key)
			occur = append(occur, value)
		}
		total := 0.0
		for _, v := range occur {
			total += v
		}
		for _, v := range occur {
			probs = append(probs, v/total)
		}
		for key, value := range stats[gram] {
			stats[gram][key] = value/total
		}
	}
}

func train(dictBytes []byte, dictLineCount int, dictLines []int, maxNgram int) map[string]map[string]float64 {
	stats := make(map[string]map[string]float64)
	log.Info("Starting train process ...")
	bar := pb.StartNew(dictLineCount)
	for i := 0; i < dictLineCount; i++ {
		line := string(dictBytes[dictLines[i]:dictLines[i+1]-1])
		for i := 0; i < maxNgram; i++ {
			for _, gram := range ngram(line+"\n", i+1) {
				prev := gram.a
				next := gram.b

				_, found := stats[prev]
				if !found {
					stats[prev] = make(map[string]float64)
				}

				_, found = stats[prev][next]
				if !found {
					stats[prev][next] = 0
				}
				stats[prev][next] += 1
			}
		}
		bar.Increment()
	}
	bar.Finish()
	return stats
}

func main() {
	// init flags
	var outFilePath = flag.String("outfile", "", "write training data to file")
	var inFilePath = flag.String("infile", "", "read training data from file")
	var maxNgram = flag.Int("max-ngram", 3, "maximum number of ngrams to train")
	flag.Parse()
	log.SetOutput(os.Stderr)

	// create and open output file
	outFile, err := os.Create(*outFilePath)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Create file failed")
	}

	// load training data
	dictBytes, dictLineCount, dictLines := loadAndScanFile(*inFilePath)

	// start train process
	stats := train(dictBytes, dictLineCount, dictLines, *maxNgram)

	// converts frequency to probabilities
	freq2prob(stats)

	// save training data
	e := json.NewEncoder(outFile)
	err = e.Encode(stats)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Failed to save training data")
	}
}
