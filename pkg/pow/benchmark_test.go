package pow

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func BenchmarkProofOfWork(b *testing.B) {
	for d := uint8(1); d <= 24; d++ {
		b.Run(fmt.Sprintf("Solve-%d", d), func(b *testing.B) {
			r := rand.New(rand.NewPCG(0, 0))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				Solve(r.Uint32(), d)
			}
		})

		b.Run(fmt.Sprintf("SolveParallel-%d", d), func(b *testing.B) {
			r := rand.New(rand.NewPCG(0, 0))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				SolveParallel(r.Uint32(), d)
			}
		})

		b.Run(fmt.Sprintf("Validate-%d", d), func(b *testing.B) {
			r := rand.New(rand.NewPCG(0, 0))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				Validate(r.Uint32(), 5403, d)
			}
		})
	}
}
