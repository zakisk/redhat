package models

import "sync"

type WordCount struct {
	Mu              sync.Mutex
	TotalWordsCount int64
	TotalFileCount  int
	WordsCountMap   map[string]int
}
