package main

import (
	"fmt"
	"gored/respser"
	"net"
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
	defer conn.Close()

	fmt.Println("handling response")
	for {
		buffer := make([]byte, 4096)
		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		in := string(buffer[:length])

		re, err := respser.RespDecode(in)
		response := ""
		if err != nil {
			fmt.Println("Error decoding:", err)
		} else {
			switch v := re.(type) {
			case *respser.SimpleString:
				fmt.Println("SimpleString:", v)
			case *respser.ErrorString:
				fmt.Println("ErrorString:", v)
			case *respser.Integer:
				fmt.Println("Integer:", v)
			case *respser.BulkString:
				fmt.Println("BulkString:", v)
			case *respser.Array:
				a := re.(*respser.Array)
				for _, e := range *a.Elements {
					switch v2 := e.(type) {
					case *respser.SimpleString:
						fmt.Println("SimpleString:", *v2)
					case *respser.ErrorString:
						fmt.Println("ErrorString:", *v2)
					case *respser.Integer:
						fmt.Println("Integer:", *v2)
					case *respser.BulkString:
						fmt.Println("BulkString:", *v2.S)
						response = handleCommand(*v2.S)
					case *respser.Array:
						fmt.Println("BulkString:", *v2)
					default:
						fmt.Println("Unknown type internal")
					}
					fmt.Println("Array:", v)
				}
			default:
				fmt.Println("Unknown type")
			}
		}

		if response == "" {
			response = "+OK\r\n"
		}
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}
		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
		}
	}
}

func handleCommand(s string) string {
	if s == "PING" {
		ss := respser.SimpleString{S: "PONG"}
		return ss.RespEncode()
	}
	return ""
}
