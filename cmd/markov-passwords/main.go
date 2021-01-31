package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/bujimuji/markov-passwords/pkg/markov"
	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
)

func loadPasswords(filePath string) <-chan string {
	// load training passwords
	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("Failed to load training passwords")
	}

	scanner := bufio.NewScanner(file)
	out := make(chan string)
	bar := pb.StartNew(-1)
	go func() {
		for scanner.Scan() {
			bar.Increment()
			out <- scanner.Text() + "\n"
		}
		close(out)
		bar.Finish()
		file.Close()
	}()
	return out
}

func main() {
	// init flags
	var inFilePath = flag.String("infile", "",
		"read training password list from file (- for stdin)")
	var ngram = flag.Int("ngram", 3,
		"number of ngrams to train markov model")
	var chans = flag.Int("j", 64,
		"number of channels for fan-out")
	var train = flag.String("train", "",
		"model path for saving trained model in json format and exit")
	var sample = flag.String("sample", "",
		"model path for generating new passwords")
	var maxAttempts = flag.Uint64("m", 1e6,
		"the amount of passwords created")
	flag.Parse()
	log.SetOutput(os.Stderr)

	// validate flags
	if *chans < 1 {
		log.WithFields(log.Fields{
			"message": strconv.Itoa(*chans),
		}).Fatal("Cannot run with this number of go routines")
	}

	if *train != "" && *sample != "" {
		log.Fatal(`Do not specify train and sample flags if you want to train a	model based on the 
		given password list following by generating new passwords.`)
	}

	c := markov.NewChain(*ngram)

	if (*train == "" && *sample == "") || *train != "" {
		passwords := loadPasswords(*inFilePath)
		pairs := make([]<-chan markov.Pair, *chans)
		for i := 0; i < *chans; i++ {
			pairs[i] = c.Pairs(passwords)
		}
		c.Add(pairs...)
	}

	if *train != "" {
		// open train output file
		var outFile *os.File
		outFile, err := os.Create(*train)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("Create file failed")
		}

		// save training data
		e := json.NewEncoder(outFile)
		err = e.Encode(c.Trans)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("Failed to save training data")
		}
		return
	}

	if *sample != "" {
		model, err := os.Open(*sample)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("Failed to open training data file")
		}
		d := json.NewDecoder(model)
		err = d.Decode(&c.Trans)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("Failed to decode training data")
		}
	}

	if *train == "" && *sample == "" || *sample != "" {
		g := markov.NewGenerator(c)
		g.Generate(*chans, *maxAttempts)
	}
}
