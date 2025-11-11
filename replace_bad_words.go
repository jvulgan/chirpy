package main

import (
	"slices"
	"strings"
)

func replaceBadWords(original string) string {
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	var clean []string
	for _, word := range strings.Split(original, " ") {
		if slices.Contains(bad_words, strings.ToLower(word)) {
			clean = append(clean, "****")
		} else {
			clean = append(clean, word)
		}
	}
	return strings.Join(clean, " ")
}
