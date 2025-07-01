package main

import (
	"strings"
)

func cleanMessage(message string) string {
	const replaceWord = "****"
	naughtyWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitMessage := strings.Split(message, " ")
	for x, word := range splitMessage {
		for _, naughtyWord := range naughtyWords {
			if word == naughtyWord {
				splitMessage[x] = replaceWord
			}
		}
	}
	return strings.Join(splitMessage, " ")
}
