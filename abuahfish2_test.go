package main

import (
	"testing"
)

func BenchmarkFish(b *testing.B) {
	//game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	searcher := NewSearcher()
	spos := parseFEN(FEN_INITIAL)
	for i := 0; i < b.N; i++ {
		for result := range searcher.search(spos) {
			if result.depth == 5 {
				break
			}
		}
	}
}
