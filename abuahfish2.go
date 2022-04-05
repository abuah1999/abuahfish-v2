package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func output(line string) {
	fmt.Println(line)
}

func main() {

	//test_fen, _ := chess.FEN("rn1qkbnr/ppp2ppp/3p4/4p3/4P1b1/2N5/PPPP1PPP/R1B1KBNR w KQkq - 0 4")
	padTables()
	//test()
	pos := parseFEN(FEN_INITIAL)
	f, _ := os.Create("log")
	defer f.Close()
	var smove string
	var stack []string
	var color int
	scanner := bufio.NewScanner(os.Stdin)
	for true {
		if len(stack) > 0 {
			smove = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			if scanner.Scan() {
				smove = scanner.Text()
			}
			//fmt.Println((smove))
		}
		f.WriteString(smove + "\n")
		if smove == "quit" {
			break
		} else if smove == "uci" {
			f.WriteString("uci received\n")
			output("id name Abuahfish .v2")
			output("id author Amaechi Abuah")
			output("option name Hash type spin default 1 min 1 max 128")
			output("uciok")
		} else if smove == "isready" {
			f.WriteString("isready received\n")
			output("readyok")
		} else if smove == "ucinewgame" {
			stack = append(stack, "position fen "+FEN_INITIAL)
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
				fen = FEN_INITIAL
			} else {
				continue
			}
			pos = parseFEN(fen)

			if strings.Fields(fen)[1] == "w" {
				color = WHITE
			} else {
				color = BLACK
			}
			for _, move := range moveslist {
				pos = pos.move(mparse(color, move))
				color = 1 - color
			}
			//output(game.Position().Board().Draw())
		} else if strings.HasPrefix(smove, "go") {
			f.WriteString("go received\n")
			depth := 4
			movetime := -1
			show_thinking := false
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
				if param == "wtime" {
					our_time, _ = strconv.Atoi(val)
				}
				/*if param == "btime" && color == 1 {
					our_time, _ = strconv.Atoi(val)
				}*/
			}

			moves_remain := 400

			start := time.Now().UnixMilli()
			var moves, moves_str string
			var ml []string
			searcher := NewSearcher()
			entry := Entry{}
			var score, sdepth int
			var usedtime int64

			for result := range searcher.search(pos) {
				moves = pv(&searcher, pos, result.depth)
				//fmt.Println(moves)
				sdepth = result.depth
				//fmt.Println(sdepth)
				if show_thinking {
					searcher.tp_scoreMutex.RLock()
					entry = searcher.tp_score[ScoreKey{pos, result.depth, true}]
					searcher.tp_scoreMutex.RUnlock()
					score = (entry.lower + entry.upper) / 2
					//fmt.Println(entry.lower, entry.upper)
					usedtime = time.Now().UnixMilli() - start
					if len(moves) < 50 {
						moves_str = moves
					} else {
						moves_str = ""
					}
					fmt.Printf("info depth %d score cp %d time %d nodes %d pv %v\n", result.depth, score, usedtime, searcher.nodes, moves_str)
					//searcher.tp_scoreMutex.Lock()
					/*for _, v := range searcher.tp_score {
						fmt.Println(v)
					}*/
					//searcher.tp_scoreMutex.Unlock()
				}

				if movetime > 0 && time.Now().UnixMilli()-start > int64(movetime/10) {
					break
				}
				if time.Now().UnixMilli()-start > int64(our_time)/int64(moves_remain) {
					break
				}
				if result.depth >= depth {
					break
				}
			}
			searcher.tp_scoreMutex.RLock()
			entry = searcher.tp_score[ScoreKey{pos, sdepth, true}]
			searcher.tp_scoreMutex.RUnlock()
			searcher.tp_moveMutex.RLock()
			_, s := searcher.tp_move[pos], entry.lower
			searcher.tp_moveMutex.RUnlock()

			if s == -MATE_UPPER {
				output("resign \n")
			} else {
				ml = strings.Fields(moves)
				if len(ml) > 1 {
					fmt.Printf("bestmove %v ponder %v \n", ml[0], ml[1])
					f.WriteString(ml[0] + "\n")
				} else {
					fmt.Printf("bestmove %v \n", ml[0])
				}
				pos = pos.move(mparse(color, ml[0]))
				color = 1 - color
			}
			//output(game.Position().Board().Draw())

		} else {
			continue
		}

	}
}
