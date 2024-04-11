package main

import (
	"fmt"
	"net"

	"gored/respser"
)

func main() {
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :6380")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buffer := make([]byte, 4096)
	length, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	in := string(buffer[:length])
	fmt.Printf("input is:\n%s\n", in)

	re, err := respser.RespDecode(in)

	if err != nil {
		fmt.Println("Error decoding:", err)
	}

	fmt.Println(re.ToString())

	args := []respser.RespEncoder{}
	if arr, ok := re.(*respser.Array); ok {
		for _, e := range *arr.Elements {
			switch v2 := e.(type) {
			case *respser.SimpleString:
				args = append(args, v2)
				fmt.Println("SimpleString:", *v2)
			case *respser.ErrorString:
				args = append(args, v2)
				fmt.Println("ErrorString:", *v2)
			case *respser.Integer:
				args = append(args, v2)
				fmt.Println("Integer:", *v2)
			case *respser.BulkString:
				args = append(args, v2)
			case *respser.Array:
				fmt.Println("BulkString:", *v2)
			default:
				fmt.Println("Unknown type internal")
			}
		}
	}

	response := handleCommand(args)
	fmt.Printf("response is: \n%s\n", response)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}

func handleCommand(args []respser.RespEncoder) string {
	if len(args) == 0 {
		ss := respser.SimpleString{S: "OK"}
		return ss.RespEncode()
	}
	fmt.Println("args are:")
	for _, a := range args {
		fmt.Println(a.ToString())
	}

	if len(args) == 0 {
		return ""
	}

	if len(args) == 1 {
		if bs, ok := args[0].(*respser.BulkString); ok {
			var command string

			if bs.S == nil {
				return ""
			} else {
				command = *bs.S
			}
			switch command {
			case "PING":
				ss := respser.SimpleString{S: "PONG"}
				return ss.RespEncode()
			}
		}

	}
	if len(args) == 2 {
		if bs, ok := args[0].(*respser.BulkString); ok {
			var command string

			if bs.S == nil {
				return ""
			} else {
				command = *bs.S
			}
			switch command {
			case "ECHO":
				return args[1].RespEncode()
			}
		}
	}
	return ""
}
