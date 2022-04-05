package main

import (
	"unicode"
)

type Result struct {
	depth int
	move  Move
	score int
}

type Move struct {
	from int
	to   int
}

type MoveScore struct {
	move  *Move
	score int
}

type ScoreKey struct {
	pos   sPosition
	depth int
	root  bool
}

type Entry struct {
	lower int
	upper int
}

type Pair struct {
	first  bool
	second bool
}

// searchPosition type
type sPosition struct {
	board string
	score int
	wc    Pair
	bc    Pair
	ep    int
	kp    int
}

func NewSPosition() sPosition {
	board := initial
	score := 0
	wc := Pair{true, true}
	bc := Pair{true, true}
	ep := 0
	kp := 0
	return sPosition{board, score, wc, bc, ep, kp}
}

func (spos sPosition) genMoves() []Move {
	res := []Move{}
	for i, p := range []rune(spos.board) {
		if !unicode.IsUpper(p) {
			continue
		}
		for _, d := range directions[p] {
			j := i + d
			for {
				q := rune(spos.board[j])
				// Stay inside the board, and off friendly pieces
				if unicode.IsSpace(q) || unicode.IsUpper(q) {
					break
				}
				// Pawn move, double move and capture
				if p == 'P' && containsInt([]int{N, N + N}, d) && q != '.' {
					break
				}
				if p == 'P' && d == N+N && (i < A1+N || spos.board[i+N] != '.') {
					break
				}
				if p == 'P' && containsInt([]int{N + W, N + E}, d) && q == '.' &&
					!containsInt([]int{spos.ep, spos.kp, spos.kp - 1, spos.kp + 1}, j) {
					break
				}
				// move it
				res = append(res, Move{i, j})
				// Stop crawlers from sliding, and sliding after captures
				if containsRune([]rune("PNK"), p) || unicode.IsLower(q) {
					break
				}
				// Castling, by sliding the rook next to the king
				if i == A1 && spos.board[j+E] == 'K' && spos.wc.first {
					res = append(res, Move{j + E, j + W})
				}
				if i == H1 && spos.board[j+W] == 'K' && spos.wc.second {
					res = append(res, Move{j + W, j + E})
				}
				j += d
			}
		}
	}
	return res
}

func (spos sPosition) rotate() sPosition {
	// Rotates the board, preserning enpassant
	var new_ep, new_kp int
	if spos.ep != 0 {
		new_ep = 119 - spos.ep
	} else {
		new_ep = 0
	}
	if spos.kp != 0 {
		new_kp = 119 - spos.ep
	} else {
		new_kp = 0
	}
	return sPosition{
		SwapCase(Reverse(spos.board)),
		-spos.score,
		spos.bc,
		spos.wc,
		new_ep,
		new_kp}
}

func (spos sPosition) nullmove() sPosition {
	// Like rotate, but clears ep and kp
	return sPosition{
		SwapCase(Reverse(spos.board)),
		-spos.score,
		spos.bc,
		spos.wc,
		0,
		0}
}

func (spos sPosition) move(mov Move) sPosition {
	i, j := mov.from, mov.to
	p := spos.board[i]
	var put = func(board string, i int, p rune) string { return board[:i] + string(p) + board[i+1:] }
	// Copy variables and reset ep and kp
	board := spos.board
	wc, bc, ep, kp := spos.wc, spos.bc, 0, 0
	score := spos.score + spos.value(mov)
	// Actual move
	board = put(board, j, rune(board[i]))
	board = put(board, i, '.')
	// Castling rights, we move the rook or capture the opponents
	if i == A1 {
		wc = Pair{false, wc.second}
	}
	if i == H1 {
		wc = Pair{wc.first, false}
	}
	if j == A8 {
		wc = Pair{bc.first, false}
	}
	if j == H8 {
		wc = Pair{false, bc.second}
	}
	// Castling
	if p == 'K' {
		wc = Pair{false, false}
		if Abs(j-i) == 2 {
			kp = (i + j) / 2
			if j < i {
				board = put(board, A1, '.')
			} else {
				board = put(board, H1, '.')
			}
			board = put(board, kp, 'R')
		}
	}
	// Pawn promotion, double move and en passant capture
	if p == 'P' {
		if A8 <= j && j <= H8 {
			board = put(board, j, 'Q')
		}
		if j-i == 2*N {
			ep = i + N
		}
		if j == spos.ep {
			board = put(board, j+S, '.')
		}
	}
	// We rotate the returned position, so its ready for the next player
	return sPosition{board, score, wc, bc, ep, kp}.rotate()
}

func (spos sPosition) value(mov Move) int {
	i, j := mov.from, mov.to
	p, q := rune(spos.board[i]), rune(spos.board[j])
	// Actual move
	score := pst[p][j] - pst[p][i]
	// Capture
	if unicode.IsLower(q) {
		score += pst[unicode.ToUpper(q)][119-j]
	}
	// Castling check detection
	if Abs(j-spos.kp) < 2 {
		score += pst['K'][119-j]
	}
	// Castling
	if p == 'K' && Abs(i-j) == 2 {
		score += pst['R'][(i+j)/2]
		if j < i {
			score -= pst['R'][A1]
		} else {
			score -= pst['R'][H1]
		}
	}
	// Special pawn stuff
	if p == 'P' {
		if A8 <= j && j <= H8 {
			score += pst['Q'][j] - pst['P'][j]
		}
		if j == spos.ep {
			score += pst['P'][119-(j+S)]
		}
	}
	return score
}
