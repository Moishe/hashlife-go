// Standard implementation of Life in Go
package hashlife

import (
	"math"
)

// Computes the inner (n / 2 x n / 2) cells, n / 4 generations in the future
func SimpleNthGeneration(board []byte) []byte {
	buffer := make([][]byte, 2)
	buffer[0] = make([]byte, len(board))
	buffer[1] = make([]byte, len(board))
	for i := 0; i < len(board); i++ {
		buffer[0][i] = board[i]
	}

	size := int(math.Sqrt(float64(len(board))))
	future_steps := size / 4
	for i := 0; i < future_steps; i++ {
		SimpleNextGeneration(&buffer[i%2], &buffer[(i+1)%2], size-(i*2))
	}
	last_size := size / 2
	last_size = last_size * last_size
	board = buffer[future_steps%2][:last_size]
	return board
}

// Computes the inner n - 2 x n - 2 cells, 1 generation in the future.
func SimpleNextGeneration(board *[]byte, new_board *[]byte, size int) {
	new_size := size - 2
	// Slice up the board for number of cores
	var cores = 4
	cells_per_core := (size - 2) * (size - 2) / cores
	done := make(chan bool)
	for core := 0; core < cores; core++ {
		this_core := core
		go func() {
			for k := 0; k < cells_per_core; k++ {
				idx := this_core*cells_per_core + k
				i := idx/(size-2) + 1
				j := idx%(size-2) + 1
				if i > size-2 || j > size-2 {
					break
				}

				c := 0
				for ox := -1; ox <= 1; ox++ {
					for oy := -1; oy <= 1; oy++ {
						if (*board)[(i+ox)+(j+oy)*size] == 1 && (ox != 0 || oy != 0) {
							c++
						}
					}
				}
				new_idx := i - 1 + (j-1)*new_size
				if c == 3 {
					(*new_board)[new_idx] = 1
				} else if c == 2 {
					(*new_board)[new_idx] = (*board)[i+j*size]
				} else {
					(*new_board)[new_idx] = 0
				}
			}
			done <- true
		}()
	}
	for i := 0; i < cores; i++ {
		<-done
	}
}
