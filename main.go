package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Szent7/medovukha-core/ipc"
	dockerpb "github.com/Szent7/medovukha-core/proto/docker/v1"
	"google.golang.org/grpc"
)

func main() {
	const socketPath = "/tmp/medovukha-core.sock"

	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			log.Fatalf("cannot remove existing socket %s: %s", socketPath, err.Error())
		}
	}

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("cannot listen on %s: %s\n", socketPath, err.Error())
	}

	if err := os.Chmod(socketPath, 0660); err != nil {
		log.Printf("warning: cannot chmod socket %s: %s\n", socketPath, err.Error())
	}

	dockerService, err := ipc.NewDockerCore()
	if err != nil {
		log.Fatalf("cannot init docker client: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	dockerpb.RegisterDockerServiceServer(grpcServer, dockerService)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	serveErrCh := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s\n", socketPath)
		serveErrCh <- grpcServer.Serve(lis)
	}()

	select {
	case sig := <-sigCh:
		log.Printf("received signal %s, shutting down...", sig)
	case err := <-serveErrCh:
		if err != nil {
			log.Fatalf("gRPC server error: %s", err.Error())
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	if err := lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("error closing listener: %s", err.Error())
	}
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		log.Printf("error removing socket file %s: %s", socketPath, err.Error())
	}
	if err := dockerService.Close(); err != nil {
		log.Printf("error while closing gRPC service (dockerService) %s: %s\n", socketPath, err.Error())
	}

	<-shutdownCtx.Done()
	log.Println("Shutdown complete")
}
