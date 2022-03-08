package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/notnil/chess"
)

func output(line string) {
	fmt.Println(line)
}

func main() {
	//test()
	//test_fen, _ := chess.FEN("rn1qkbnr/ppp2ppp/3p4/4p3/4P1b1/2N5/PPPP1PPP/R1B1KBNR w KQkq - 0 4")
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	//padTables()
	var smove string
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		if scanner.Scan() {
			smove = scanner.Text()
		}
		//fmt.Println((smove))

		if smove == "quit" {
			break
		} else if smove == "uci" {
			output("id name Abuahfish .v2")
			output("id author Amaechi Abuah")
			output("uciok")
		} else if smove == "isready" {
			output("readyok")
		} else if smove == "ucinewgame" {
			game = chess.NewGame(chess.UseNotation(chess.UCINotation{}))
			//output(game.String())
		} else if strings.HasPrefix(smove, "position") {
			params := strings.Fields(smove)
			idx := strings.Index(smove, "moves")
			var moveslist []string
			if idx >= 0 {
				moveslist = strings.Fields(smove[idx:])[1:]
			} /*else {
				moveslist = []string{}
			}*/
			var fen string
			if params[1] == "fen" {
				var fenpart string
				if idx >= 0 {
					fenpart = smove[:idx]
				} else {
					fenpart = smove
				}
				fen = strings.Join(strings.Fields(fenpart)[2:], " ")
				//output(fen)
			} else if params[1] == "startpos" {
				fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			} else {
				continue
			}
			our_fen, _ := chess.FEN(fen)
			game = chess.NewGame(our_fen, chess.UseNotation(chess.UCINotation{}))

			for _, move := range moveslist {
				game.MoveStr(move)
			}
			//output(game.Position().Board().Draw())
		} else if strings.HasPrefix(smove, "go") {
			depth := 1000
			movetime := -1
			show_thinking := true
			//var our_time int
			our_time := 3600000 // one hour
			params := strings.Fields(smove)[1:]
			var param, val string
			for i := 0; i < len(params)-1; i++ {
				param = params[i]
				val = params[i+1]
				if param == "depth" {
					depth, _ = strconv.Atoi(val)
				}
				if param == "movetime" {
					movetime, _ = strconv.Atoi(val)
				}
				if param == "wtime" && game.Position().Turn() == 1 {
					our_time, _ = strconv.Atoi(val)
				}
				if param == "btime" && game.Position().Turn() == 2 {
					our_time, _ = strconv.Atoi(val)
				}
			}

			moves_remain := 40

			start := time.Now().UnixMilli()
			var moves, moves_str string
			var ml []string
			searcher := NewSearcher()
			entry := Entry{}
			var score, sdepth int
			var usedtime int64
			spos := NewSPosition(*game.Position())
			for result := range searcher.search(spos) {
				moves = pv(&searcher, spos)
				sdepth = result.depth
				//fmt.Println(sdepth)
				if show_thinking {
					entry = searcher.tp_score[ScoreKey{spos.pos.Hash(), result.depth, true}]
					score = (entry.lower + entry.upper) / 2
					usedtime = time.Now().UnixMilli() - start
					if len(moves) < 50 {
						moves_str = moves
					} else {
						moves_str = ""
					}
					fmt.Printf("info depth %d score cp %d time %d nodes %d pv %v \n", result.depth, score, usedtime, searcher.nodes, moves_str)
				}

				if movetime > 0 && time.Now().UnixMilli()-start > int64(movetime) {
					break
				}
				if time.Now().UnixMilli()-start > int64(our_time)/int64(moves_remain) {
					break
				}
				if result.depth >= depth {
					break
				}
			}
			entry = searcher.tp_score[ScoreKey{spos.pos.Hash(), sdepth, true}]
			_, s := searcher.tp_move[spos.pos.Hash()], entry.lower

			if s == -MATE_UPPER {
				output("resign \n")
			} else {
				ml = strings.Fields(moves)
				if len(ml) > 1 {
					fmt.Printf("bestmove %v ponder %v \n", ml[0], ml[1])
				} else {
					fmt.Printf("bestmove %v \n", ml[0])
				}
				game.MoveStr(ml[0])
			}
			//output(game.Position().Board().Draw())

		} else {
			continue
		}

	}
}
