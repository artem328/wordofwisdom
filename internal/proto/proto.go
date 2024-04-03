package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	headerSize  = 1 + 2     // 1 for message type and 2 for body size
	maxBodySize = 1<<16 - 1 // 2 bytes uint
)

var ErrMalformedMessage = errors.New("malformed message")
var ErrUnknownMsgType = errors.New("unknown message type")
var ErrBodyTooBig = errors.New("body too big")

type Proto struct {
	conn      net.Conn
	headerBuf []byte
}

func NewProto(conn net.Conn) *Proto {
	return &Proto{
		conn:      conn,
		headerBuf: make([]byte, headerSize),
	}
}

func (p *Proto) Read() (Msg, error) {
	return p.ReadTimeout(0)
}

func (p *Proto) ReadTimeout(timeout time.Duration) (Msg, error) {
	if err := p.setReadTimeout(timeout); err != nil {
		return nil, fmt.Errorf("set read deadline: %w", err)
	}

	n, err := p.conn.Read(p.headerBuf)
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	if n != headerSize {
		return nil, fmt.Errorf("read header: %w", ErrMalformedMessage)
	}

	size := int(binary.BigEndian.Uint16(p.headerBuf[1:3]))
	body := make([]byte, size)

	if err = p.setReadTimeout(timeout); err != nil {
		return nil, fmt.Errorf("set read deadline: %w", err)
	}

	n, err = p.conn.Read(body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if n != size {
		return nil, fmt.Errorf("read body: %w", ErrMalformedMessage)
	}

	return Decode(MsgType(p.headerBuf[0]), body)
}

func (p *Proto) Write(msg Msg) error {
	return p.WriteTimeout(0, msg)
}

func (p *Proto) WriteTimeout(timeout time.Duration, msg Msg) error {
	msgType, body, err := Encode(msg)
	if err != nil {
		return err
	}

	size := len(body)
	if size > maxBodySize {
		return ErrBodyTooBig
	}

	p.headerBuf[0] = byte(msgType)
	binary.BigEndian.PutUint16(p.headerBuf[1:3], uint16(size))

	if err := p.setWriteTimeout(timeout); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	out := io.MultiReader(bytes.NewReader(p.headerBuf), bytes.NewReader(body))

	_, err = io.Copy(p.conn, out)

	return err
}

func (p *Proto) setReadTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return nil
	}

	return p.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (p *Proto) setWriteTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return nil
	}

	return p.conn.SetWriteDeadline(time.Now().Add(timeout))
}
