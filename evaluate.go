package main

import (
	"fmt"
	"unicode"
)

func test() {
	//padTables()
	//fmt.Println('e'==rune("yes"[1]))
	//position := parseFEN(FEN_INITIAL)
	//position = position.move(Move{81, 61})
	//fmt.Printf(position.board)

	fmt.Println(render(64))

}

//position2 := position.rotate()

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/*func contains[T comparable](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}*/

func containsInt(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func containsRune(s []rune, str rune) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func SwapCase(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if unicode.IsLower(runes[i]) {
			runes[i] = unicode.ToUpper(runes[i])
		} else if unicode.IsUpper(runes[i]) {
			runes[i] = unicode.ToLower(runes[i])
		}
	}
	return string(runes)
}

var piece = map[rune]int{
	'P': 100, 'N': 280, 'B': 320, 'R': 479, 'Q': 929, 'K': 60000}

var pst = map[rune][]int{
	'P': {0, 0, 0, 0, 0, 0, 0, 0,
		78, 83, 86, 73, 102, 82, 85, 90,
		7, 29, 21, 44, 40, 31, 44, 7,
		-17, 16, -2, 15, 14, 0, 15, -13,
		-26, 3, 10, 9, 6, 1, 0, -23,
		-22, 9, 5, -11, -10, -2, 3, -19,
		-31, 8, -7, -37, -36, -14, 3, -31,
		0, 0, 0, 0, 0, 0, 0, 0},

	'N': {-66, -53, -75, -75, -10, -55, -58, -70,
		-3, -6, 100, -36, 4, 62, -4, -14,
		10, 67, 1, 74, 73, 27, 62, -2,
		24, 24, 45, 37, 33, 41, 25, 17,
		-1, 5, 31, 21, 22, 35, 2, 0,
		-18, 10, 13, 22, 18, 15, 11, -14,
		-23, -15, 2, 0, 2, 0, -23, -20,
		-74, -23, -26, -24, -19, -35, -22, -69},

	'B': {-59, -78, -82, -76, -23, -107, -37, -50,
		-11, 20, 35, -42, -39, 31, 2, -22,
		-9, 39, -32, 41, 52, -10, 28, -14,
		25, 17, 20, 34, 26, 25, 15, 10,
		13, 10, 17, 23, 17, 16, 0, 7,
		14, 25, 24, 15, 8, 25, 20, 15,
		19, 20, 11, 6, 7, 6, 20, 16,
		-7, 2, -15, -12, -14, -15, -10, -10},

	'R': {35, 29, 33, 4, 37, 33, 56, 50,
		55, 29, 56, 67, 55, 62, 34, 60,
		19, 35, 28, 33, 45, 27, 25, 15,
		0, 5, 16, 13, 18, -4, -9, -6,
		-28, -35, -16, -21, -13, -29, -46, -30,
		-42, -28, -42, -25, -25, -35, -26, -46,
		-53, -38, -31, -26, -29, -43, -44, -53,
		-30, -24, -18, 5, -2, -18, -31, -32},

	'Q': {6, 1, -8, -104, 69, 24, 88, 26,
		14, 32, 60, -10, 20, 76, 57, 24,
		-2, 43, 32, 60, 72, 63, 43, 2,
		1, -16, 22, 17, 25, 20, -13, -6,
		-14, -15, -2, -5, -1, -10, -20, -22,
		-30, -6, -13, -11, -16, -11, -16, -27,
		-36, -18, 0, -19, -15, -15, -21, -38,
		-39, -30, -31, -13, -31, -36, -34, -42},

	'K': {4, 54, 47, -99, -99, 60, 83, -62,
		-32, 10, 55, 56, 56, 55, 10, 3,
		-62, 12, -57, 44, -67, 28, 37, -31,
		-55, 50, 11, -4, -19, 13, 0, -49,
		-55, -43, -52, -28, -51, -47, -8, -50,
		-47, -42, -43, -79, -64, -32, -29, -32,
		-4, 3, -14, -50, -57, -18, 13, 4,
		17, 30, -3, -14, 6, -1, 40, 18}}

func padTables() {
	for k, table := range pst {
		padrow := func(row []int) []int {
			for i := 0; i < len(row); i++ {
				row[i] += piece[k]
			}
			//add 0 to the front
			row = append([]int{0}, row...)
			//add 0 to the back
			row = append(row, 0)
			return row
		}
		new_table := []int{}
		for i := 0; i < 8; i++ {
			new_table = append(new_table, padrow(table[i*8:i*8+8])...)
		}
		for i := 0; i < 20; i++ {
			//add 0 to the front
			new_table = append(new_table, 0)
			copy(new_table[1:], new_table)
			new_table[0] = 0
			//add 0 to the back
			new_table = append(new_table, 0)
		}
		pst[k] = new_table
	}
}

const A1, H1, A8, H8 int = 91, 98, 21, 28
const initial string = "         \n" + //0-9
	"         \n" + //10-19
	" rnbqkbnr\n" + //20-29
	" pppppppp\n" + //30-39
	" ........\n" + //40-49
	" ........\n" + //50-59
	" ........\n" + //60-69
	" ........\n" + //70-79
	" PPPPPPPP\n" + //80-89
	" RNBQKBNR\n" + //90-99
	"         \n" + //100-109
	"         \n" //110-119

const N, E, S, W int = -10, 1, 10, -1

var directions = map[rune][]int{
	'P': {N, N + N, N + W, N + E},
	'N': {N + N + E, E + N + E, E + S + E, S + S + E, S + S + W, W + S + W, W + N + W, N + N + W},
	'B': {N + E, S + E, S + W, N + W},
	'R': {N, E, S, W},
	'Q': {N, E, S, W, N + E, S + E, S + W, N + W},
	'K': {N, E, S, W, N + E, S + E, S + W, N + W}}
