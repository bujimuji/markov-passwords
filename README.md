# Markov-chain passwords generator
The Markov-chain password generator is one of the password guessing methods. The trained model represents the probability of transition between specific characters. 

The [original implementation](https://github.com/brannondorsey/markov-passwords) has written in Python, and It is very slow for experimentation purposes. This implementation is written in Go and is more faster.

## license
**markov-passwords** is licensed under the MIT license.

## Get Requirements
```
go get github.com/sirupsen/logrus
```

## Speed Test
Trained on [Top304Thousand-probable-v2](https://weakpass.com/wordlist/1859) password list.

The `pv` results:
```shell
./bin/sample -infile top304K.json | pv -lta > /dev/null
00:30 [ 716k/s]
```
Python implementation:
```shell
python sample.py | pv -lta > /dev/null
00:30 [5.30k/s]
```
## Usage
### Clone and Build
```
go get github.com/bujimuji/markov-passwords
cd $GOPATH/src/github.com/bujimuji/markov-passwords
make
```
### Train
```
Usage of ./bin/train:
  -infile string
    	read training data from file
  -max-ngram int
    	maximum number of ngrams to train (default 3)
  -outfile string
    	write training data to file
```
Example:
```
./bin/train -infile Top304Thousand-probable-v2 -outfile top304K.json
```

### Generate
```
Usage of ./bin/sample:
  -b int
    	block size of prints (default 3000)
  -infile string
    	read training data from file
  -j int
    	number of threads (default 32)
  -max-ngram int
    	maximum number of ngrams to train (default 3)
```
Eample: 
```
./bin/sample -infile top304K.json
```

## TODO
- [ ] Use an insecure and fast random number generator