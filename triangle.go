package main

import (
	"fmt"
	"net/rpc"
	"os"
)

const usage = `Usage: triangle <command>
Available commands:
  server    Run the triangle server
  info      Print information about the played track
  toggle    Toggle play/pause
  next      Play the next track
  previous  Play the previous track`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}
	command := os.Args[1]
	switch command {
	case "server":
		server()
	case "toggle":
		toggle()
	case "info":
		info()
	default:
		fmt.Println("Error: Unknown command", command)
	}
}

func info() {
	var reply TriangleStatus
	err := call("TriangleServer.Info", &reply)
	if err != nil {
		fmt.Println("Error: Couldn't connect to Triangle server.", err)
		return
	}
	fmt.Println(reply)
}

func toggle() {
	call("TriangleServer.Toggle", nil)
}

func call(method string, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		return err
	}
	err = client.Call(method, 0, reply)
	if err != nil {
		return err
	}
	return nil
}
