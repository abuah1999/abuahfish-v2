package main

import (
	"strings"
)

func pv(s *Searcher, spos sPosition) string {
	res := []string{}
	for true {
		move, present := s.tp_move[spos.pos.Hash()]
		if !present {
			break
		}
		res = append(res, strings.ToLower(move.String()))
		spos = spos.move(move)
	}
	return strings.Join(res, " ")
}
