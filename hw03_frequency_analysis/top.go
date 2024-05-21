package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	// Split text to words by space
	words := strings.Fields(text)

	// Get words count map
	count := make(map[string]int)
	for _, word := range words {
		count[word]++
	}

	// Create a structure with the words and counts
	type wordCount struct {
		word  string
		count int
	}

	// Place each word into a slice of wordsCount structures
	wordsFrequencies := make([]wordCount, 0, len(count))
	for key, value := range count {
		wordsFrequencies = append(wordsFrequencies, wordCount{key, value})
	}

	// Sort wordsFrequencies slice
	sort.Slice(wordsFrequencies, func(i, j int) bool {
		// If count equal, sort lexicographically
		if wordsFrequencies[i].count == wordsFrequencies[j].count {
			return wordsFrequencies[i].word < wordsFrequencies[j].word
		}

		// Sort by count
		return wordsFrequencies[i].count > wordsFrequencies[j].count
	})

	// Top 10 words
	topWords := make([]string, 0, 10)
	for i := range wordsFrequencies {
		topWords = append(topWords, wordsFrequencies[i].word)

		if i == 9 {
			break
		}
	}

	return topWords
}
