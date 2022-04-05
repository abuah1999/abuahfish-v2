package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const FEN_INITIAL string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
const WHITE, BLACK int = 0, 1

func get_color(pos sPosition) int {
	if strings.HasPrefix(pos.board, "\n") {
		return BLACK
	} else {
		return WHITE
	}
}

func parse(c string) int {
	fil, rank := int(c[0]-'a'), int(c[1]-'1')
	return A1 + fil - 10*rank
}

func render(i int) string {
	rank, fil := (i-H1)/10, (i-A8)%10
	return string(fil+'a') + strconv.Itoa(1-rank)
}

func mrender(pos sPosition, m Move) string {
	var p string
	if A8 <= m.to && m.to <= H8 && pos.board[m.from] == 'P' {
		p = "q"
	} else {
		p = ""
	}
	if get_color(pos) == BLACK {
		m = Move{119 - m.from, 119 - m.to}
	}
	return render(m.from) + render(m.to) + p
}

func mparse(color int, move string) Move {
	m := Move{parse(move[0:2]), parse(move[2:4])}
	if color == WHITE {
		return m
	} else {
		return Move{119 - m.from, 119 - m.to}
	}
}

func can_kill_king(pos sPosition) bool {
	for _, m := range pos.genMoves() {
		if pos.value(m) >= MATE_LOWER {
			return true
		}
	}
	return false
}

func parseFEN(fen string) sPosition {
	info := strings.Fields(fen)
	board, color, castling, enpas := info[0], info[1], info[2], info[3]
	//byte_board := []byte(board)
	digit_exp := regexp.MustCompile("[0-9]")
	board = digit_exp.ReplaceAllStringFunc(board, func(s string) string {
		//s = s[1:len(s)-1]
		return (strings.Repeat(".", int(s[0]-'0')))
	})
	rune_board := []rune(strings.Repeat(" ", 21) + strings.Join(strings.Split(board, "/"), "  ") + strings.Repeat(" ", 21))
	for i := 9; i < 120; i += 10 {
		rune_board[i] = '\n'
	}
	board = string(rune_board)
	wc := Pair{strings.Contains(castling, "Q"), strings.Contains(castling, "K")}
	bc := Pair{strings.Contains(castling, "k"), strings.Contains(castling, "q")}
	var ep int
	if enpas != "-" {
		ep = parse(enpas)
	} else {
		ep = 0
	}
	score := 0
	for i, p := range []rune(board) {
		if unicode.IsUpper(p) {
			score += pst[p][i]
		} else if unicode.IsLower(p) {
			score -= pst[unicode.ToUpper(p)][119-i]
		}
	}
	pos := sPosition{board, score, wc, bc, ep, 0}
	if color == "w" {
		return pos
	} else {
		return pos.rotate()
	}
}

func pv(s *Searcher, spos sPosition, depth int) string {

	res := []string{}
	for true {
		s.tp_moveMutex.RLock()
		move, present := s.tp_move[spos]
		s.tp_moveMutex.RUnlock()
		//fmt.Println(depth)
		if !present || can_kill_king(spos.move(move)) {

			//fmt.Println("oops")

			break
		}
		res = append(res, mrender(spos, move))
		spos = spos.move(move)
	}
	/*for _, v := range s.tp_score {
		fmt.Println(v)
	}*/
	return strings.Join(res, " ")
}
