// cmd/prism/main.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/tamzrod/prism/internal/config"
	"github.com/tamzrod/prism/internal/engine"
	"github.com/tamzrod/prism/internal/protocol"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	var req protocol.Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		fmt.Fprintf(conn, `{"error_code":1}`)
		return
	}
	// Process already validates internally but we can catch early
	resp, err := engine.Process(&req)
	if err != nil {
		// should not happen; Process only returns marshal errors
		fmt.Fprintf(conn, `{"error_code":1}`)
		return
	}
	conn.Write(resp)
}

func main() {
	portFlag := flag.Int("port", 0, "port to listen on (overrides config)")
	flag.Parse()

	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config:", err)
		os.Exit(1)
	}
	port := cfg.Port
	if *portFlag != 0 {
		port = *portFlag
	}
	addr := fmt.Sprintf(":%d", port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	fmt.Println("listening on", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "accept:", err)
			continue
		}
		go handleConn(conn)
	}
}
