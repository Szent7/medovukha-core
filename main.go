package main

import (
	"log"
	"medovukha/ipc"
	"net"
	"net/rpc"
	"os"
)

func main() {
	const socketPath = "/tmp/medovukha-core.sock"
	const serverName = "MedovukhaCore"

	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("cannot listen on %s: %s\n", socketPath, err.Error())
	}
	defer ln.Close()
	os.Chmod(socketPath, 0660)

	core, err := ipc.NewDockerCore()
	if err != nil {
		log.Fatalf("cannot init docker client: %s", err.Error())
	}
	defer core.Close()

	rpcServer := rpc.NewServer()
	if err := rpcServer.RegisterName(serverName, core); err != nil {
		log.Fatalf("cannot register rpc: %s", err.Error())
	}

	log.Printf("%s listening on %s\n", serverName, socketPath)
	for {
		conn, err := ln.Accept()
		log.Printf("request from: %s\n", conn.RemoteAddr().String())
		if err != nil {
			log.Printf("accept error: %s\n", err.Error())
			continue
		}
		go rpcServer.ServeConn(conn)
	}
}
