package proto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChallenge_Encode(t *testing.T) {
	e, err := (&Challenge{Challenge: 1<<4<<23 + 1<<3<<15 + 1<<2<<7 + 1, Difficulty: 59}).Encode()

	assert.NoError(t, err)
	assert.Equal(t, []byte{0b1000, 0b100, 0b10, 0b1, 59}, e)
}

func TestChallenge_Decode(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		var c Challenge

		err := c.Decode([]byte{0b1000, 0b100, 0b10, 0b1, 59})

		assert.NoError(t, err)
		assert.Equal(t, Challenge{Challenge: 1<<4<<23 + 1<<3<<15 + 1<<2<<7 + 1, Difficulty: 59}, c)
	})

	invalidLenTestCases := [][]byte{
		nil, make([]byte, 1), make([]byte, 2), make([]byte, 3), make([]byte, 4), make([]byte, 6),
	}

	for _, tt := range invalidLenTestCases {
		t.Run(fmt.Sprintf("Invalid Length=%d", len(tt)), func(t *testing.T) {
			var c Challenge

			err := c.Decode(tt)

			assert.ErrorIs(t, err, ErrMalformedMessage)
			assert.Empty(t, c)
		})
	}
}

func TestChallengeSolution_Encode(t *testing.T) {
	e, err := (&ChallengeSolution{Nonce: 1<<4<<23 + 1<<3<<15 + 1<<2<<7 + 1}).Encode()

	assert.NoError(t, err)
	assert.Equal(t, []byte{0b1000, 0b100, 0b10, 0b1}, e)
}

func TestChallengeSolution_Decode(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		var cs ChallengeSolution

		err := cs.Decode([]byte{0b1000, 0b100, 0b10, 0b1})

		assert.NoError(t, err)
		assert.Equal(t, ChallengeSolution{Nonce: 1<<4<<23 + 1<<3<<15 + 1<<2<<7 + 1}, cs)
	})

	invalidLenTestCases := [][]byte{
		nil, make([]byte, 1), make([]byte, 2), make([]byte, 3), make([]byte, 5),
	}

	for _, tt := range invalidLenTestCases {
		t.Run(fmt.Sprintf("Invalid Length=%d", len(tt)), func(t *testing.T) {
			var cs ChallengeSolution

			err := cs.Decode(tt)

			assert.ErrorIs(t, err, ErrMalformedMessage)
			assert.Empty(t, cs)
		})
	}
}

func TestQuoteRequest_Encode(t *testing.T) {
	e, err := (&QuoteRequest{}).Encode()

	assert.NoError(t, err)
	assert.Empty(t, e)
}

func TestQuoteRequest_Decode(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		var qr QuoteRequest

		err := qr.Decode(nil)

		assert.NoError(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		var qr QuoteRequest

		err := qr.Decode(make([]byte, 1))

		assert.ErrorIs(t, err, ErrMalformedMessage)
	})
}

func TestQuote_Encode(t *testing.T) {
	e, err := (&Quote{Quote: "Hello"}).Encode()

	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello"), e)
}

func TestQuote_Decode(t *testing.T) {
	var q Quote

	err := q.Decode([]byte("Hello"))

	assert.NoError(t, err)
	assert.Equal(t, Quote{Quote: "Hello"}, q)
}

func TestError_Encode(t *testing.T) {
	e, err := (&Error{Msg: "Good Bye"}).Encode()

	assert.NoError(t, err)
	assert.Equal(t, []byte("Good Bye"), e)
}

func TestError_Decode(t *testing.T) {
	var e Error

	err := e.Decode([]byte("Good Bye"))

	assert.NoError(t, err)
	assert.Equal(t, Error{Msg: "Good Bye"}, e)
}
