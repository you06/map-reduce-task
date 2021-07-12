package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/c2h5oh/datasize"
)

var (
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	mustRepeatWords    = [][]byte{}
	bytesPerFile       uint64
	firstNonRepeatFile int
	firstNonRepeatWord []byte
)

func Run(dataSize, outputPath string, partition int) {
	var (
		v   datasize.ByteSize
		err error
	)
	err = v.UnmarshalText([]byte(dataSize))
	if err != nil {
		fmt.Println("Parse datasize error", err)
		return
	}
	if _, e := os.Stat(outputPath); os.IsNotExist(e) {
		err = os.Mkdir(outputPath, 0755)
		if err != nil {
			fmt.Println("Mkdir error", err)
			return
		}
	}
	bytes := v.Bytes()
	bytesPerFile = bytes / uint64(partition)
	firstNonRepeatFile = globalRand.Intn(partition)
	firstNonRepeatWord = randWord(globalRand, randLen(globalRand, MAX_LEN))
	fmt.Printf("The first non-repeat word is \"%s\",\nStart generating test data.\n", string(firstNonRepeatWord))

	for i := 0; i < partition; i++ {
		err = oneFile(outputPath, i, partition)
		if err != nil {
			fmt.Println("Generate error", err)
			return
		}
	}
	fmt.Println("Generate success")
}

func oneFile(dir string, seq int, total int) error {
	fileBytes := make([]byte, bytesPerFile+MAX_LEN)
	bytesLeft := bytesPerFile
	index := uint64(0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	afterFirst := seq > firstNonRepeatFile

	for {
		mustRepeat := !afterFirst
		var word []byte

		if seq == firstNonRepeatFile && !afterFirst && r.Float64() < float64(index)/float64(bytesPerFile) {
			word = firstNonRepeatWord
			mustRepeat = false
			afterFirst = true
		} else {
			if r.Float64() < float64(len(mustRepeatWords))/float64(bytesPerFile/MAX_LEN) ||
				(seq == total-1 && len(mustRepeatWords) > 0) {
				// get word from mustRepeatWords
				index := r.Intn(len(mustRepeatWords))
				word = mustRepeatWords[index]
				mustRepeatWords = append(mustRepeatWords[:index], mustRepeatWords[index+1:]...)
				mustRepeat = false
			} else {
				// rand word
				word = randWord(r, randLen(r, bytesLeft))
			}
		}

		if mustRepeat {
			mustRepeatWords = append(mustRepeatWords, word)
		}
		word = append(word, ' ')

		copy(fileBytes[index:], word)
		index += uint64(len(word))

		if bytesLeft <= uint64(len(word)) {
			// last word in this file
			break
		}
		bytesLeft -= uint64(len(word))
	}
	if seq == firstNonRepeatFile && !afterFirst {
		copy(fileBytes[index:], firstNonRepeatWord)
	}
	trimLen := len(fileBytes) - 1
	for fileBytes[trimLen] == ' ' || fileBytes[trimLen] == 0 {
		trimLen--
	}
	fileBytes = fileBytes[:trimLen+1]
	return ioutil.WriteFile(path.Join(dir, fmt.Sprintf("%d.txt", seq)), fileBytes, 0644)
}

const MAX_LEN uint64 = 100

func randLen(r *rand.Rand, maxLen uint64) int {
	if maxLen > MAX_LEN {
		maxLen = MAX_LEN
	}
	return 1 + r.Intn(int(maxLen))
}

var (
	letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func randWord(r *rand.Rand, l int) []byte {
	b := make([]byte, l)
	for {
		for i := range b {
			b[i] = letterRunes[r.Intn(len(letterRunes))]
		}
		if !bytes.Equal(b, firstNonRepeatWord) {
			break
		}
	}
	return b
}
