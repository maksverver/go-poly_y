// Poly-Y is the abstract board game Y played on a polygonal board.
package poly_y

import (
	"io"
	"strconv"
)

// A move is either a positive (1-based) field index, or the number -1 which
// requests swapping positions (only allowed as the second move).
type Move int

// The Poly-Y game state consists of a board descriptor (which is not modified
// by any method) and an array of valid moves executed so far.
type State struct {
	Board *Board
	Moves []Move
}

func (m Move) String() string {
	return strconv.FormatInt(int64(m), 10)
}

// Returns whether the game is over.
func (s *State) Over() bool {
	// See if there are any unoccupied fields left:
	unoccupied := len(s.Board.Fields) - len(s.Moves)
	if len(s.Moves) >= 2 && s.Moves[1] == -1 {
		unoccupied++
	}
	if unoccupied <= 0 {
		return false
	}

	// Check if either player has captured a majority of the corners:
	a, b := s.Scores()
	n := len(s.Board.Sides) / 2
	return a > n || b > n
}

// Returns which player (0 or 1) must play next.
func (s *State) Next() int {
	return len(s.Moves) % 2
}

// Returns whether the field with the given 1-based index is occupied.
func (s *State) occupied(field int) bool {
	for _, m := range s.Moves {
		i := int(m)
		if i != -1 && i == int(field) {
			return true
		}
	}
	return false
}

// Lists all valid moves in the current state.
func (s *State) ListMoves() []interface{} {
	moves := make([]interface{}, 0)
	if len(s.Moves) == 1 {
		moves = append(moves, -1)
	}
	for i := 1; i <= len(s.Board.Fields); i++ {
		if !s.occupied(i) {
			moves = append(moves, Move(i))
		}
	}
	return moves
}

// Returns whether the given move is valid in the current state.
func (s *State) valid(m Move) bool {
	i := int(m)
	if i == -1 {
		return len(s.Moves) == 1
	}
	return 1 <= i && i <= len(s.Board.Fields) && !s.occupied(i)
}

// If the argument is a valid move, it is executed, and true is returned.
// Otherwise, the game state is unchanged, and false is returned.
func (s *State) Execute(arg interface{}) bool {
	if m, ok := arg.(Move); ok && s.valid(m) {
		s.Moves = append(s.Moves, m)
		return true
	}
	return false
}

// Returns the current score (in number of corners captured) for both players.
func (s *State) Scores() (int, int) {
	score1, score2 := 0, 0
	colors := make([]int, len(s.Board.Fields))
	sides := make([]uint, len(s.Board.Fields))
	visited := make([]bool, len(s.Board.Fields))
	for i, m := range s.Moves {
		j := int(m)
		if j == -1 {
			j = int(s.Moves[i-1])
		}
		colors[j-1] = 1 - i%2*2
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

// Writes the game transcript.  The format is a sequence of 1-based field
// indices (with -1 indicating swap) with lines folded to 79 columns.
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
