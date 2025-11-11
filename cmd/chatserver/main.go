package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	chat "tcp-server/internal/chat"
)

func main() {
	defaultAddr := getEnv("CHATSERVER_ADDR", ":8080")
	addr := flag.String("addr", defaultAddr, "TCP address to listen on")
	flag.Parse()

	server := chat.NewServer()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		log.Println("Shutting down...")
		cancel()
		listener.Close()
	}()

	log.Printf("Server listening on %s", *addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}

		go server.HandleConnection(ctx, conn)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
