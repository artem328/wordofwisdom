package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"wordofwisdom/internal/quote/zenquotes"
	"wordofwisdom/internal/server"
)

const defaultChallengeDifficulty uint = 20

func main() {
	var (
		port                int
		debug               bool
		challengeDifficulty uint
	)

	flag.IntVar(&port, "port", 9000, "The server's port")
	flag.BoolVar(&debug, "debug", false, "Start server in debug mode (more verbose logs)")
	flag.UintVar(&challengeDifficulty, "difficulty", defaultChallengeDifficulty, "The challenge difficulty (1-255)")
	flag.Parse()

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if challengeDifficulty < 1 || challengeDifficulty > 255 {
		slog.Warn(fmt.Sprintf("Invalid challenge difficulty %d. Setting to default value %d", challengeDifficulty, defaultChallengeDifficulty))
		challengeDifficulty = defaultChallengeDifficulty
	}

	failureChan := make(chan error, 1)

	stop, err := start(port, uint8(challengeDifficulty), failureChan)
	if err != nil {
		slog.Error("Could not start the server", "error", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		slog.Info(fmt.Sprintf("Received %s. Shutting Down. Press Ctrl+C to force stop", sig))
	case err := <-failureChan:
		slog.Error("Unexpected error. Shutting Down. Press Ctrl+C to force stop", "error", err)
	}

	done := make(chan struct{})
	go func(done chan<- struct{}) {
		stop()
		close(done)
	}(done)

	select {
	case <-sigCh:
		os.Exit(1)
	case <-done:
	}
}

func start(port int, challengeDifficulty uint8, failureChan chan<- error) (stop func(), err error) {
	addr := net.JoinHostPort("", strconv.Itoa(port))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not listen on %s: %w", addr, err)
	}

	qp := zenquotes.NewProvider()
	srv := server.New(l, qp, challengeDifficulty)

	go func(f chan<- error) {
		if err := srv.Start(); err != nil && !errors.Is(err, server.ErrStopped) {
			f <- err
		}
	}(failureChan)

	slog.Info(fmt.Sprintf("Listening for incoming connections on %s", addr))

	return func() {
		if err := srv.Stop(); err != nil {
			slog.Error("Failed to stop the server", "error", err)
		}
	}, nil
}
