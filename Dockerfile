FROM golang:1.21 as builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /chatserver ./cmd/chatserver

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /chatserver /usr/local/bin/chatserver
EXPOSE 8080
ENV CHATSERVER_ADDR=:8080
ENTRYPOINT ["/usr/local/bin/chatserver"]
