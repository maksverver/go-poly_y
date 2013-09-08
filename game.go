package poly_y

import (
	"io"
	"strconv"
)

type Move int

type State struct {
	Board *Board
	Moves []Move
}

func (m Move) String() string {
	return strconv.FormatInt(int64(m), 10)
}

func (s *State) Over() bool {
	if len(s.Moves) >= len(s.Board.Fields) {
		return true
	}
	a, b := s.Scores()
	n := len(s.Board.Sides) / 2
	return a > n || b > n
}

func (s *State) Next() int {
	return len(s.Moves) % 2
}

func (s *State) Occupied(field int) bool {
	for _, m := range s.Moves {
		if field == int(m) {
			return true
		}
	}
	return false
}

func (s *State) ListMoves() []interface{} {
	moves := make([]interface{}, 0)
	for i := 1; i <= len(s.Board.Fields); i++ {
		if !s.Occupied(i) {
			moves = append(moves, Move(i))
		}
	}
	return moves
}

func (s *State) Execute(arg interface{}) bool {
	if m, ok := arg.(Move); ok {
		i := int(m)
		if i > 0 && i <= len(s.Board.Fields) && !s.Occupied(i) {
			s.Moves = append(s.Moves, m)
			return true
		}
	}
	return false
}

func (s *State) Scores() (int, int) {
	score1, score2 := 0, 0
	colors := make([]int, len(s.Board.Fields))
	sides := make([]uint, len(s.Board.Fields))
	visited := make([]bool, len(s.Board.Fields))
	for i, m := range s.Moves {
		colors[int(m)-1] = 1 - i%2*2
	}
	for i, side := range s.Board.Sides {
		for _, j := range side {
			sides[j] |= 1 << uint(i)
		}
	}
	var dfs func(i int) uint
	dfs = func(i int) uint {
		visited[i] = true
		res := sides[i]
		for _, j := range s.Board.Fields[i] {
			if colors[i] == colors[j] && !visited[j] {
				res |= dfs(j)
			}
		}
		return res
	}
	for i, c := range colors {
		if c != 0 && !visited[i] {
			x := dfs(i)
			y := x & (x - 1)
			z := y & (y - 1)
			if z != 0 {
				// Chain touches at least three sides!  Check which ones.
				for j := range s.Board.Sides {
					var m uint = 1<<uint(j) | 1<<uint((j+1)%len(s.Board.Sides))
					if x&m == m {
						if colors[i] > 0 {
							score1++
						} else {
							score2++
						}
					}
				}
			}
		}
	}
	return score1, score2
}

func (s *State) WriteLog(w io.Writer) {
	text := ""
	col := 0
	for _, move := range s.Moves {
		s := move.String()
		if col > 0 {
			if col+1+len(s) < 80 {
				text += " "
				col++
			} else {
				text += "\n"
				col = 0
			}
		}
		text += s
		col += len(s)
	}
	if col > 0 {
		text += "\n"
	}
	w.Write([]byte(text))
}
