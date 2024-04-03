package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"
	"wordofwisdom/internal/proto"
	"wordofwisdom/pkg/pow"
)

var errFailed = errors.New("failed")

func main() {
	var (
		addr    string
		verbose bool
	)

	flag.StringVar(&addr, "addr", "", "The server's address")
	flag.BoolVar(&verbose, "verbose", false, "Print extra logs")
	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if addr == "" {
		slog.Error("Empty address")
		os.Exit(1)
	}

	startTime := time.Now()

	if err := start(addr); err != nil {
		os.Exit(1)
	}

	slog.Debug(fmt.Sprintf("Finished in %s", time.Since(startTime)))
}

func start(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		slog.Error("Could not connect", "addr", addr, "error", err)
		return errFailed
	}

	slog.Debug("Connected", "addr", conn.RemoteAddr().String())

	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("Failed to close connection", "error", err)
			return
		}

		slog.Debug("Disconnected")
	}()

	p := proto.NewProto(conn)

	for {
		msg, err := p.Read()
		if err != nil {
			slog.Error("Could not read message", "error", err)
			return errFailed
		}

		switch m := msg.(type) {
		case *proto.Quote:
			fmt.Println("The quote received from the server")
			fmt.Println()
			fmt.Println(m.Quote)
			return nil
		case *proto.Challenge:
			slog.Debug("Received challenge", "challenge", m.Challenge, "difficulty", m.Difficulty)
			startTime := time.Now()
			var nonce uint32
			if m.Difficulty <= 6 {
				// with difficulty up to 6 bits single-thread solution is faster
				// (determined empirically on 8 CPUs)
				nonce = pow.Solve(m.Challenge, m.Difficulty)
			} else {
				nonce = pow.SolveParallel(m.Challenge, m.Difficulty)
			}

			slog.Debug(fmt.Sprintf("The challenge solved in %s", time.Since(startTime)), "solution", nonce)

			if err := p.Write(&proto.ChallengeSolution{Nonce: nonce}); err != nil {
				slog.Error("Failed to write challenge response", "error", err)
				return errFailed
			}

			if err := p.Write(&proto.QuoteRequest{}); err != nil {
				slog.Error("Failed to send request", "error", err)
				return errFailed
			}
		case *proto.Error:
			slog.Error("Server returned error", "msg", m.Msg)
			return errFailed
		default:
			slog.Debug(fmt.Sprintf("Received unexpected msg %T", msg))
		}
	}
}
