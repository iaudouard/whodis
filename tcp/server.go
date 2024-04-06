package tcp

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/iaudouard/whodis/commands"
	"github.com/iaudouard/whodis/store"
)

type TCPServer struct {
	Listener net.Listener
	Store    store.KVStore
}

func (t TCPServer) Listen() {
	defer t.close()

	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			slog.Warn("error accepting connection: %w", err)
		}
		t.handleConn(conn)
	}
}

func (t TCPServer) close() {
	err := t.Store.WriteToDisk()
	if err != nil {
		slog.Error("failed to write to disk", "error", err)
	}
	t.Listener.Close()
}

func (t TCPServer) handleConn(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 1024)

Loop:
	for {
		n, err := conn.Read(b)
		if err != nil && err != io.EOF {
			slog.Error("connection goofed: %w", err)
		}

		data := parseData(b[:n])
		args := strings.Split(data, " ")
		command := commands.ParseCommand(args[0])

		res := ""
		switch command {
		case commands.Get:
			if len(args) != 2 {
				res = fmt.Sprintf("got %d arguments, expected 2", len(args))
				break
			}
			res = t.Store.Get(args[1])
			if res == "" {
				res = "key not found"
			}
		case commands.Set:
			if len(args) < 3 {
				res = fmt.Sprintf("got %d arguments, expected 3", len(args))
				break
			}
			t.Store.Set(args[1], args[2])
			res = "great success"
		case commands.Exit:
			break Loop
		case -1:
			res = "bad command"
		default:
			res = "not implemented"
		}

		conn.Write([]byte(fmt.Sprintf("%s\n", res)))
	}
}

func parseData(b []byte) string {
	b = bytes.Trim(b, "\x00")
	s := string(b)
	s = strings.TrimSpace(s)
	return s
}

func NewTCPServer(port uint16) (*TCPServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	slog.Info("started listening", "port", port)
	s := TCPServer{
		Listener: listener,
		Store:    store.NewKVStore(),
	}
	return &s, nil
}
