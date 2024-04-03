package pow

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolve(t *testing.T) {
	const challenge uint32 = 962031355

	t.Parallel()

	for i := range 20 {
		t.Run(fmt.Sprintf("Difficulty %d", i+1), func(t *testing.T) {
			t.Parallel()
			assert.True(t, Validate(challenge, Solve(challenge, uint8(i+1)), uint8(i+1)))
		})
	}
}

func TestSolveParallel(t *testing.T) {
	const challenge uint32 = 962031355

	t.Parallel()

	for i := range 20 {
		t.Run(fmt.Sprintf("Difficulty %d", i+1), func(t *testing.T) {
			t.Parallel()
			assert.True(t, Validate(challenge, SolveParallel(challenge, uint8(i+1)), uint8(i+1)))
		})
	}
}

func Test_target(t *testing.T) {
	slice := func(n int, b byte) []byte {
		s := make([]byte, 32)
		s[n] = b
		return s
	}

	testCases := []struct {
		difficulty uint8
		target     []byte
	}{
		{difficulty: 1, target: slice(0, 0b_1000_0000)},
		{difficulty: 2, target: slice(0, 0b_0100_0000)},
		{difficulty: 3, target: slice(0, 0b_0010_0000)},
		{difficulty: 4, target: slice(0, 0b_0001_0000)},
		{difficulty: 5, target: slice(0, 0b_0000_1000)},
		{difficulty: 6, target: slice(0, 0b_0000_0100)},
		{difficulty: 7, target: slice(0, 0b_0000_0010)},
		{difficulty: 8, target: slice(0, 0b_0000_0001)},
		{difficulty: 9, target: slice(1, 0b_1000_0000)},
		{difficulty: 10, target: slice(1, 0b_0100_0000)},
		{difficulty: 11, target: slice(1, 0b_0010_0000)},
		{difficulty: 12, target: slice(1, 0b_0001_0000)},
		{difficulty: 13, target: slice(1, 0b_0000_1000)},
		{difficulty: 14, target: slice(1, 0b_0000_0100)},
		{difficulty: 15, target: slice(1, 0b_0000_0010)},
		{difficulty: 16, target: slice(1, 0b_0000_0001)},
		{difficulty: 17, target: slice(2, 0b_1000_0000)},
		{difficulty: 18, target: slice(2, 0b_0100_0000)},
		{difficulty: 19, target: slice(2, 0b_0010_0000)},
		{difficulty: 20, target: slice(2, 0b_0001_0000)},
		{difficulty: 21, target: slice(2, 0b_0000_1000)},
		{difficulty: 22, target: slice(2, 0b_0000_0100)},
		{difficulty: 23, target: slice(2, 0b_0000_0010)},
		{difficulty: 24, target: slice(2, 0b_0000_0001)},
		// the pattern should repeat
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("difficulty=%d", tt.difficulty), func(t *testing.T) {
			assert.Equal(t, tt.target, target(tt.difficulty))
		})
	}
}
