package main

import (
	"math"
	"sort"
	"strings"

	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/notnil/chess"
)

var MATE_UPPER int = piece["k"] + 10*piece["q"]
var MATE_LOWER int = piece["k"] - 10*piece["q"]

const TABLE_SIZE int = 1e7
const QS_LIMIT int = 219
const EVAL_ROUGHNESS int = 13
const DRAW_TEST bool = true

type Searcher struct {
	tp_score      map[ScoreKey]Entry
	tp_move       map[[16]byte]chess.Move
	tp_scoreMutex *sync.RWMutex
	tp_moveMutex  *sync.RWMutex
	history       mapset.Set
	nodes         int
}

func NewSearcher() Searcher {
	tp_score := make(map[ScoreKey]Entry)
	tp_move := make(map[[16]byte]chess.Move)
	tp_scoreMutex := sync.RWMutex{}
	tp_moveMutex := sync.RWMutex{}
	history := mapset.NewSet()
	nodes := 0
	return Searcher{tp_score, tp_move, &tp_scoreMutex, &tp_moveMutex, history, nodes}
}

func (s *Searcher) moves(spos sPosition, gamma int, depth int, root bool) chan MoveScore {
	//res := []MoveScore{}
	c := make(chan MoveScore)
	pos_hash := spos.pos.Hash()
	hasRook := strings.Contains(spos.pos.String(), "R")
	hasRookb := strings.Contains(spos.pos.String(), "r")
	hasKnight := strings.Contains(spos.pos.String(), "N")
	hasKnightb := strings.Contains(spos.pos.String(), "n")
	hasBishop := strings.Contains(spos.pos.String(), "B")
	hasBishopb := strings.Contains(spos.pos.String(), "b")
	hasQueen := strings.Contains(spos.pos.String(), "Q")
	hasQueenb := strings.Contains(spos.pos.String(), "q")
	s.tp_moveMutex.Lock()
	killer, killerpresent := s.tp_move[pos_hash]
	s.tp_moveMutex.Unlock()
	go func() {
		if depth > 0 && !root && (hasRook || hasRookb || hasKnight || hasKnightb || hasBishop || hasBishopb || hasQueen || hasQueenb) {
			c <- MoveScore{nil, -s.bound(spos.nullmove(), 1-gamma, depth-3, false)}
		}
		if depth == 0 {
			c <- MoveScore{nil, spos.score}
		}

		if killerpresent && (depth > 0 || spos.value(killer) >= QS_LIMIT) {
			c <- MoveScore{&killer, -s.bound(spos.move(killer), 1-gamma, depth-1, false)}
		}
		other_moves := spos.genMoves()
		sort.Slice(other_moves, func(i, j int) bool {
			return spos.value(*other_moves[i]) > spos.value(*other_moves[j])
		})
		i := 0
		for {
			if i == len(other_moves) {
				close(c)
				return
			}
			if depth > 0 || spos.value(*other_moves[i]) >= QS_LIMIT {
				c <- MoveScore{other_moves[i], -s.bound(spos.move(*other_moves[i]), 1-gamma, depth-1, false)}
			}
			i++
		}
	}()
	return c
}

/*if depth > 0 && !root && (hasRook || hasRookb || hasKnight || hasKnightb || hasBishop || hasBishopb || hasQueen || hasQueenb) {
	res = append(res, MoveScore{chess.Move{}, -s.bound(spos.nullmove(), 1-gamma, depth-3, false)})
}
if depth == 0 {
	res = append(res, MoveScore{chess.Move{}, spos.score})
}
killer, killerpresent := s.tp_move[pos_hash]
if killerpresent && (depth > 0 || spos.value(killer) >= QS_LIMIT) {
	res = append(res, MoveScore{killer, -s.bound(spos.move(killer), 1-gamma, depth-1, false)})
}
other_moves := spos.genMoves()
sort.Slice(other_moves, func(i, j int) bool {
	return spos.value(*other_moves[i]) > spos.value(*other_moves[j])
})
for _, move := range other_moves {
	if depth > 0 || spos.value(*move) >= QS_LIMIT {
		res = append(res, MoveScore{*move, -s.bound(spos.move(*move), 1-gamma, depth-1, false)})
	}
}
return res
*/

