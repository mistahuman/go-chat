package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"tcp-server/blackjack"
)

const maxMessageSize = 1024 * 1024

type Client struct {
	conn      net.Conn
	nick      string
	room      *Room
	server    *Server
	joinedAt  time.Time
	mu        sync.RWMutex
	writeChan chan string
	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
	game      *blackjack.Game
}

func NewClient(ctx context.Context, conn net.Conn, server *Server) *Client {
	clientCtx, cancel := context.WithCancel(ctx)
	c := &Client{
		conn:      conn,
		nick:      conn.RemoteAddr().String(),
		server:    server,
		joinedAt:  time.Now(),
		writeChan: make(chan string, 100),
		ctx:       clientCtx,
		cancel:    cancel,
	}
	go func() {
		<-clientCtx.Done()
		_ = conn.Close()
	}()
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

func (c *Client) JoinedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.joinedAt
}

func (c *Client) Room() *Room {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.room
}

func (c *Client) setRoom(room *Room) {
	c.mu.Lock()
	c.room = room
	c.mu.Unlock()
}

func (c *Client) setGame(game *blackjack.Game) {
	c.mu.Lock()
	c.game = game
	c.mu.Unlock()
}

func (c *Client) currentGame() *blackjack.Game {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.game
}

func (c *Client) Send(msg string) {
	select {
	case c.writeChan <- msg:
	case <-c.ctx.Done():
		return
	default:
	}
}

func (c *Client) writer() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg, ok := <-c.writeChan:
			if !ok {
				return
			}
			if _, err := fmt.Fprintln(c.conn, msg); err != nil {
				c.Close()
				return
			}
		}
	}
}

func (c *Client) Handle() {
	defer c.Close()

	scanner := bufio.NewScanner(c.conn)
	scanner.Buffer(make([]byte, 0, 1024), maxMessageSize)
	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
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

	if err := scanner.Err(); err != nil {
		log.Printf("client %s disconnected: %v", c.Nick(), err)
	}
}

func (c *Client) handleMessage(msg string) {
	room := c.Room()
	if room == nil {
		c.Send("Join a room first with /join <room>")
		return
	}
	room.Broadcast(fmt.Sprintf("<%s> %s", c.Nick(), msg))
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
		if room := c.Room(); room != nil {
			room.Broadcast(fmt.Sprintf("* %s is now known as %s", oldNick, c.Nick()))
		}

	case "/join":
		if len(parts) < 2 {
			c.Send("Usage: /join <room>")
			return
		}

		if current := c.Room(); current != nil {
			current.Leave(c)
			current.Broadcast(fmt.Sprintf("* %s left", c.Nick()))
		}

		newRoom := c.server.GetOrCreateRoom(parts[1])
		newRoom.Join(c)
		c.setRoom(newRoom)
		c.Send(fmt.Sprintf("Joined room: %s", parts[1]))
		newRoom.Broadcast(fmt.Sprintf("* %s joined", c.Nick()))

	case "/leave":
		current := c.Room()
		if current == nil {
			c.Send("Not in a room")
			return
		}
		current.Leave(c)
		current.Broadcast(fmt.Sprintf("* %s left", c.Nick()))
		c.setRoom(nil)
		c.Send("Left room")

	case "/rooms":
		rooms := c.server.ListRooms()
		c.Send("Available rooms:")
		for _, name := range rooms {
			c.Send(fmt.Sprintf("  - %s", name))
		}

	case "/list":
		current := c.Room()
		if current == nil {
			c.Send("Not in a room")
			return
		}
		users := current.ListUsers()
		c.Send(fmt.Sprintf("Users in %s:", current.Name()))
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
		duration := time.Since(target.JoinedAt())
		roomName := "none"
		if room := target.Room(); room != nil {
			roomName = room.Name()
		}
		c.Send(fmt.Sprintf("%s - room: %s, online: %v",
			target.Nick(), roomName, duration.Round(time.Second)))

	case "/blackjack":
		c.handleBlackjack(parts)

	default:
		c.Send("Unknown command")
	}
}

func (c *Client) handleBlackjack(parts []string) {
	if len(parts) < 2 {
		c.Send("Usage: /blackjack <start|hit|stand|status|quit>")
		return
	}

	switch parts[1] {
	case "start":
		game := blackjack.NewGame(nil)
		c.setGame(game)
		c.Send("Blackjack game started.")
		c.Send(fmt.Sprintf("Dealer shows %s and [?]", game.Dealer[0].String()))
		c.Send(formatBlackjackHand("Your hand", game.Player))
		if game.IsFinished() {
			c.finishBlackjack(game)
			return
		}
		c.Send("Use /blackjack hit or /blackjack stand.")

	case "hit":
		game := c.currentGame()
		if game == nil {
			c.Send("Start a game first with /blackjack start")
			return
		}
		if game.IsFinished() {
			c.finishBlackjack(game)
			return
		}
		game.PlayerHit()
		c.Send("You drew a card.")
		c.Send(formatBlackjackHand("Your hand", game.Player))
		if game.IsFinished() {
			c.finishBlackjack(game)
			return
		}
		c.Send(fmt.Sprintf("Dealer shows %s and [?]", game.Dealer[0].String()))

	case "stand":
		game := c.currentGame()
		if game == nil {
			c.Send("Start a game first with /blackjack start")
			return
		}
		if game.IsFinished() {
			c.finishBlackjack(game)
			return
		}
		game.PlayerStand()
		c.finishBlackjack(game)

	case "status":
		game := c.currentGame()
		if game == nil {
			c.Send("No active blackjack game.")
			return
		}
		c.Send(formatBlackjackHand("Your hand", game.Player))
		c.Send(fmt.Sprintf("Dealer shows %s and [?]", game.Dealer[0].String()))

	case "quit":
		if c.currentGame() == nil {
			c.Send("No active blackjack game.")
			return
		}
		c.setGame(nil)
		c.Send("Blackjack game discarded.")

	default:
		c.Send("Usage: /blackjack <start|hit|stand|status|quit>")
	}
}

func (c *Client) finishBlackjack(game *blackjack.Game) {
	c.Send("Round finished.")
	c.Send(formatBlackjackHand("Dealer", game.Dealer))
	c.Send(formatBlackjackHand("Player", game.Player))
	switch game.Result {
	case blackjack.ResultPlayerWin:
		c.Send("Result: you win!")
	case blackjack.ResultDealerWin:
		c.Send("Result: dealer wins.")
	case blackjack.ResultPush:
		c.Send("Result: push.")
	default:
		c.Send("Result: no winner.")
	}
	c.setGame(nil)
}

func formatBlackjackHand(label string, hand blackjack.Hand) string {
	return fmt.Sprintf("%s: %s", label, hand.String())
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		_ = c.conn.Close()
	})
}
