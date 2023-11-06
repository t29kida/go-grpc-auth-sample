package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/t29kida/go-grpc-auth-sample/internal/pb"
	"github.com/t29kida/go-grpc-auth-sample/internal/server"
	"github.com/t29kida/go-grpc-auth-sample/internal/server/interceptor"
	"github.com/t29kida/go-grpc-auth-sample/internal/service/auth"
	"github.com/t29kida/go-grpc-auth-sample/internal/service/hash"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(interceptor.RecoveryFunc),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_auth.UnaryServerInterceptor(interceptor.AuthFunc),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	pb.RegisterBackendServiceServer(s, server.New(auth.NewAuth(), hash.NewHash()))
	reflection.Register(s)

	go func() {
		log.Println("listening server with port:8080")
		if err := s.Serve(ln); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Println("stopping gRPC server...")
	s.GracefulStop()
	log.Println("grpc server shutdown completed")
}
