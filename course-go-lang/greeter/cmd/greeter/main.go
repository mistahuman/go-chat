package main

import (
    "os"
    "greeter/pkg/greeter"
)

func main() {
    name := "World"
    if len(os.Args) > 1 {
        name = os.Args[1]
    }
    greet.SayHello(name)
}
