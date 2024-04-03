package proto

import (
	"encoding/binary"
)

type Msg interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

var _ Msg = (*Challenge)(nil)

type Challenge struct {
	Challenge  uint32
	Difficulty uint8
}

func (c *Challenge) Encode() ([]byte, error) {
	buf := make([]byte, 5)
	binary.BigEndian.PutUint32(buf[:4], c.Challenge)
	buf[4] = c.Difficulty

	return buf, nil
}

func (c *Challenge) Decode(b []byte) error {
	if len(b) != 5 {
		return ErrMalformedMessage
	}

	c.Challenge = binary.BigEndian.Uint32(b[:4])
	c.Difficulty = b[4]

	return nil
}

var _ Msg = (*ChallengeSolution)(nil)

type ChallengeSolution struct {
	Nonce uint32
}

func (s *ChallengeSolution) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, s.Nonce)

	return buf, nil
}

func (s *ChallengeSolution) Decode(b []byte) error {
	if len(b) != 4 {
		return ErrMalformedMessage
	}

	s.Nonce = binary.BigEndian.Uint32(b)

	return nil
}

var _ Msg = (*QuoteRequest)(nil)

type QuoteRequest struct{}

func (*QuoteRequest) Encode() ([]byte, error) {
	return nil, nil
}

func (*QuoteRequest) Decode(b []byte) error {
	if len(b) > 0 {
		return ErrMalformedMessage
	}

	return nil
}

var _ Msg = (*Quote)(nil)

type Quote struct {
	Quote string
}

func (q *Quote) Encode() ([]byte, error) {
	return []byte(q.Quote), nil
}

func (q *Quote) Decode(b []byte) error {
	q.Quote = string(b)

	return nil
}

var _ Msg = (*Error)(nil)

type Error struct {
	Msg string
}

func (e *Error) Encode() ([]byte, error) {
	return []byte(e.Msg), nil
}

func (e *Error) Decode(b []byte) error {
	e.Msg = string(b)

	return nil
}
