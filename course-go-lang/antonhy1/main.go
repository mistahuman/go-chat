package main

import (
	"fmt"
)

// func main() {
// 	userch := make(chan string, 2)

// 	// go func() {
// 	// 	time.Sleep(1 * time.Second)
// 	// 	userch <- "Bob"
// 	// }()

// 	userch <- "Bob"
// 	userch <- "Alice"

// 	user := <-userch

// 	fmt.Println(user)

// 	// user = <-userch
// 	// sendMessage(userch)

// }

func main() {
	server := NewServer()
	server.Start()

	server.quitch <- struct{}{}
}

type Server struct {
	users  map[string]string
	userch chan string
	quitch chan struct{}
}

func NewServer() *Server {
	return &Server{
		users:  make(map[string]string),
		userch: make(chan string),
		quitch: make(chan struct{}),
	}
}

func (s *Server) Start() {
	go s.loop()
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.userch:
			fmt.Println(msg)
		case <-s.quitch:
			fmt.Println("server need to quint")
			return
		default:
			return
		}

		// user := <-s.userch
		// s.users[user] = user

		// // fmt.Printf("adding new user %s\n %s", user, s.users)
		// fmt.Printf("adding new user %s\n", user)
	}
}

func (s *Server) addUser(user string) {
	s.users[user] = user
}

func sendMessage(msgch chan<- string) {
	// msg := <-msgch
	msgch <- "hello!"
	// fmt.Println(msg)
}

func readMessage(msgch chan string) {
	msg := <-msgch
	fmt.Println(msg)
}
