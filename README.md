# Go Chat Server

A modular TCP chat server written in Go with support for rooms, private messages and a blackjack mini-game.

## Features

- Lobby plus ad-hoc rooms with live room metadata (topics, population, creation time) and automatic cleanup when empty.
- Per-room topic management with `/topic`, `/rooms` and `/roominfo` commands.
- Embedded blackjack mini-game featuring chip bankrolls, configurable bets, payouts and blackjack bonuses.
- Extensible game registry for adding new interactive games.
- Container-first workflow with Docker and docker-compose.

## Getting Started

### Local development

```bash
make build
./bin/chatserver -addr :8080
```

Connect with `nc localhost 8080` or any TCP client.

### Docker

Build and run using docker-compose:

```bash
docker compose up --build
```

The server listens on `localhost:8080` by default and can be changed with the `CHATSERVER_ADDR` environment variable.

### Commands overview

```
/nick <name>         Change nickname
/join <room>         Join or create a room
/rooms               List available rooms with topics
/roominfo [room]     Display metadata for the current or provided room
/topic [new topic]   Show or set the current room topic
/list                Show users in the current room
/msg <user> <text>   Send a private message
/who <user>          Inspect another user
/blackjack ...       Play the blackjack mini-game (start, bet, hit, stand, status, bankroll, reset)
/games               List registered games
```

## Project layout

- `cmd/chatserver`: Application entrypoint.
- `internal/chat`: Core chat domain (clients, rooms, server lifecycle).
- `internal/games`: Game registry abstractions.
- `internal/games/blackjack`: Blackjack implementation and tests.

Run tests with `make test`.
