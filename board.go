package poly_y

import (
	"errors"
	"fmt"
	"io"
)

// Maximum number of fields on a Poly-Y board.
const MAX_FIELDS = 1000000

// Maximum number of sides for a Poly-Y board.
// (Chosen to be less than the width of an integer, to facilitate bitmasking.)
const MAX_SIDES = 31

// A Poly-Y board is an undirected graph with fields (indexed from 0) and sides
// (which are subsets of fields).  In the Board struct, Fields is an array
// of adjacency lists; i.e. Fields[i][j] == k implies that fields i and j
// are adjacent.  Sides is an array of arrays of field indices; i.e.
// Sides[i][j] == k implies that field k lies on the i-th side.  Note that
// fields may lay on multiple sides!
type Board struct {
	Fields [][]int
	Sides  [][]int
}

// Reads a list of field indices, converted from 1-based to 0-based.
// The first integer read indicates the number of indices that follow.
// e.g. if the input contains "3 2 4 6" readIndices() returns [1 3 5].
func readIndices(r io.Reader) []int {
	var i int
	if n, _ := fmt.Fscanf(r, "%d", &i); n != 1 {
		return nil
	}
	res := make([]int, i)
	for i = range res {
		if n, _ := fmt.Fscanf(r, "%d", &res[i]); n != 1 {
			return nil
		}
		res[i]-- // decrement values
	}
	return res
}

// Reads a Poly-Y board description.
func (b *Board) Read(r io.Reader) error {
	var nfield, nside int
	if n, _ := fmt.Fscanf(r, "Poly-Y %d %d", &nfield, &nside); n != 2 {
		return errors.New("board: invalid header")
	}
	if nfield < 1 || nfield > MAX_FIELDS {
		return errors.New("board: invalid number of fields")
	}
	if nside < 3 || nside > MAX_SIDES {
		return errors.New("board: invalid number of sides")
	}
	fields := make([][]int, nfield)
	sides := make([][]int, nside)
	for i := range fields {
		fields[i] = readIndices(r)
		if fields[i] == nil {
			return errors.New("board: invalid neighbour indices")
		}
	}
	for i := range sides {
		sides[i] = readIndices(r)
		if sides[i] == nil {
			return errors.New("board: invalid side indices")
		}
	}
	b.Fields = fields
	b.Sides = sides
	return nil
}
