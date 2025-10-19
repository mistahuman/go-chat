package main

import (
    "context"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    server := NewServer()
    
    listener, err := net.Listen("tcp", ":8080")
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
    
    log.Println("Server listening on :8080")
    
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