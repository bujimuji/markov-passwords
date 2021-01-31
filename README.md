# Markov-chain passwords generator

[![Go Reference](https://pkg.go.dev/badge/github.com/bujimuji/markov-passwords.svg)](https://pkg.go.dev/github.com/bujimuji/markov-passwords)
[![Go Report Card](https://goreportcard.com/badge/github.com/bujimuji/markov-passwords?style=flat-square)](https://goreportcard.com/report/github.com/bujimuji/markov-passwords)
![GitHub](https://img.shields.io/github/license/bujimuji/markov-passwords)

The Markov-chain password generator is one of the password guessing methods. The trained model represents the probability of transition between specific characters.
Password guessing attacks using n-grams (i.e., substrings of length n appearing in a training set) have been originally proposed by [Narayanan and Shmatikov](https://www.cs.cornell.edu/~shmat/shmat_ccs05pwd.pdf).

There is another implementation by *brannondorsey* [here](https://github.com/brannondorsey/markov-passwords).
It has written in Python, and It is very slow for experimentation purposes. This implementation is written in Go and is much faster.

## License
**markov-passwords** is licensed under the MIT license.

## Tests

### Speed
Trained on [Top304Thousand-probable-v2](https://weakpass.com/wordlist/1859) password list.

The `pv` results:
```shell
./bin/mkpass -sample top304K.json -m 100000000 | pv -lta > /dev/null
00:30 [1.44M/s]
```
Python implementation:
```shell
python sample.py | pv -lta > /dev/null
00:30 [5.30k/s]
```

### Unique Passwords
The `sort` results:
```shell
./bin/mkpass -sample top304K.json -m 1000000 | sort -u | wc -l
786577
```
Python implementation:
```shell
python sample.py | head -1000000 | sort -u | wc -l
541842
```

## Usage
### Get and Build
```shell
go get github.com/bujimuji/markov-passwords/...
cd $GOPATH/src/github.com/bujimuji/markov-passwords
make
```

### Train and Generate Passwords
```
Usage of ./bin/mkpass:
  -infile string
    	read training password list from file (- for stdin)
  -j int
    	number of channels for fan-out (default 64)
  -m uint
    	the amount of passwords created (default 1000000)
  -ngram int
    	number of ngrams to train markov model (default 3)
  -sample string
    	model path for generating new passwords
  -train string
    	model path for saving trained model in json format and exit
```
Example:
```shell
./bin/mkpass -infile assets/Top304Thousand-probable-v2.txt -m 100
```
#### Only Train
Example:
```shell
./bin/mkpass -infile assets/Top304Thousand-probable-v2.txt -train top304k.json
```

#### Only Generate
Example: 
```shell
./bin/mkpass -sample top304K.json -m 1000
```
