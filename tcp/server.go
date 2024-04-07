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
			break Loop
		}

		data := parseData(b[:n])
		split := strings.Split(data, " ")
		args := split[1:]
		command := commands.ParseCommand(split[0])

		res := command.ValidateArgs(args)
		if res != "" {
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
			continue
		}

		switch command {
		case commands.Get:
			res = t.Store.Get(args[0])
			if res == "" {
				res = "key not found"
			}
		case commands.Set:
			t.Store.Set(args[0], args[1])
			res = "great success"
		case commands.Delete:
			t.Store.Delete(args[0])
			res = "great success"
		case commands.Exit:
			err = t.Store.WriteToDisk()
			if err != nil {
				slog.Error(err.Error())
			}
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
	store := store.NewKVStore()
	s := TCPServer{
		Listener: listener,
		Store:    store,
	}
	return &s, nil
}
