package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func startRedisProtocol(host string) {
	conn, err := net.DialTimeout("tcp4", host, 300*time.Millisecond)
	if err != nil {
		host = "not connected"
	} else {
		defer conn.Close()
	}

	for {
		fmt.Printf("%s>", host)

		stdio := bufio.NewReader(os.Stdin)

		argsData, _, _ := stdio.ReadLine()
		args := strings.Fields(string(argsData))
		if len(args) < 1 {
			continue
		}

		if strings.EqualFold(args[0], "exit") || strings.EqualFold(args[0], "quit") {
			break
		}

		if err != nil {
			continue
		}

		protoData := fmt.Sprintf("*%d\r\n", len(args))
		for _, arg := range args {
			protoData += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
		}

		conn.Write([]byte(protoData))

		r := bufio.NewReader(conn)
		conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))

		prefix, _ := r.ReadByte()
		switch prefix {
		case '+', '-', ':':
			result, _, _ := r.ReadLine()
			fmt.Printf("%s\n", result)
		case '$':
			data, _, _ := r.ReadLine()
			bytesCount, _ := strconv.Atoi(string(data[:]))
			if bytesCount < 1 {
				fmt.Println("(nil)")
				continue
			}
			result, _, _ := r.ReadLine()
			fmt.Printf("\"%s\"\n", result)
		case '*':
			data, _, _ := r.ReadLine()
			resultCount, _ := strconv.Atoi(string(data[:]))
			for i := 1; i <= resultCount; i++ {
				r.ReadLine()
				result, _, _ := r.ReadLine()
				fmt.Printf("%d) \"%s\"\n", i, result)
			}
		default:
			fmt.Println("Error Procotol")
		}
	}
}

func main() {
	var host = "127.0.0.1:6379"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	startRedisProtocol(host)
}
