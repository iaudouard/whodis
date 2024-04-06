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
	listener net.Listener
	store    store.KVStore
}

func (t TCPServer) Listen() {
	defer t.close()

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			slog.Warn("error accepting connection: %w", err)
		}
		t.handleConn(conn)
	}
}

func (t TCPServer) close() {
	err := t.store.WriteToDisk()
	if err != nil {
		slog.Error("failed to write to disk", "error", err)
	}
	t.listener.Close()
}

func (t TCPServer) handleConn(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 1024)

Loop:
	for {
		_, err := conn.Read(b)
		if err != nil && err != io.EOF {
			slog.Error("connection goofed: %w", err)
		}

		data := parseData(b)
		args := strings.Split(data, " ")
		command := commands.ParseCommand(args[0])

		res := ""
		switch command {
		case commands.Get:
			if len(args) < 2 {
				res = "missing key to get"
				break
			}
			res = t.store.Get(args[1])
			if res == "" {
				res = "key not found"
			}
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
		listener: listener,
		store:    store.NewKVStore(),
	}
	return &s, nil
}
