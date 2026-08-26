package search

import (
	"sort"
	"strings"
)

type Result struct {
	Entry SearchEntry
	Score int
}

func Search(index *Index, query string) []Result {
	if query == "" {
		results := make([]Result, len(index.Entries))
		for i, entry := range index.Entries {
			results[i] = Result{Entry: entry, Score: 0}
		}
		return results
	}

	query = strings.ToLower(query)
	var results []Result

	for _, entry := range index.Entries {
		score := scoreEntry(entry, query)
		if score > 0 {
			results = append(results, Result{Entry: entry, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func scoreEntry(entry SearchEntry, query string) int {
	score := 0
	words := strings.Fields(query)

	title := strings.ToLower(entry.Title)
	filename := strings.ToLower(entry.Filename)
	dir := strings.ToLower(entry.Dir)
	content := strings.ToLower(entry.Content)

	for _, word := range words {
		wordScore := 0

		if strings.Contains(title, word) {
			wordScore += 100
			if title == word {
				wordScore += 50
			}
		}

		if strings.Contains(filename, word) {
			wordScore += 80
			if filename == word {
				wordScore += 40
			}
		}

		if strings.Contains(dir, word) {
			wordScore += 30
		}

		if strings.Contains(content, word) {
			wordScore += 10
			wordScore += min(strings.Count(content, word), 10)
		}

		if wordScore == 0 {
			return 0
		}
		score += wordScore
	}

	return score
}
