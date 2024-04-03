package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"
	"time"
	"wordofwisdom/internal/proto"
	"wordofwisdom/pkg/pow"
)

const readWriteTimeout = 3 * time.Second

var ErrStopped = errors.New("server stopped")

type Server struct {
	l                   net.Listener
	qp                  quoteProvider
	stopped             atomic.Bool
	challengeDifficulty uint8
}

func New(l net.Listener, qp quoteProvider, challengeDifficulty uint8) *Server {
	return &Server{l: l, qp: qp, challengeDifficulty: challengeDifficulty}
}

func (s *Server) Start() error {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			if s.stopped.Load() {
				return ErrStopped
			}

			return fmt.Errorf("accept connection: %s", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) Stop() error {
	s.stopped.Store(true)

	if err := s.l.Close(); err != nil {
		return fmt.Errorf("close listener: %w", err)
	}

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	logger := slog.With("addr", conn.RemoteAddr().String())

	logger.Info("New connection")

	p := proto.NewProto(conn)

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Failed to close connection", "error", err)
			return
		}

		logger.Info("Disconnected")
	}()

	if err := doChallenge(p, s.challengeDifficulty, logger); err != nil {
		logger.Error("Challenge failed", "error", err)
		return
	}

	msg, err := p.ReadTimeout(readWriteTimeout)
	if err != nil {
		logger.Error("Failed to read message", "error", err)
		return
	}

	s.handleMsg(p, msg, logger)
}

func doChallenge(p *proto.Proto, difficulty uint8, logger *slog.Logger) error {
	challenge, err := pow.Challenge()
	if err != nil {
		return fmt.Errorf("do challenge: prepare challenge: %w", err)
	}

	logger.Debug("Sending challenge", "challenge", challenge, "difficulty", difficulty)

	if err := p.WriteTimeout(readWriteTimeout, &proto.Challenge{
		Challenge:  challenge,
		Difficulty: difficulty,
	}); err != nil {
		return fmt.Errorf("do challenge: %w", err)
	}

	msg, err := p.ReadTimeout(readWriteTimeout)
	if err != nil {
		return fmt.Errorf("do challenge: %w", err)
	}

	resp, ok := msg.(*proto.ChallengeSolution)
	if !ok {
		return fmt.Errorf("do challenge: unexpected response")
	}

	logger.Debug("Validating challenge")
	validationStart := time.Now()
	valid := pow.Validate(challenge, resp.Nonce, difficulty)
	logger.Debug(fmt.Sprintf("Validation took %s", time.Since(validationStart)))
	if !valid {
		return fmt.Errorf("do challenge: invalid solution")
	}

	return nil
}

func (s *Server) handleMsg(p *proto.Proto, msg proto.Msg, logger *slog.Logger) {
	switch msg.(type) {
	case *proto.QuoteRequest:
		// should be a context related to the connection,
		// so long operations can be cancelled if connection was closed
		q, err := s.qp.GetQuote(context.TODO())
		if err != nil {
			logger.Error("Failed to fetch quote", "error", err)

			if err := p.WriteTimeout(readWriteTimeout, &proto.Error{Msg: "Something went wrong"}); err != nil {
				logger.Error("Failed to send error response", "error", err)
			}
			return
		}

		if err := p.WriteTimeout(readWriteTimeout, &proto.Quote{Quote: q}); err != nil {
			logger.Error("Failed to send quote", "error", err)
		}
	}
}
