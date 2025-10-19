package main

import (
    "bufio"
    "context"
    "fmt"
    "net"
    "strings"
    "sync"
    "time"
)

type Client struct {
    conn      net.Conn
    nick      string
    room      *Room
    server    *Server
    joinedAt  time.Time
    mu        sync.RWMutex
    writeChan chan string
}

func NewClient(conn net.Conn, server *Server) *Client {
    c := &Client{
        conn:      conn,
        nick:      conn.RemoteAddr().String(),
        server:    server,
        joinedAt:  time.Now(),
        writeChan: make(chan string, 100),
    }
    go c.writer()
    return c
}

func (c *Client) Nick() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.nick
}

func (c *Client) SetNick(nick string) {
    c.mu.Lock()
    c.nick = nick
    c.mu.Unlock()
}

func (c *Client) Send(msg string) {
    select {
    case c.writeChan <- msg:
    default:
    }
}

func (c *Client) writer() {
    for msg := range c.writeChan {
        fmt.Fprintln(c.conn, msg)
    }
}

func (c *Client) Handle(ctx context.Context) {
    defer c.conn.Close()
    defer close(c.writeChan)
    
    scanner := bufio.NewScanner(c.conn)
    for scanner.Scan() {
        select {
        case <-ctx.Done():
            return
        default:
        }
        
        msg := scanner.Text()
        if strings.HasPrefix(msg, "/") {
            c.handleCommand(msg)
        } else {
            c.handleMessage(msg)
        }
    }
}

func (c *Client) handleMessage(msg string) {
    if c.room == nil {
        c.Send("Join a room first with /join <room>")
        return
    }
    c.room.Broadcast(fmt.Sprintf("<%s> %s", c.Nick(), msg))
}

func (c *Client) handleCommand(cmd string) {
    parts := strings.Fields(cmd)
    if len(parts) == 0 {
        return
    }
    
    switch parts[0] {
    case "/nick":
        if len(parts) < 2 {
            c.Send("Usage: /nick <name>")
            return
        }
        oldNick := c.Nick()
        c.SetNick(parts[1])
        if c.room != nil {
            c.room.Broadcast(fmt.Sprintf("* %s is now known as %s", oldNick, c.Nick()))
        }
        
    case "/join":
        if len(parts) < 2 {
            c.Send("Usage: /join <room>")
            return
        }
        
        if c.room != nil {
            c.room.Leave(c)
            c.room.Broadcast(fmt.Sprintf("* %s left", c.Nick()))
        }
        
        newRoom := c.server.GetOrCreateRoom(parts[1])
        newRoom.Join(c)
        c.room = newRoom
        c.Send(fmt.Sprintf("Joined room: %s", parts[1]))
        newRoom.Broadcast(fmt.Sprintf("* %s joined", c.Nick()))
        
    case "/leave":
        if c.room == nil {
            c.Send("Not in a room")
            return
        }
        c.room.Leave(c)
        c.room.Broadcast(fmt.Sprintf("* %s left", c.Nick()))
        c.room = nil
        c.Send("Left room")
        
    case "/rooms":
        rooms := c.server.ListRooms()
        c.Send("Available rooms:")
        for _, name := range rooms {
            c.Send(fmt.Sprintf("  - %s", name))
        }
        
    case "/list":
        if c.room == nil {
            c.Send("Not in a room")
            return
        }
        users := c.room.ListUsers()
        c.Send(fmt.Sprintf("Users in %s:", c.room.name))
        for _, nick := range users {
            c.Send(fmt.Sprintf("  - %s", nick))
        }
        
    case "/msg":
        if len(parts) < 3 {
            c.Send("Usage: /msg <user> <message>")
            return
        }
        target := c.server.FindClientByNick(parts[1])
        if target == nil {
            c.Send(fmt.Sprintf("User %s not found", parts[1]))
            return
        }
        message := strings.Join(parts[2:], " ")
        target.Send(fmt.Sprintf("[PM from %s] %s", c.Nick(), message))
        c.Send(fmt.Sprintf("[PM to %s] %s", parts[1], message))
        
    case "/who":
        if len(parts) < 2 {
            c.Send("Usage: /who <user>")
            return
        }
        target := c.server.FindClientByNick(parts[1])
        if target == nil {
            c.Send(fmt.Sprintf("User %s not found", parts[1]))
            return
        }
        duration := time.Since(target.joinedAt)
        roomName := "none"
        if target.room != nil {
            roomName = target.room.name
        }
        c.Send(fmt.Sprintf("%s - room: %s, online: %v", 
            target.Nick(), roomName, duration.Round(time.Second)))
        
    default:
        c.Send("Unknown command")
    }
}