package main

import (
	"math"

	"github.com/notnil/chess"
)

type Result struct {
	depth int
	move  chess.Move
	score int
}

type MoveScore struct {
	move  *chess.Move
	score int
}

type ScoreKey struct {
	hash  [16]byte
	depth int
	root  bool
}

type Entry struct {
	lower int
	upper int
}

// searchPosition type
type sPosition struct {
	pos   chess.Position
	score int
}

func NewSPosition(pos chess.Position) sPosition {
	score := eval(pos)
	return sPosition{pos, score}
}

func (spos sPosition) genMoves() []*chess.Move {
	return spos.pos.ValidMoves()
}

func (spos sPosition) nullmove() sPosition {
	new_pos := spos.pos.Update(&chess.Move{})
	return sPosition{*new_pos, -spos.score}
}

func (spos sPosition) move(mov chess.Move) sPosition {
	new_pos := spos.pos.Update(&mov)
	new_score := spos.score + spos.value(mov)
	return sPosition{*new_pos, -new_score}
}

func (spos sPosition) value(mov chess.Move) int {
	i, j := mov.S1(), mov.S2()
	p, q := spos.pos.Board().SquareMap()[i], spos.pos.Board().SquareMap()[j]
	qrs, krs := chess.A1, chess.H1
	//ps1, ps2 := chess.A8, chess.H8
	south := -8
	if p.Color().String() == "b" {
		i = 63 - i
		j = 63 - j
		qrs = 63 - chess.A8
		krs = 63 - chess.H8
		//ps1 = chess.A1
		//ps2 = chess.H1
		south = 8
	}
	// Actual move
	score := pst[p.Type().String()][j] - pst[p.Type().String()][i]
	// Capture
	if q != chess.NoPiece {
		score += pst[q.Type().String()][63-j]
	}
	// Castling
	if p.Type().String() == "k" && math.Abs(float64(i-j)) == 2 {
		score += pst["r"][(i+j)/2]
		if j > i {
			score -= pst["r"][krs]
		} else {
			score -= pst["r"][qrs]
		}
	}
	// Specail pawn stuff
	if p.Type().String() == "p" {
		if chess.A8 <= j && j <= chess.H8 {
			score += pst["q"][j] - pst["p"][j]
		}
		if mov.HasTag(chess.EnPassant) {
			score += pst["p"][63-(int(j)+south)]
		}
	}
	return score
}
