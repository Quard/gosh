package main

import (
	"bytes"
	"fmt"
)

var seedChars = []byte("abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789")

// GenerateNextSequence generate next sequence based on given last generated
func GenerateNextSequence(lastSeqeunce string) (string, error) {
	sequence := []byte(lastSeqeunce)
	for i := len(sequence) - 1; i >= 0; i-- {
		position := bytes.IndexByte(seedChars, sequence[i])
		if position == -1 {
			return "", fmt.Errorf("not valid char '%v' in identifier", sequence[i])
		}
		if position == len(seedChars)-1 {
			sequence[i] = seedChars[0]
		}
		if position < len(seedChars)-1 {
			sequence[i] = seedChars[position+1]
			break
		}
	}

	if bytes.Equal(sequence, bytes.Repeat(seedChars[0:1], len(sequence))) {
		return "a" + string(sequence), nil
	}

	return string(sequence), nil
}