func (s *Searcher) bound(spos sPosition, gamma int, depth int, root bool) int {
	s.nodes += 1
	pos_hash := spos.pos.Hash()
	depth = int(math.Max(float64(depth), 0))

	if spos.score <= -MATE_LOWER {
		return -MATE_UPPER
	}

	if DRAW_TEST {
		if !root && s.history.Contains(pos_hash) {
			return 0
		}
	}
	s.tp_scoreMutex.RLock()
	entry, scorepresent := s.tp_score[ScoreKey{pos_hash, depth, root}]
	s.tp_scoreMutex.RUnlock()
	s.tp_moveMutex.RLock()
	_, movepresent := s.tp_move[pos_hash]
	s.tp_moveMutex.RUnlock()
	if !scorepresent {
		entry = Entry{-MATE_UPPER, MATE_UPPER}
	}
	if entry.lower >= gamma && (!root || movepresent) {
		return entry.lower
	}
	if entry.upper < gamma {
		return entry.upper
	}

	best := -MATE_UPPER
	for ms := range s.moves(spos, gamma, depth, root) {
		best = int(math.Max(float64(best), float64(ms.score)))
		if best >= gamma {
			if len(s.tp_move) > TABLE_SIZE {
				s.tp_move = map[[16]byte]chess.Move{}
			}

			if ms.move != nil {
				s.tp_moveMutex.Lock()
				s.tp_move[pos_hash] = *ms.move
				s.tp_moveMutex.Unlock()
			}

			break
		}
	}

	if best < gamma && best < 0 && depth > 0 {
		is_dead := func(sp sPosition) bool {
			for _, m := range sp.genMoves() {
				if sp.value(*m) >= MATE_LOWER {
					return true
				}
			}
			return false
		}
		all_is_dead := true
		for _, m := range spos.genMoves() {
			if !is_dead(spos.move(*m)) {
				all_is_dead = false
			}
		}
		if all_is_dead {
			in_check := is_dead(spos.nullmove())
			if in_check {
				best = -MATE_UPPER
			} else {
				best = 0
			}
		}
	}
	if len(s.tp_score) > TABLE_SIZE {
		s.tp_score = map[ScoreKey]Entry{}
	}
	if best >= gamma {
		s.tp_scoreMutex.Lock()
		s.tp_score[ScoreKey{pos_hash, depth, root}] = Entry{best, entry.upper}
		s.tp_scoreMutex.Unlock()
	}
	if best < gamma {
		s.tp_scoreMutex.Lock()
		s.tp_score[ScoreKey{pos_hash, depth, root}] = Entry{entry.lower, best}
		s.tp_scoreMutex.Unlock()
	}
	return best
}

func (s *Searcher) search(spos sPosition) chan Result {
	c := make(chan Result)
	s.nodes = 0
	depth := 1
	pos_hash := spos.pos.Hash()
	//var res []Result
	if DRAW_TEST {
		s.history = mapset.NewSet()
		s.tp_score = map[ScoreKey]Entry{}
	}

	go func() {
		for {
			if depth == 1000 {
				close(c)
				return
			}
			lower, upper := -MATE_UPPER, MATE_UPPER
			for lower < upper-EVAL_ROUGHNESS {
				gamma := (lower + upper + 1) / 2
				score := s.bound(spos, gamma, depth, true)
				if score >= gamma {
					lower = score
				}
				if score < gamma {
					upper = score
				}
			}
			_ = s.bound(spos, lower, depth, true)
			c <- Result{depth, s.tp_move[pos_hash], s.tp_score[ScoreKey{pos_hash, depth, true}].lower}
			depth++
		}
	}()
	return c
}

/*for depth := 1; depth < 1000; depth++ {
		lower, upper := -MATE_UPPER, MATE_UPPER
		for lower < upper-EVAL_ROUGHNESS {
			gamma := (lower + upper + 1) / 2
			score := s.bound(spos, gamma, depth, true)
			if score >= gamma {
				lower = score
			}
			if score < gamma {
				upper = score
			}
		}
		_ = s.bound(spos, lower, depth, true)
		res = append(res, Result{depth, s.tp_move[pos_hash], s.tp_score[ScoreKey{pos_hash, depth, true}].lower})
	}

}*/
