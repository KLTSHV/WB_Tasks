package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("not enough arguments")
	}
	if n, err := strconv.Atoi(args[1]); err != nil || n < 0 || n > 65536 {
		log.Fatal("port must be positive integer in the range 0-65536")
	}
	addr := args[0] + ":" + args[1]
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(*timeout))
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to", addr)

	var wg sync.WaitGroup

	wg.Add(2)
	//stdin
	go func() {
		defer wg.Done()
		_, _ = io.Copy(conn, os.Stdin)
		conn.Close() //CTRl+D

	}()
	//stdout
	go func() {
		defer wg.Done()
		_, _ = io.Copy(os.Stdout, conn)
	}()
	wg.Wait()
}
