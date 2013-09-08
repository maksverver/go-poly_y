package poly_y

import (
	"errors"
	"fmt"
	"io"
)

const MAX_FIELDS = 1000000
const MAX_SIDES = 31

type Board struct {
	Fields [][]int
	Sides  [][]int
}

func readIndices(r io.Reader) []int {
	is := make([]int, 0)
	for {
		var i int
		n, _ := fmt.Fscanf(r, "%d", &i)
		if n != 1 {
			return nil
		}
		if i == 0 {
			break
		}
		is = append(is, i-1)
	}
	return is
}

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
		var j int
		var x, y float64
		if n, _ := fmt.Fscanf(r, "%d %f %f", &j, &x, &y); n < 3 || j != i+1 {
			return errors.New("board: invalid field header")
		}
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
