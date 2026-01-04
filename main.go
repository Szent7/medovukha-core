package main

import (
	"log"
	"net"
	"os"

	"github.com/Szent7/medovukha-core/ipc"
	dockerpb "github.com/Szent7/medovukha-core/proto/docker/v1"
	"google.golang.org/grpc"
)

func main() {
	const socketPath = "/tmp/medovukha-core.sock"

	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("cannot listen on %s: %s\n", socketPath, err.Error())
	}
	defer lis.Close()
	os.Chmod(socketPath, 0660)

	dockerService, err := ipc.NewDockerCore()
	if err != nil {
		log.Fatalf("cannot init docker client: %s", err.Error())
	}
	defer dockerService.Close()

	s := grpc.NewServer()
	dockerpb.RegisterDockerServiceServer(s, dockerService)

	log.Printf("gRPC server listening on %s\n", socketPath)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %s", err.Error())
	}
}
