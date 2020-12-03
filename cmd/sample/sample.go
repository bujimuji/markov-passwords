package main

import (
	"encoding/json"
	"errors"
	"flag"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func genChar(stats map[string]map[string]float64, src rand.Source, ngram string) (string, error) {
	dic, found := stats[ngram]
	if found {
		var p = float64(src.Int63()&0xff) / 255.0
		t := 0.0
		var lastKey string
		for key, value := range dic {
			if p < t+value {
				return key, nil
			} else {
				t += value
			}
			lastKey = key
		}
		return lastKey, nil
	} else {
		if len(ngram) > 0 {
			return genChar(stats, src, ngram[0:len(ngram)-1])
		} else {
			return "", errors.New("not found")
		}
	}
}

func genPassword(stats map[string]map[string]float64, src rand.Source, n int, blockSize int) []byte {
	output := strings.Repeat("`", n)
	outBytes := make([]byte, blockSize)
	k := 0
	for i := 0; i < 100; i++ {
		candid, err := genChar(stats, src, output[i:i+n])
		if err != nil {
			output = output[:n]
			continue
		}
		output += candid
		if output[len(output)-1] == '\n' {
			if len(output) > blockSize-k {
				return outBytes[:k]
			}
			k += copy(outBytes[k:], output[n:])
			output = output[:n]
			i = -1
		}
	}
	return []byte("")
}

func main() {
	// init flags
	var inFilePath = flag.String("infile", "", "read training data from file")
	var maxNgram = flag.Int("max-ngram", 3, "maximum number of ngrams to train")
	var blockSize = flag.Int("b", 3000, "block size of prints")
	var numThreads = flag.Int("j", 32, "number of threads")
	flag.Parse()

	// initialize
	rand.Seed(time.Now().Unix())
	log.SetOutput(os.Stderr)

	// load and decode training data
	inFile, err := os.Open(*inFilePath)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Failed to open training data file")
	}
	stats := make(map[string]map[string]float64)
	e := json.NewDecoder(inFile)
	err = e.Decode(&stats)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Failed to decode training data")
	}

	// spawn workers
	var tasks = make(chan int, *numThreads)
	var wg sync.WaitGroup
	wg.Add(*numThreads)

	for i := 0; i < *numThreads; i++ {
		go func() {
			defer wg.Done()
			for {
				src := rand.NewSource(time.Now().Unix())
				pwl := genPassword(stats, src, *maxNgram, *blockSize)
				_, _ = os.Stdout.Write(pwl)
			}
		}()
	}
	close(tasks)
	wg.Wait()
}
