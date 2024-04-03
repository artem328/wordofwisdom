package pow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"sync"
)

// Challenge generates random uint32 challenge
func Challenge() (uint32, error) {
	chBuf := make([]byte, 4)
	n, err := rand.Read(chBuf)
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("generate challenge: %w", err)
	}

	if n != len(chBuf) {
		return 0, fmt.Errorf("generate challenge: not enough random bytes")
	}

	return binary.BigEndian.Uint32(chBuf), nil
}

// Validate tests if challenge and nonce produces the hash that satisfies the target difficulty
func Validate(challenge, nonce uint32, difficulty uint8) bool {
	return validate(difficulty, challenge, nonce, nil)
}

// Solve finds a solution for the challenge for the target difficulty
func Solve(challenge uint32, difficulty uint8) uint32 {
	var nonce uint32
	buf := newBuf()

	for !validate(difficulty, challenge, nonce, buf) {
		nonce++
		if nonce == 0 { // overflow
			return 0
		}
	}

	return nonce
}

// SolveParallel finds a solution for the challenge for the target difficulty on multiple CPUs
func SolveParallel(challenge uint32, difficulty uint8) uint32 {
	var (
		solution = make(chan uint32, 1)
		done     = make(chan struct{})
	)
	defer close(done)

	var (
		workers = runtime.GOMAXPROCS(-1)
		chunk   = uint32(1 << 32 / workers)

		from uint32
		to   = from + chunk - 1
	)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(nonce, nonceMax uint32, buf []byte) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
				}

				if validate(difficulty, challenge, nonce, buf) {
					break
				}

				nonce++
				if nonce > nonceMax {
					return
				}
			}

			select {
			case <-done:
				return
			case solution <- nonce:
			}
		}(from, to, newBuf())

		from += chunk
		to += chunk
	}

	go func() {
		wg.Wait()
		select {
		case <-done:
		case solution <- 0: // neither worker found a solution
		}
	}()

	return <-solution
}

var targets [][]byte

func init() {
	targets = make([][]byte, 255)

	for i := uint8(0); i < 255; i++ {
		targets[i] = target(i + 1)
	}
}

// target creates a byte slice which contains single 1-bit at the position of target difficulty.
func target(difficulty uint8) []byte {
	if difficulty == 0 {
		panic("difficulty must be positive number")
	}

	t := make([]byte, 32)

	n := (difficulty - 1) / 8 // number of whole bytes containing zeros
	bitsShift := difficulty - n*8

	t[n] = 1 << (8 - bitsShift)

	return t
}

const bufLen = 8

func validate(difficulty uint8, challenge, nonce uint32, buf []byte) bool {
	if difficulty <= 0 {
		return true
	}

	if len(buf) < bufLen {
		buf = newBuf()
	}

	binary.BigEndian.PutUint32(buf[:4], challenge)
	binary.BigEndian.PutUint32(buf[4:], nonce)
	s := sha256.Sum256(buf)

	// test that resulting hash sum is less than target
	// meaning has more zero MSB than target
	// e.g. assume difficulty = 6
	// the target is 00000100 00000000 ...
	// then solution should have at least 6 most significant bits equal to zero,
	// thus it will be always less than target
	return bytes.Compare(s[:], targets[difficulty-1]) == -1
}

func newBuf() []byte {
	return make([]byte, bufLen)
}
